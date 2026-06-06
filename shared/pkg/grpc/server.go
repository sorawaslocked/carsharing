package grpc

import (
	pkglog "carsharing/shared/pkg/log"
	"crypto/tls"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type ServerConfig struct {
	Host string `yaml:"host" env:"GRPC_SERVER_HOST" env-required:"true"`
	Port int    `yaml:"port" env:"GRPC_SERVER_PORT" env-required:"true"`

	TLSCertFile string `yaml:"tls_cert_file" env:"GRPC_SERVER_TLS_CERT_FILE"`
	TLSKeyFile  string `yaml:"tls_key_file" env:"GRPC_SERVER_TLS_KEY_FILE"`

	MaxMsgSizeMb int `yaml:"max_msg_size_mb" env:"GRPC_SERVER_MAX_MSG_SIZE_MB" env-default:"32"`

	MaxConnectionIdle  time.Duration `yaml:"max_connection_idle" env:"GRPC_SERVER_MAX_CONNECTION_IDLE" env-default:"15m"`
	MaxConnectionAge   time.Duration `yaml:"max_connection_age" env:"GRPC_SERVER_MAX_CONNECTION_AGE" env-default:"30m"`
	MaxConnectionGrace time.Duration `yaml:"max_connection_grace" env:"GRPC_SERVER_MAX_CONNECTION_GRACE" env-default:"30s"`
}

func NewServer(
	log *slog.Logger,
	cfg ServerConfig,
	unaryInterceptors grpc.ServerOption,
	streamInterceptors grpc.ServerOption,
) (*grpc.Server, error) {
	log = pkglog.WithMethod(log, "grpc.NewServer")

	var tlsCreds credentials.TransportCredentials

	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			log.Error(
				"creating tls certificate",
				pkglog.Err(err),
				slog.String("tlsCertFile", cfg.TLSCertFile),
				slog.String("tlsKeyFile", cfg.TLSKeyFile),
			)

			return nil, ErrInvalidTLSFiles
		}

		tlsCreds = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			MinVersion:   tls.VersionTLS12,
		})
	} else {
		tlsCreds = insecure.NewCredentials()
	}

	maxMsgSize := cfg.MaxMsgSizeMb * 1024 * 1024

	kaParams := keepalive.ServerParameters{
		MaxConnectionIdle:     cfg.MaxConnectionIdle,
		MaxConnectionAge:      cfg.MaxConnectionAge,
		MaxConnectionAgeGrace: cfg.MaxConnectionGrace,
	}

	unknownHandler := grpc.UnknownServiceHandler(
		func(srv any, stream grpc.ServerStream) error {
			return status.Errorf(codes.Unimplemented, "unknown service or method")
		},
	)

	s := grpc.NewServer(
		grpc.Creds(tlsCreds),
		unaryInterceptors,
		streamInterceptors,
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
		grpc.KeepaliveParams(kaParams),
		unknownHandler,
	)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthSrv)
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(s)

	return s, nil
}
