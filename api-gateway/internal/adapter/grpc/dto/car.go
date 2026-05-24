package dto

import (
	"carsharing/api-gateway/internal/model"
	basecarpb "carsharing/protos/gen/base/car"
)

func CarFromProto(c *basecarpb.Car) model.Car {
	if c == nil {
		return model.Car{}
	}
	car := model.Car{
		ID:               c.GetId(),
		ModelID:          c.GetModelId(),
		VIN:              c.GetVin(),
		LicensePlate:     c.GetLicensePlate(),
		Color:            c.GetColor(),
		YearManufactured: int16(c.GetYearManufactured()),
		MileageKM:        c.GetMileageKm(),
		FuelLevel:        c.FuelLevel,
		BatteryLevel:     c.BatteryLevel,
		Location:         LocationFromProto(c.GetLocation()),
		TelemetryID:      c.GetTelemetryId(),
		ZoneID:           c.ZoneId,
		FuelStatus:       c.GetFuelStatus(),
		Status:           c.GetStatus(),
		IsRetired:        c.GetIsRetired(),
		Notes:            c.Notes,
		ImageURLs:        c.GetImageUrls(),
	}
	if c.GetLastSeenAt() != nil {
		car.LastSeenAt = c.GetLastSeenAt().AsTime()
	}
	if c.GetCreatedAt() != nil {
		car.CreatedAt = c.GetCreatedAt().AsTime()
	}
	if c.GetUpdatedAt() != nil {
		car.UpdatedAt = c.GetUpdatedAt().AsTime()
	}
	return car
}

func CarStatusReadingFromProto(r *basecarpb.CarStatusReading) model.CarStatusReading {
	return model.CarStatusReading{
		ID:         r.GetId(),
		CarID:      r.GetCarId(),
		FromStatus: r.GetFromStatus(),
		ToStatus:   r.GetToStatus(),
		ActorType:  r.GetActorType(),
		ActorID:    r.ActorId,
		Reason:     r.Reason,
		Metadata:   structToMap(r.GetMetadata()),
		RecordedAt: r.GetRecordedAt().AsTime(),
	}
}

func CarTelemetryReadingFromProto(r *basecarpb.CarTelemetryReading) model.CarTelemetryReading {
	reading := model.CarTelemetryReading{
		ID:           r.GetId(),
		CarID:        r.GetCarId(),
		FuelPct:      r.FuelPct,
		FuelRawPct:   r.FuelRawPct,
		BatteryLevel: r.BatteryLevel,
		MileageKM:    r.MileageKm,
		ActorType:    r.GetActorType(),
		ActorID:      r.ActorId,
		Reason:       r.Reason,
		Metadata:     structToMap(r.GetMetadata()),
	}
	if r.GetLocation() != nil {
		loc := LocationFromProto(r.GetLocation())
		reading.Location = &loc
	}
	if r.GetRecordedAt() != nil {
		reading.RecordedAt = r.GetRecordedAt().AsTime()
	}
	return reading
}
