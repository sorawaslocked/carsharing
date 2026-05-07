package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
)

type HealthHandler struct {
	checkers []HealthChecker
}

func NewHealthHandler(checkers []HealthChecker) *HealthHandler {
	return &HealthHandler{checkers: checkers}
}

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
