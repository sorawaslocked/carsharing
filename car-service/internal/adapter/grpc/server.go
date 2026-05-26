package grpc

import (
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/handler"
	"carsharing/car-service/internal/adapter/grpc/interceptor"
	authinterceptor "carsharing/car-service/internal/adapter/grpc/interceptor/auth"
	carsvc "carsharing/protos/gen/service/car"
	pkggrpc "carsharing/shared/pkg/grpc"
	"google.golang.org/grpc"
)

func NewServer(
	log *slog.Logger,
	cfg pkggrpc.ServerConfig,
	carModelService handler.CarModelService,
	carService handler.CarService,
	carInsuranceService handler.CarInsuranceService,
	carMaintenanceService handler.CarMaintenanceService,
	zoneService handler.ZoneService,
	telemetrySubscriber handler.TelemetrySubscriber,
	statusSubscriber handler.StatusSubscriber,
	healthHandler *handler.HealthHandler,
) (*grpc.Server, error) {
	baseInterceptor := interceptor.NewBaseInterceptor()
	loggerInterceptor := interceptor.NewLoggerInterceptor(log)
	authInterceptor := authinterceptor.NewInterceptor(log)

	s, err := pkggrpc.NewServer(
		log,
		cfg,
		grpc.ChainUnaryInterceptor(
			baseInterceptor.Unary,
			loggerInterceptor.Unary,
			authInterceptor.Unary,
		),
		grpc.ChainStreamInterceptor(
			baseInterceptor.Stream,
			loggerInterceptor.Stream,
			authInterceptor.Stream,
		),
	)
	if err != nil {
		return nil, err
	}

	carsvc.RegisterCarModelServiceServer(s, handler.NewCarModelHandler(log, carModelService))
	carsvc.RegisterCarServiceServer(s, handler.NewCarHandler(log, carService))
	carsvc.RegisterCarInsuranceServiceServer(s, handler.NewCarInsuranceHandler(log, carInsuranceService))
	carsvc.RegisterCarMaintenanceServiceServer(s, handler.NewCarMaintenanceHandler(log, carMaintenanceService))
	carsvc.RegisterZoneServiceServer(s, handler.NewZoneHandler(log, zoneService))
	carsvc.RegisterCarStreamServiceServer(s, handler.NewCarStreamHandler(log, carService, telemetrySubscriber, statusSubscriber))
	carsvc.RegisterHealthServiceServer(s, healthHandler)

	return s, nil
}
