package metricsmw

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestGrpcMiddleware_ServerMetricsInterceptor(t *testing.T) {
	//не понятно, как получить фактическое значение метрики
	metr, _ := metrics.NewGrpcMetrics("test_grpc")

	tests := []struct {
		name    string
		wantErr error
		handler func(ctx context.Context, req interface{}) (interface{}, error)
	}{
		{
			name:    "TestServerMetricsInterceptor_Success",
			wantErr: nil,
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "test response", nil
			},
		},
		{

			name:    "TestServerMetricsInterceptor_Fail",
			wantErr: errors.New("error"),
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, errors.New("error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			interceptor := NewGrpcMw(*metr)

			ctx := context.Background()
			req := httptest.NewRequest(http.MethodGet, "http://example/", bytes.NewReader([]byte{}))
			info := &grpc.UnaryServerInfo{FullMethod: "/package.Service/Method"}
			_, err := interceptor.ServerMetricsInterceptor(ctx, req, info, tt.handler)
			assert.Equal(t, err, tt.wantErr)

		})
	}
}

func Test_mapStatusCodes(t *testing.T) {
	type args struct {
		Err string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test_Case1",
			args: args{
				Err: "user not found",
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "Test_Case2",
			args: args{
				Err: "wrong password",
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "Test_Case3",
			args: args{
				Err: "error test",
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapStatusCodes(tt.args.Err); got != tt.want {
				t.Errorf("mapStatusCodes() = %v, want %v", got, tt.want)
			}
		})
	}
}
