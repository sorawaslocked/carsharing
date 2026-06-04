package dto

import (
	"carsharing/api-gateway/internal/model"
	basecarpb "carsharing/protos/gen/base/car"
)

func CarModelFromProto(m *basecarpb.CarModel) model.CarModel {
	cm := model.CarModel{
		ID:           m.GetId(),
		Brand:        m.GetBrand(),
		Model:        m.GetModel(),
		Year:         int16(m.GetYear()),
		FuelType:     m.GetFuelType(),
		Transmission: m.GetTransmission(),
		BodyType:     m.GetBodyType(),
		Class:        m.GetClass(),
		Seats:        int8(m.GetSeats()),
		EngineVolume: m.EngineVolume,
		RangeKM:      m.GetRangeKm(),
		Features:     m.GetFeatures(),
		Images:       ImagesFromProto(m.GetImages()),
	}
	if m.GetCreatedAt() != nil {
		cm.CreatedAt = m.GetCreatedAt().AsTime()
	}
	if m.GetUpdatedAt() != nil {
		cm.UpdatedAt = m.GetUpdatedAt().AsTime()
	}
	return cm
}
