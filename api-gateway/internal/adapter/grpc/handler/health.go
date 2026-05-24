package handler

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	svchelpers "carsharing/protos/gen/service"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type healthClient interface {
	Health(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*svchelpers.ServiceHealthResponse, error)
}

type HealthHandler struct {
	name   string
	client healthClient
	log    *slog.Logger
}

func NewHealthHandler(name string, client healthClient, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		name:   name,
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.HealthHandler"),
	}
}

func (h *HealthHandler) Health(ctx context.Context) (model.ServiceHealth, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Health"), utils.MetadataFromCtx(ctx))

	res, err := h.client.Health(ctx, &emptypb.Empty{})
	if err != nil {
		log.Warn("checking service health", slog.String("service", h.name), pkglog.Err(err))

		return model.ServiceHealth{Name: h.name, Status: "degraded"}, nil
	}

	return dto.ServiceHealthFromProto(res), nil
}
