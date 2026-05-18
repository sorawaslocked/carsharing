package handler

import (
	"carsharing/api-gateway/internal/adapter/http/dto"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	checkers []HealthChecker
}

func NewHealthHandler(checkers []HealthChecker) *HealthHandler {
	return &HealthHandler{checkers: checkers}
}

// Health godoc
// @Summary      Gateway health
// @Description  Returns the health status of the API gateway and all upstream services.
// @Tags         health
// @Produce      json
// @Success      200  {object}  dto.HealthResponse
// @Router       /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	var services []dto.ServiceHealthResponse
	overallStatus := "healthy"

	for _, checker := range h.checkers {
		health, err := checker.Health(c)
		if err != nil {
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
