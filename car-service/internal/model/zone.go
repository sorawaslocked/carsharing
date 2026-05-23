package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type ZoneType string

const (
	ZoneTypeOperating ZoneType = "operating"
	ZoneTypeNoDrop    ZoneType = "no_drop"
	ZoneParkingHub    ZoneType = "parking_hub"
	ZoneTypeSurcharge ZoneType = "surcharge"
)

var validZoneTypes = map[ZoneType]struct{}{
	ZoneTypeOperating: {},
	ZoneTypeNoDrop:    {},
	ZoneParkingHub:    {},
	ZoneTypeSurcharge: {},
}

func ZoneTypeFromString(s string) (ZoneType, bool) {
	zt := ZoneType(s)
	if _, ok := validZoneTypes[zt]; !ok {
		return "", false
	}
	return zt, true
}

func (t ZoneType) String() string {
	return string(t)
}

type Zone struct {
	ID              string
	Name            string
	Type            ZoneType
	BoundaryGeoJSON string
	FeeAdjustment   int32
	IsActive        bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ZoneFilter struct {
	Type     *ZoneType
	IsActive *bool

	Pagination *sharedmodel.Pagination
}

type ZoneUpdate struct {
	Name            *string
	Type            *ZoneType
	BoundaryGeoJSON *string
	FeeAdjustment   *int32
	IsActive        *bool
	UpdatedAt       time.Time
}
