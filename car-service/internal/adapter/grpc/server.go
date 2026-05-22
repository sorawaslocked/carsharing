package grpc

import (
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/handler"
	"carsharing/car-service/internal/adapter/grpc/interceptor"
	authinterceptor "carsharing/car-service/internal/adapter/grpc/interceptor/auth"
	pkggrpc "carsharing/shared/pkg/grpc"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
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
	telematicsSubscriber handler.TelematicsSubscriber,
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
		grpc.ChainStreamInterceptor(),
	)
	if err != nil {
		return nil, err
	}

	carsvc.RegisterCarModelServiceServer(s, handler.NewCarModelHandler(carModelService, log))
	carsvc.RegisterCarServiceServer(s, handler.NewCarHandler(carService, log))
	carsvc.RegisterCarInsuranceServiceServer(s, handler.NewCarInsuranceHandler(carInsuranceService, log))
	carsvc.RegisterCarMaintenanceServiceServer(s, handler.NewCarMaintenanceHandler(carMaintenanceService, log))
	carsvc.RegisterZoneServiceServer(s, handler.NewZoneHandler(zoneService, log))
	carsvc.RegisterCarStreamServiceServer(s, handler.NewCarStreamHandler(carService, telematicsSubscriber, log))
	carsvc.RegisterHealthServiceServer(s, healthHandler)

	return s, nil
}
