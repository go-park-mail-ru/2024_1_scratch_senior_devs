package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/hub"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/path"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/protection"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/recover"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/loadtls"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"

	_ "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	metricsmw "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/metrics"

	grpcAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
	grpcNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"

	authDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/http"

	noteDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/http"
	noteRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/repo"

	attachDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach/delivery/http"
	attachRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach/repo"
	attachUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach/usecase"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
}

// @title 			YouNote API
// @version 		1.0
// @description 	API for YouNote service
// @host 			you-note.ru
func main() {
	logFile, err := os.OpenFile(os.Getenv("MAIN_LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("error opening log file: " + err.Error())
		return
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(logFile, os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.LoadConfig(os.Getenv("CONFIG_FILE"), logger)

	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("error connecting to postgres: " + err.Error())
		return
	}
	defer db.Close()

	tlsCredentials, err := loadtls.LoadTLSClientCredentials()
	if err != nil {
		logger.Error("fail to load TLS credentials" + err.Error())
	}

	authConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.Grpc.AuthIP, cfg.Grpc.AuthPort),

		grpc.WithTransportCredentials(tlsCredentials),
	)
	if err != nil {
		logger.Error("fail grpc.Dial auth: " + err.Error())
		return
	}
	defer authConn.Close()

	noteConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.Grpc.NoteIP, cfg.Grpc.NotePort),
		grpc.WithTransportCredentials(tlsCredentials),
	)
	if err != nil {
		logger.Error("fail grpc.Dial note: " + err.Error())
		return
	}
	defer noteConn.Close()

	Metrics, err := metrics.NewHttpMetrics("main")
	if err != nil {
		logger.Error("can`t create metrics (main http): " + err.Error())
	}

	postgresMetrics, err := metrics.NewDatabaseMetrics("postgres", "main")
	if err != nil {
		logger.Error("can`t create metrics (main postgres): " + err.Error())
	}

	websocketMetrics, err := metrics.NewWebsocketMetrics()
	if err != nil {
		logger.Error("can`t create metrics (websockets): " + err.Error())
	}

	NoteBaseRepo := noteRepo.CreateNotePostgres(db, &postgresMetrics)
	NoteHub := hub.NewHub(NoteBaseRepo, cfg.Hub, websocketMetrics)

	AttachRepo := attachRepo.CreateAttachRepo(db, &postgresMetrics)
	AttachUsecase := attachUsecase.CreateAttachUsecase(AttachRepo, NoteBaseRepo)
	AttachDelivery := attachDelivery.CreateAttachHandler(AttachUsecase, cfg.Attach)

	AuthClient := grpcAuth.NewAuthClient(authConn)
	NoteClient := grpcNote.NewNoteClient(noteConn)

	AuthDelivery := authDelivery.CreateAuthHandler(AuthClient, NoteClient, cfg.AuthHandler, cfg.Validation)
	NoteDelivery := noteDelivery.CreateNotesHandler(NoteClient, AuthClient, NoteHub)

	JwtMiddleware := protection.CreateJwtMiddleware(cfg.AuthHandler.Jwt)
	JwtWebsocketMiddleware := protection.CreateJwtWebsocketMiddleware(cfg.AuthHandler.Jwt)
	CsrfMiddleware := protection.CreateCsrfMiddleware(cfg.AuthHandler.Csrf)
	MetricsMiddleware := metricsmw.CreateHttpMetricsMiddleware(Metrics)
	LogMiddleware := log.CreateLogMiddleware(logger)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Error("query to path: " + r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})

	r.Use(
		LogMiddleware,
		MetricsMiddleware,
		protection.CorsMiddleware,
		recover.RecoverMiddleware,
	)

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet, http.MethodOptions)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.Handle("/signup", http.HandlerFunc(AuthDelivery.SignUp)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/login", http.HandlerFunc(AuthDelivery.SignIn)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/logout", JwtMiddleware(CsrfMiddleware(http.HandlerFunc(AuthDelivery.LogOut)))).Methods(http.MethodDelete, http.MethodOptions)
		auth.Handle("/check_user", JwtMiddleware(http.HandlerFunc(AuthDelivery.CheckUser))).Methods(http.MethodGet, http.MethodOptions)
		auth.Handle("/get_qr", JwtMiddleware(http.HandlerFunc(AuthDelivery.GetQRCode))).Methods(http.MethodGet, http.MethodOptions)
		auth.Handle("/disable_2fa", JwtMiddleware(CsrfMiddleware(http.HandlerFunc(AuthDelivery.DisableSecondFactor)))).Methods(http.MethodDelete, http.MethodOptions)
	}

	note := r.PathPrefix("/note").Subrouter()
	note.Use(protection.ReadAndCloseBody, CsrfMiddleware)
	{
		note.Handle("/get_all", JwtMiddleware(http.HandlerFunc(NoteDelivery.GetAllNotes))).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/{id}", JwtMiddleware(http.HandlerFunc(NoteDelivery.GetNote))).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/add", JwtMiddleware(http.HandlerFunc(NoteDelivery.AddNote))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/edit", JwtMiddleware(http.HandlerFunc(NoteDelivery.UpdateNote))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/delete", JwtMiddleware(http.HandlerFunc(NoteDelivery.DeleteNote))).Methods(http.MethodDelete, http.MethodOptions)
		note.Handle("/{id}/add_attach", JwtMiddleware(http.HandlerFunc(AttachDelivery.AddAttach))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/add_subnote", JwtMiddleware(http.HandlerFunc(NoteDelivery.CreateSubNote))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/add_collaborator", JwtMiddleware(http.HandlerFunc(NoteDelivery.AddCollaborator))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/subscribe_on_updates", JwtWebsocketMiddleware(http.HandlerFunc(NoteDelivery.SubscribeOnUpdates))).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/{id}/add_tag", JwtMiddleware(http.HandlerFunc(NoteDelivery.AddTag))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/delete_tag", JwtMiddleware(http.HandlerFunc(NoteDelivery.DeleteTag))).Methods(http.MethodDelete, http.MethodOptions)
		note.Handle("/{id}/set_icon", JwtMiddleware(http.HandlerFunc(NoteDelivery.SetIcon))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/set_header", JwtMiddleware(http.HandlerFunc(NoteDelivery.SetHeader))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/add_fav", JwtMiddleware(http.HandlerFunc(NoteDelivery.AddFavorite))).Methods(http.MethodPut, http.MethodOptions)
		note.Handle("/{id}/del_fav", JwtMiddleware(http.HandlerFunc(NoteDelivery.DeleteFavorite))).Methods(http.MethodPut, http.MethodOptions)
		note.Handle("/{id}/set_public", JwtMiddleware(http.HandlerFunc(NoteDelivery.SetPublic))).Methods(http.MethodPut, http.MethodOptions)
		note.Handle("/{id}/set_private", JwtMiddleware(http.HandlerFunc(NoteDelivery.SetPrivate))).Methods(http.MethodPut, http.MethodOptions)
		note.Handle("/{id}/make_zip", JwtMiddleware(http.HandlerFunc(NoteDelivery.ExportZip))).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/subscribe/on_invites", JwtWebsocketMiddleware(http.HandlerFunc(NoteDelivery.SubscribeOnInvites))).Methods(http.MethodGet, http.MethodOptions)
	}

	shared := r.PathPrefix("/shared").Subrouter()
	{
		shared.Handle("/note/{id}", http.HandlerFunc(NoteDelivery.GetPublicNote)).Methods(http.MethodGet, http.MethodOptions)
		shared.Handle("/note/{id}/make_zip", http.HandlerFunc(NoteDelivery.ExportZip)).Methods(http.MethodPost, http.MethodOptions)
		shared.Handle("/attach/{id}", http.HandlerFunc(AttachDelivery.GetSharedAttach)).Methods(http.MethodGet, http.MethodOptions)
	}

	profile := r.PathPrefix("/profile").Subrouter()
	profile.Use(protection.ReadAndCloseBody, JwtMiddleware, CsrfMiddleware)
	{
		profile.Handle("/get", http.HandlerFunc(AuthDelivery.GetProfile)).Methods(http.MethodGet, http.MethodOptions)
		profile.Handle("/update", http.HandlerFunc(AuthDelivery.UpdateProfile)).Methods(http.MethodPost, http.MethodOptions)
		profile.Handle("/update_avatar", http.HandlerFunc(AuthDelivery.UpdateProfileAvatar)).Methods(http.MethodPost, http.MethodOptions)
	}

	attach := r.PathPrefix("/attach").Subrouter()
	attach.Use(JwtMiddleware, CsrfMiddleware)
	{
		attach.Handle("/{id}", http.HandlerFunc(AttachDelivery.GetAttach)).Methods(http.MethodGet, http.MethodOptions)
		attach.Handle("/{id}/delete", http.HandlerFunc(AttachDelivery.DeleteAttach)).Methods(http.MethodDelete, http.MethodOptions)
	}

	tags := r.PathPrefix("/tags").Subrouter()
	tags.Use(JwtMiddleware, CsrfMiddleware)
	{
		tags.Handle("", http.HandlerFunc(NoteDelivery.GetTags)).Methods(http.MethodGet, http.MethodOptions)
		tags.Handle("/remember", http.HandlerFunc(NoteDelivery.RememberTag)).Methods(http.MethodPost, http.MethodOptions)
		tags.Handle("/forget", http.HandlerFunc(NoteDelivery.ForgetTag)).Methods(http.MethodDelete, http.MethodOptions)
		tags.Handle("/update", http.HandlerFunc(NoteDelivery.UpdateTag)).Methods(http.MethodPost, http.MethodOptions)
	}

	export := r.PathPrefix("/export_to_pdf").Subrouter()
	export.Use(protection.ReadAndCloseBody)
	{
		export.Handle("", http.HandlerFunc(NoteDelivery.ExportToPDF)).Methods(http.MethodPost, http.MethodOptions)
	}

	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	http.Handle("/", r)

	go NoteHub.Run(context.WithValue(context.Background(), config.LoggerContextKey, logger))
	go NoteHub.StartCache(context.WithValue(context.Background(), config.LoggerContextKey, logger))
	go NoteHub.StartCacheMain(context.WithValue(context.Background(), config.LoggerContextKey, logger))

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Handler:           path.PathMiddleware(r),
		Addr:              fmt.Sprintf(":%s", cfg.Main.Port),
		ReadTimeout:       cfg.Main.ReadTimeout,
		WriteTimeout:      cfg.Main.WriteTimeout,
		ReadHeaderTimeout: cfg.Main.ReadHeaderTimeout,
		IdleTimeout:       cfg.Main.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Info("Server stopped")
		}
	}()
	logger.Info("Server started")

	sig := <-signalCh
	logger.Info("Received signal: " + sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Main.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed: " + err.Error())
	}
}
