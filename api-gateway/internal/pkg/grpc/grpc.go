package grpc

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"time"
)

type (
	Config struct {
		Client Client `yaml:"client" env-required:"true"`
	}

	Client struct {
		UserServiceURL    string `yaml:"user_service_url" env:"GRPC_USER_SERVICE_URL" env-required:"true"`
		CarServiceURL     string `yaml:"car_service_url" env:"GRPC_CAR_SERVICE_URL" env-required:"true"`
		BookingServiceURL string `yaml:"booking_service_url" env:"GRPC_BOOKING_SERVICE_URL" env-required:"true"`
		TripServiceURL    string `yaml:"trip_service_url" env:"GRPC_TRIP_SERVICE_URL" env-required:"true"`

		MaxReceiveSizeMb int           `yaml:"max_receive_size_mb" env:"GRPC_MAX_RECEIVE_SIZE_MB" env-default:"4"`
		TimeKeepAlive    time.Duration `yaml:"time_keep_alive" env:"GRPC_TIME_KEEP_ALIVE" env-default:"1m"`
		Timeout          time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"10s"`
	}
)

func Connect(target string, clientCfg Client) (*grpc.ClientConn, error) {
	keepAliveParams := keepalive.ClientParameters{
		Time:                clientCfg.TimeKeepAlive,
		Timeout:             clientCfg.Timeout,
		PermitWithoutStream: true,
	}

	maxReceiveSizeBytes := 1024 * 1024 * clientCfg.MaxReceiveSizeMb

	baseClientInterceptor := interceptor.NewBaseInterceptor()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepAliveParams),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxReceiveSizeBytes)),
		grpc.WithUnaryInterceptor(baseClientInterceptor.Unary),
		grpc.WithStreamInterceptor(baseClientInterceptor.Stream),
	}

	return grpc.NewClient(target, opts...)
}

func PingServer(conn *grpc.ClientConn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	healthClient := grpc_health_v1.NewHealthClient(conn)

	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "",
	})
	if err != nil {
		return err
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return ErrGrpcServerOffline
	}

	return nil
}
