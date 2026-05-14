package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	osrm "github.com/gojuno/go.osrm"
	natsgo "github.com/nats-io/nats.go"

	"github.com/sorawaslocked/car-rental-telematics/internal/config"
	"github.com/sorawaslocked/car-rental-telematics/internal/handler"
	internalnats "github.com/sorawaslocked/car-rental-telematics/internal/nats"
	"github.com/sorawaslocked/car-rental-telematics/internal/server"
	"github.com/sorawaslocked/car-rental-telematics/internal/service"
)

func main() {
	cfg := config.Load()

	nc, err := natsgo.Connect(cfg.NATSUrl)
	if err != nil {
		slog.Error("NATS connect failed", "url", cfg.NATSUrl, "error", err)
		os.Exit(1)
	}
	defer nc.Close()
	slog.Info("connected to NATS", "url", cfg.NATSUrl)

	osrmClient := osrm.NewFromURL(cfg.OSRMUrl)
	simSvc := service.NewSimulationService(osrmClient, cfg.OSRMProfile)

	sub := internalnats.NewSubscriber(nc, simSvc, cfg)
	if err := sub.Subscribe(); err != nil {
		slog.Error("NATS subscribe failed", "error", err)
		os.Exit(1)
	}

	grpcServer := server.NewGRPCServer(handler.NewTelematicsHandler(simSvc))

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
