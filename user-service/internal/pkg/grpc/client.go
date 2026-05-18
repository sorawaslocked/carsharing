package grpc

import (
	"time"

	"carsharing/user-service/internal/adapter/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type DocumentAnalyzerConfig struct {
	Addr             string        `yaml:"addr"                env:"DOCUMENT_ANALYZER_ADDR"                 env-required:"true"`
	MaxReceiveSizeMb int           `yaml:"max_receive_size_mb" env:"DOCUMENT_ANALYZER_MAX_RECEIVE_SIZE_MB" env-default:"4"`
	TimeKeepAlive    time.Duration `yaml:"time_keep_alive"     env:"DOCUMENT_ANALYZER_TIME_KEEP_ALIVE"     env-default:"1m"`
	Timeout          time.Duration `yaml:"timeout"             env:"DOCUMENT_ANALYZER_TIMEOUT"             env-default:"10s"`
}

func NewClientConn(cfg DocumentAnalyzerConfig) (*grpc.ClientConn, error) {
	keepAliveParams := keepalive.ClientParameters{
		Time:                cfg.TimeKeepAlive,
		Timeout:             cfg.Timeout,
		PermitWithoutStream: true,
	}

	maxReceiveSizeBytes := 1024 * 1024 * cfg.MaxReceiveSizeMb

	baseClientInterceptor := interceptor.NewClientBaseInterceptor()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepAliveParams),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxReceiveSizeBytes)),
		grpc.WithUnaryInterceptor(baseClientInterceptor.Unary),
	}

	return grpc.NewClient(cfg.Addr, opts...)
}
