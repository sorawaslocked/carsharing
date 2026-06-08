package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	osrm "github.com/gojuno/go.osrm"

	"carsharing/telematics-service/internal/config"
	"carsharing/telematics-service/internal/db"
	"carsharing/telematics-service/internal/handler"
	natssub "carsharing/telematics-service/internal/nats"
	"carsharing/telematics-service/internal/server"
	"carsharing/telematics-service/internal/service"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", "path", *configPath, "error", err)
		os.Exit(1)
	}

	interval, err := time.ParseDuration(cfg.TelemetryInterval)
	if err != nil {
		slog.Error("invalid telemetry_interval", "value", cfg.TelemetryInterval, "error", err)
		os.Exit(1)
	}

	carRepo, err := db.NewCarRepository(cfg.DBUrl)
	if err != nil {
		slog.Error("failed to connect to database", "url", cfg.DBUrl, "error", err)
		os.Exit(1)
	}
	defer carRepo.Close()
	slog.Info("connected to database")

	osrmClient := osrm.NewFromURL(cfg.OSRMUrl)
	simSvc := service.NewSimulationService(osrmClient, cfg.OSRMProfile, interval, cfg.SpeedKmh, cfg.FuelConsumPerKm, cfg.BatteryConsumPerKm)

	tripSub, err := natssub.NewTripSubscriber(cfg.NATSUrl, cfg.TripStartedSubject, cfg.TripEndedSubject, cfg.TripCancelledSubject, simSvc)
	if err != nil {
		slog.Error("failed to connect to NATS", "url", cfg.NATSUrl, "error", err)
		os.Exit(1)
	}
	defer tripSub.Close()

	if err := tripSub.Subscribe(); err != nil {
		slog.Error("failed to subscribe to trip events", "error", err)
		os.Exit(1)
	}

	grpcServer := server.NewGRPCServer(handler.NewTelematicsHandler(simSvc, carRepo))

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		slog.Error("listen failed", "addr", cfg.GRPCPort, "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		slog.Info("shutting down gracefully")
		grpcServer.GracefulStop()
	}()

	slog.Info("gRPC server listening", "addr", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("gRPC server error", "error", err)
		os.Exit(1)
	}
}
