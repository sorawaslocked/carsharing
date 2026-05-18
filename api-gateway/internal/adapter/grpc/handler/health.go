package handler

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/utils"
	svchelpers "github.com/sorawaslocked/car-rental-protos/gen/service"
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
