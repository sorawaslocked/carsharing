package grpc

import (
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ClientConfig struct {
	Host string `yaml:"host" env:"GRPC_CLIENT_HOST" env-required:"true"`
	Port int    `yaml:"port" env:"GRPC_CLIENT_PORT" env-required:"true"`

	TLSCertFile string `yaml:"tls_cert_file" env:"GRPC_CLIENT_TLS_CERT_FILE"`
	TLSKeyFile  string `yaml:"tls_key_file" env:"GRPC_CLIENT_TLS_KEY_FILE"`

	MaxMsgSizeMb int `yaml:"max_msg_size_mb" env:"GRPC_CLIENT_MAX_MSG_SIZE" env-default:"32"`
}

func NewClientConn(
	log *slog.Logger,
	cfg ClientConfig,
	unaryInterceptors grpc.DialOption,
	streamInterceptors grpc.DialOption,
) (*grpc.ClientConn, error) {
	log = pkglog.WithMethod(log, "grpc.NewClientConn")

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
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(tlsCreds),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
		unaryInterceptors,
		streamInterceptors,
	)
	if err != nil {
		log.Error("dialing grpc client", pkglog.Err(err), slog.String("addr", addr))

		return nil, ErrFailedConnection
	}

	return conn, nil
}

func PingClient(ctx context.Context, log *slog.Logger, conn *grpc.ClientConn) (string, error) {
	log = pkglog.WithMethod(log, "grpc.PingClient")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	healthClient := grpc_health_v1.NewHealthClient(conn)

	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		log.Error(
			"checking grpc client health",
			pkglog.Err(err),
			slog.String("target", conn.Target()),
			slog.String("status", convertHealthStatus(resp.Status)),
		)

		return "", ErrFailedConnection
	}

	return convertHealthStatus(resp.Status), nil
}
