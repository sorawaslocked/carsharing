package handler

import (
	"log/slog"

	carsvc "carsharing/protos/gen/service/car"
	pkglog "carsharing/shared/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CarMaintenanceStreamHandler struct {
	log                        *slog.Logger
	maintenanceEventSubscriber MaintenanceEventSubscriber

	carsvc.UnimplementedCarMaintenanceStreamServiceServer
}

func NewCarMaintenanceStreamHandler(log *slog.Logger, maintenanceEventSubscriber MaintenanceEventSubscriber) *CarMaintenanceStreamHandler {
	return &CarMaintenanceStreamHandler{
		log:                        pkglog.WithComponent(log, "adapter.grpc.handler.CarMaintenanceStreamHandler"),
		maintenanceEventSubscriber: maintenanceEventSubscriber,
	}
}

func (h *CarMaintenanceStreamHandler) StreamMaintenanceEvents(_ *emptypb.Empty, stream grpc.ServerStreamingServer[carsvc.MaintenanceEvent]) error {
	ctx := stream.Context()

	ch, unsub := h.maintenanceEventSubscriber.SubscribeMaintenanceEvents()
	defer unsub()

	for {
		select {
		case event, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(&carsvc.MaintenanceEvent{
				CarId:      event.CarID,
				TemplateId: event.TemplateID,
				RecordId:   event.RecordID,
				EventType:  event.EventType,
				OccurredAt: timestamppb.New(event.OccurredAt),
			}); err != nil {
				h.log.Error("failed to send maintenance event", slog.String("error", err.Error()))
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
