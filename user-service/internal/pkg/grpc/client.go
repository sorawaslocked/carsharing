package grpc

import (
	"time"

	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type (
	Config struct {
		Client Client `yaml:"client" env-required:"true"`
	}

	Client struct {
		DocumentAnalyzerURL string        `yaml:"document_analyzer_url" env:"GRPC_DOCUMENT_ANALYZER_URL" env-required:"true"`
		MaxReceiveSizeMb    int           `yaml:"max_receive_size_mb" env:"GRPC_MAX_RECEIVE_SIZE_MB" env-default:"4"`
		TimeKeepAlive       time.Duration `yaml:"time_keep_alive" env:"GRPC_TIME_KEEP_ALIVE" env-default:"1m"`
		Timeout             time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"10s"`
	}
)

func Connect(target string, clientCfg Client) (*grpc.ClientConn, error) {
	keepAliveParams := keepalive.ClientParameters{
		Time:                clientCfg.TimeKeepAlive,
		Timeout:             clientCfg.Timeout,
		PermitWithoutStream: true,
	}

	maxReceiveSizeBytes := 1024 * 1024 * clientCfg.MaxReceiveSizeMb

	baseClientInterceptor := interceptor.NewClientBaseInterceptor()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepAliveParams),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxReceiveSizeBytes)),
		grpc.WithUnaryInterceptor(baseClientInterceptor.Unary),
	}

	return grpc.NewClient(target, opts...)
}
