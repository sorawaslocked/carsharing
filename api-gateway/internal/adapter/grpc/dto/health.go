package dto

import (
	"carsharing/api-gateway/internal/model"
	svchelpers "github.com/sorawaslocked/car-rental-protos/gen/service"
)

func ServiceHealthFromProto(r *svchelpers.ServiceHealthResponse) model.ServiceHealth {
	deps := make([]model.DependencyHealth, len(r.GetDependencies()))
	for i, d := range r.GetDependencies() {
		deps[i] = DependencyHealthFromProto(d)
	}
	return model.ServiceHealth{
		Name:          r.GetName(),
		Status:        r.GetStatus(),
		Version:       r.GetVersion(),
		Timestamp:     r.GetTimestamp().AsTime(),
		UptimeSeconds: r.GetUptimeSeconds(),
		Dependencies:  deps,
	}
}

func DependencyHealthFromProto(d *svchelpers.DependencyHealth) model.DependencyHealth {
	return model.DependencyHealth{
		Name:      d.GetName(),
		Status:    d.GetStatus(),
		LatencyMS: d.LatencyMs,
		Error:     d.Error,
	}
}
