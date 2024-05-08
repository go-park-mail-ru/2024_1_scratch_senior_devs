package log

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestLogMiddleware_ServerLogsInterceptor(t *testing.T) {

	tests := []struct {
		name    string
		handler func(ctx context.Context, req interface{}) (interface{}, error)
		wantErr error
	}{
		{
			name: "TestLogInterceptor_Success",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, nil
			},
			wantErr: nil,
		},
		{
			name: "TestLogInterceptor_Fail",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, errors.New("error")
			},
			wantErr: errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))

			interceptor := NewGrpcLogMw(logger)

			ctx := context.Background()
			req := httptest.NewRequest(http.MethodGet, "http://example/", bytes.NewReader([]byte{}))
			info := &grpc.UnaryServerInfo{FullMethod: "/package.Service/Method"}
			_, err := interceptor.ServerLogsInterceptor(ctx, req, info, tt.handler)
			assert.Equal(t, err, tt.wantErr)

		})
	}
}
