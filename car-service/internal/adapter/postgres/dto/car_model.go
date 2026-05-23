package dto

import (
	"fmt"
	"time"

	"carsharing/car-service/internal/model"
)

type carModelRow struct {
	ID           string
	Brand        string
	Model        string
	Year         int16
	FuelType     string
	Transmission string
	BodyType     string
	Class        string
	Seats        int8
	EngineVolume *float32
	RangeKM      int32
	Features     []string
	ImageKeys    []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r carModelRow) toDomain() model.CarModel {
	return model.CarModel{
		ID:           r.ID,
		Brand:        r.Brand,
		Model:        r.Model,
		Year:         r.Year,
		FuelType:     model.CarFuelType(r.FuelType),
		Transmission: model.CarTransmission(r.Transmission),
		BodyType:     model.CarBodyType(r.BodyType),
		Class:        model.CarClass(r.Class),
		Seats:        r.Seats,
		EngineVolume: r.EngineVolume,
		RangeKM:      r.RangeKM,
		Features:     r.Features,
		Images:       ImageKeysToImages(r.ImageKeys),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func ScanCarModelRow(s scanner) (model.CarModel, error) {
	var r carModelRow

	err := s.Scan(
		&r.ID, &r.Brand, &r.Model, &r.Year,
		&r.FuelType, &r.Transmission, &r.BodyType, &r.Class,
		&r.Seats, &r.EngineVolume, &r.RangeKM, &r.Features, &r.ImageKeys,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.CarModel{}, err
	}

	return r.toDomain(), nil
}

func WhereClausesFromCarModelFilter(f model.CarModelFilter, args []any, n int, tableAlias string) ([]string, []any, int) {
	var clauses []string

	if f.ID != nil {
		n++
		args = append(args, *f.ID)
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column(tableAlias, "id"), n))
	}
	if f.Brand != nil {
		n++
		args = append(args, *f.Brand)
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column(tableAlias, "brand"), n))
	}
	if f.Model != nil {
		n++
		args = append(args, *f.Model)
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column(tableAlias, "model"), n))
	}
	if f.FuelType != nil {
		n++
		args = append(args, string(*f.FuelType))
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column(tableAlias, "fuel_type"), n))
	}
	if f.Transmission != nil {
		n++
		args = append(args, string(*f.Transmission))
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column(tableAlias, "transmission"), n))
	}
	if f.BodyType != nil {
		n++
		args = append(args, string(*f.BodyType))
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column(tableAlias, "body_type"), n))
	}
	if f.Class != nil {
		n++
		args = append(args, string(*f.Class))
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column(tableAlias, "class"), n))
	}
	if f.MinSeats != nil {
		n++
		args = append(args, *f.MinSeats)
		clauses = append(clauses, fmt.Sprintf("%s >= $%d", column(tableAlias, "seats"), n))
	}

	return clauses, args, n
}

func SetClausesFromCarModelUpdate(update model.CarModelUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	n := 0

	if update.Brand != nil {
		n++
		args = append(args, *update.Brand)
		clauses = append(clauses, fmt.Sprintf("brand = $%d", n))
	}
	if update.Model != nil {
		n++
		args = append(args, *update.Model)
		clauses = append(clauses, fmt.Sprintf("model = $%d", n))
	}
	if update.Year != nil {
		n++
		args = append(args, *update.Year)
		clauses = append(clauses, fmt.Sprintf("year = $%d", n))
	}
	if update.FuelType != nil {
		n++
		args = append(args, string(*update.FuelType))
		clauses = append(clauses, fmt.Sprintf("fuel_type = $%d", n))
	}
	if update.Transmission != nil {
		n++
		args = append(args, string(*update.Transmission))
		clauses = append(clauses, fmt.Sprintf("transmission = $%d", n))
	}
	if update.BodyType != nil {
		n++
		args = append(args, string(*update.BodyType))
		clauses = append(clauses, fmt.Sprintf("body_type = $%d", n))
	}
	if update.Class != nil {
		n++
		args = append(args, string(*update.Class))
		clauses = append(clauses, fmt.Sprintf("class = $%d", n))
	}
	if update.Seats != nil {
		n++
		args = append(args, *update.Seats)
		clauses = append(clauses, fmt.Sprintf("seats = $%d", n))
	}
	if update.EngineVolume != nil {
		n++
		args = append(args, *update.EngineVolume)
		clauses = append(clauses, fmt.Sprintf("engine_volume = $%d", n))
	}
	if update.RangeKM != nil {
		n++
		args = append(args, *update.RangeKM)
		clauses = append(clauses, fmt.Sprintf("range_km = $%d", n))
	}
	if update.Features != nil {
		n++
		args = append(args, update.Features)
		clauses = append(clauses, fmt.Sprintf("features = $%d", n))
	}
	if update.ImageKeys != nil {
		n++
		args = append(args, update.ImageKeys)
		clauses = append(clauses, fmt.Sprintf("image_keys = $%d", n))
	}

	n++
	args = append(args, update.UpdatedAt)
	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", n))

	return clauses, args, n
}
