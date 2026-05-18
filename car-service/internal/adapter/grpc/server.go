package grpc

import (
	"log/slog"

	"carsharing/car-service/internal/adapter/grpc/handler"
	"carsharing/car-service/internal/adapter/grpc/interceptor"
	authinterceptor "carsharing/car-service/internal/adapter/grpc/interceptor/auth"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	s   *grpc.Server
	log *slog.Logger
}

func NewServer(
	log *slog.Logger,
	carModelService handler.CarModelService,
	carService handler.CarService,
	carInsuranceService handler.CarInsuranceService,
	carMaintenanceService handler.CarMaintenanceService,
	zoneService handler.ZoneService,
	telematicsSubscriber handler.TelematicsSubscriber,
	dbPinger handler.DBPinger,
	natsChecker handler.NATSChecker,
) *Server {
	server := &Server{log: log}

	server.register(
		carModelService, carService, carInsuranceService, carMaintenanceService, zoneService,
		telematicsSubscriber, dbPinger, natsChecker,
		log,
	)

	return server
}

func (s *Server) GRPCServer() *grpc.Server {
	return s.s
}

func (s *Server) register(
	carModelService handler.CarModelService,
	carService handler.CarService,
	carInsuranceService handler.CarInsuranceService,
	carMaintenanceService handler.CarMaintenanceService,
	zoneService handler.ZoneService,
	telematicsSubscriber handler.TelematicsSubscriber,
	dbPinger handler.DBPinger,
	natsChecker handler.NATSChecker,
	log *slog.Logger,
) {
	baseInterceptor := interceptor.NewBaseInterceptor()
	authInterceptor := authinterceptor.NewAuthInterceptor(log)

	s.s = grpc.NewServer(grpc.ChainUnaryInterceptor(
		baseInterceptor.Unary,
		authInterceptor.Unary,
	))

	carsvc.RegisterCarModelServiceServer(s.s, handler.NewCarModelHandler(carModelService, log))
	carsvc.RegisterCarServiceServer(s.s, handler.NewCarHandler(carService, log))
	carsvc.RegisterCarInsuranceServiceServer(s.s, handler.NewCarInsuranceHandler(carInsuranceService, log))
	carsvc.RegisterCarMaintenanceServiceServer(s.s, handler.NewCarMaintenanceHandler(carMaintenanceService, log))
	carsvc.RegisterZoneServiceServer(s.s, handler.NewZoneHandler(zoneService, log))
	carsvc.RegisterCarStreamServiceServer(s.s, handler.NewCarStreamHandler(carService, telematicsSubscriber, log))
	carsvc.RegisterHealthServiceServer(s.s, handler.NewHealthHandler(dbPinger, natsChecker, log))

	reflection.Register(s.s)
}
