package dto

import (
	"car-rental-car-service/internal/model"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type carModelRow struct {
	ID           string          `db:"id"`
	Brand        string          `db:"brand"`
	Model        string          `db:"model"`
	Year         int16           `db:"year"`
	FuelType     string          `db:"fuel_type"`
	Transmission string          `db:"transmission"`
	BodyType     string          `db:"body_type"`
	Class        string          `db:"class"`
	Seats        int8            `db:"seats"`
	EngineVolume sql.NullFloat64 `db:"engine_volume"`
	RangeKM      int32           `db:"range_km"`
	Features     pq.StringArray  `db:"features"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
}

func (r carModelRow) toDomain() model.CarModel {
	cm := model.CarModel{
		ID:           r.ID,
		Brand:        r.Brand,
		Model:        r.Model,
		Year:         r.Year,
		FuelType:     model.CarFuelType(r.FuelType),
		Transmission: model.CarTransmission(r.Transmission),
		BodyType:     model.CarBodyType(r.BodyType),
		Class:        model.CarClass(r.Class),
		Seats:        r.Seats,
		RangeKM:      r.RangeKM,
		Features:     []string(r.Features),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}

	if r.EngineVolume.Valid {
		cm.EngineVolume = new(float32(r.EngineVolume.Float64))
	}

	return cm
}

func ScanCarModelRow(s scanner) (model.CarModel, error) {
	var r carModelRow

	err := s.Scan(
		&r.ID, &r.Brand, &r.Model, &r.Year,
		&r.FuelType, &r.Transmission, &r.BodyType, &r.Class,
		&r.Seats, &r.EngineVolume, &r.RangeKM, &r.Features,
		&r.CreatedAt, &r.UpdatedAt,
	)

	if err != nil {
		return model.CarModel{}, err
	}

	return r.toDomain(), nil
}

func BuildCarModelWhereClauses(b *ArgsBuilder, f model.CarModelFilter, tableAlias string) []string {
	var clauses []string

	if f.ID != nil {
		clauses = append(clauses, fmt.Sprintf("%s = %s", column(tableAlias, "id"), b.Add(*f.ID)))
	}
	if f.Brand != nil {
		clauses = append(clauses, fmt.Sprintf("%s = %s", column(tableAlias, "brand"), b.Add(*f.Brand)))
	}
	if f.Model != nil {
		clauses = append(clauses, fmt.Sprintf("%s = %s", column(tableAlias, "model"), b.Add(*f.Model)))
	}
	if f.FuelType != nil {
		clauses = append(clauses, fmt.Sprintf("%s = %s", column(tableAlias, "fuel_type"), b.Add(string(*f.FuelType))))
	}
	if f.Transmission != nil {
		clauses = append(clauses, fmt.Sprintf("%s = %s", column(tableAlias, "transmission"), b.Add(string(*f.Transmission))))
	}
	if f.BodyType != nil {
		clauses = append(clauses, fmt.Sprintf("%s = %s", column(tableAlias, "body_type"), b.Add(string(*f.BodyType))))
	}
	if f.Class != nil {
		clauses = append(clauses, fmt.Sprintf("%s = %s", column(tableAlias, "class"), b.Add(string(*f.Class))))
	}
	if f.MinSeats != nil {
		clauses = append(clauses, fmt.Sprintf("%s >= %s", column(tableAlias, "min_seats"), b.Add(*f.MinSeats)))
	}

	return clauses
}

func BuildCarModelSetClauses(u model.CarModelUpdate, b *ArgsBuilder) []string {
	var clauses []string

	if u.Brand != nil {
		clauses = append(clauses, fmt.Sprintf("brand = %s", b.Add(*u.Brand)))
	}
	if u.Model != nil {
		clauses = append(clauses, fmt.Sprintf("model = %s", b.Add(*u.Model)))
	}
	if u.Year != nil {
		clauses = append(clauses, fmt.Sprintf("year = %s", b.Add(*u.Year)))
	}
	if u.FuelType != nil {
		clauses = append(clauses, fmt.Sprintf("fuel_type = %s", b.Add(string(*u.FuelType))))
	}
	if u.Transmission != nil {
		clauses = append(clauses, fmt.Sprintf("transmission = %s", b.Add(string(*u.Transmission))))
	}
	if u.BodyType != nil {
		clauses = append(clauses, fmt.Sprintf("body_type = %s", b.Add(string(*u.BodyType))))
	}
	if u.Class != nil {
		clauses = append(clauses, fmt.Sprintf("class = %s", b.Add(string(*u.Class))))
	}
	if u.Seats != nil {
		clauses = append(clauses, fmt.Sprintf("seats = %s", b.Add(*u.Seats)))
	}
	if u.EngineVolume != nil {
		clauses = append(clauses, fmt.Sprintf("engine_volume = %s", b.Add(*u.EngineVolume)))
	}
	if u.RangeKM != nil {
		clauses = append(clauses, fmt.Sprintf("range_km = %s", b.Add(*u.RangeKM)))
	}
	if u.Features != nil {
		clauses = append(clauses, fmt.Sprintf("features = %s", b.Add(pq.StringArray(u.Features))))
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = %s", b.Add(u.UpdatedAt)))

	return clauses
}
