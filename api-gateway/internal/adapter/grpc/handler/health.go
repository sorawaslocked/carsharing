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
	client healthClient
	log    *slog.Logger
}

func NewHealthHandler(client healthClient, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.HealthHandler"),
	}
}

func (h *HealthHandler) Health(ctx context.Context) (model.ServiceHealth, error) {
	logger := pkglog.WithMethod(h.log, "Health")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.Health(ctx, &emptypb.Empty{})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.ServiceHealth{}, dto.FromGrpcErr(err)
	}

	return dto.ServiceHealthFromProto(res), nil
}
