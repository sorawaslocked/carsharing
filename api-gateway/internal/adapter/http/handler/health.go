package handler

import (
	"log/slog"

	"carsharing/api-gateway/internal/adapter/http/dto"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	checkers []HealthChecker
	log      *slog.Logger
}

func NewHealthHandler(checkers []HealthChecker, log *slog.Logger) *HealthHandler {
	return &HealthHandler{
		checkers: checkers,
		log:      pkglog.WithComponent(log, "http.HealthHandler"),
	}
}

// Health godoc
// @Summary      Gateway health
// @Description  Returns the health status of the API gateway and all upstream services.
// @Tags         health
// @Produce      json
// @Success      200  {object}  dto.HealthResponse
// @Router       /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Health"), utils.MetadataFromCtx(c))

	var services []dto.ServiceHealthResponse
	overallStatus := "healthy"

	for _, checker := range h.checkers {
		health, err := checker.Health(c)
		if err != nil {
			log.Warn("health checker returned error", pkglog.Err(err))

			overallStatus = "degraded"
			continue
		}
		if health.Status != "healthy" {
			overallStatus = "degraded"
		}
		services = append(services, dto.ServiceHealthFromModel(health))
	}

	dto.Ok(c, dto.HealthResponse{
		Status:   overallStatus,
		Services: services,
	})
}
