package dto

import (
	"carsharing/api-gateway/internal/model"
	basecarpb "github.com/sorawaslocked/car-rental-protos/gen/base/car"
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
		TelematicsID:     c.GetTelematicsId(),
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

func CarStatusEntryFromProto(r *basecarpb.CarStatusReading) model.CarStatusReading {
	return model.CarStatusReading{
		ID:         r.GetId(),
		CarID:      r.GetCarId(),
		FromStatus: r.GetFromStatus(),
		ToStatus:   r.GetToStatus(),
		ActorType:  r.GetActorType(),
		ActorID:    r.ActorId,
		Reason:     r.Reason,
		Metadata:   structToMap(r.GetMetadata()),
		ChangedAt:  r.GetRecordedAt().AsTime(),
	}
}

func CarFuelReadingFromProto(r *basecarpb.CarFuelReading) model.CarFuelReading {
	return model.CarFuelReading{
		ID:         r.GetId(),
		CarID:      r.GetCarId(),
		FuelPct:    r.GetFuelPct(),
		RawPct:     r.GetRawPct(),
		ActorType:  r.GetActorType(),
		ActorID:    r.ActorId,
		Reason:     r.Reason,
		Metadata:   structToMap(r.GetMetadata()),
		RecordedAt: r.GetRecordedAt().AsTime(),
	}
}

func CarLocationEntryFromProto(r *basecarpb.CarLocationReading) model.CarLocationReading {
	return model.CarLocationReading{
		ID:         r.GetId(),
		CarID:      r.GetCarId(),
		Location:   LocationFromProto(r.GetLocation()),
		ActorType:  r.GetActorType(),
		ActorID:    r.ActorId,
		Reason:     r.Reason,
		Metadata:   structToMap(r.GetMetadata()),
		RecordedAt: r.GetRecordedAt().AsTime(),
	}
}

func CarBatteryReadingFromProto(r *basecarpb.CarBatteryReading) model.CarBatteryReading {
	return model.CarBatteryReading{
		ID:           r.GetId(),
		CarID:        r.GetCarId(),
		BatteryLevel: r.GetBatteryLevel(),
		ActorType:    r.GetActorType(),
		ActorID:      r.ActorId,
		Reason:       r.Reason,
		Metadata:     structToMap(r.GetMetadata()),
		RecordedAt:   r.GetRecordedAt().AsTime(),
	}
}

func CarMileageEntryFromProto(r *basecarpb.CarMileageReading) model.CarMileageReading {
	return model.CarMileageReading{
		ID:         r.GetId(),
		CarID:      r.GetCarId(),
		MileageKM:  r.GetMileageKm(),
		ActorType:  r.GetActorType(),
		ActorID:    r.ActorId,
		Reason:     r.Reason,
		Metadata:   structToMap(r.GetMetadata()),
		RecordedAt: r.GetRecordedAt().AsTime(),
	}
}
