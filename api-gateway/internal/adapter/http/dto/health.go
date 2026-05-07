package dto

import "github.com/sorawaslocked/car-rental-api-gateway/internal/model"

type HealthResponse struct {
	Status   string                  `json:"status"`
	Services []ServiceHealthResponse `json:"services"`
}

type ServiceHealthResponse struct {
	Name          string                     `json:"name"`
	Status        string                     `json:"status"`
	Version       string                     `json:"version,omitempty"`
	UptimeSeconds uint64                     `json:"uptimeSeconds"`
	Dependencies  []DependencyHealthResponse `json:"dependencies,omitempty"`
}

type DependencyHealthResponse struct {
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	LatencyMS *uint32 `json:"latencyMS,omitempty"`
	Error     *string `json:"error,omitempty"`
}

func ServiceHealthFromModel(h model.ServiceHealth) ServiceHealthResponse {
	deps := make([]DependencyHealthResponse, len(h.Dependencies))
	for i, d := range h.Dependencies {
		deps[i] = DependencyHealthResponse{
			Name:      d.Name,
			Status:    d.Status,
			LatencyMS: d.LatencyMS,
			Error:     d.Error,
		}
	}
	return ServiceHealthResponse{
		Name:          h.Name,
		Status:        h.Status,
		Version:       h.Version,
		UptimeSeconds: h.UptimeSeconds,
		Dependencies:  deps,
	}
}
