package service

import (
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/utils"
	"log/slog"
)

const (
	defaultPaginationLimit  int64 = 20
	defaultPaginationOffset int64 = 0
)

func defaultLogger(oldLog *slog.Logger, method string) *slog.Logger {
	return oldLog.With(
		slog.Group("src",
			slog.String("method", method),
		),
	)
}

func loggerWithMetadata(oldLog *slog.Logger, md utils.Metadata) *slog.Logger {
	return oldLog.With(
		slog.Group("metadata",
			slog.String("clientIP", md.ClientIP),
			slog.String("requestID", md.RequestID),
			slog.Group("user",
				slog.String("id", md.UserID),
				slog.Any("roles", md.UserRoles),
				slog.Bool("verified", md.UserVerified),
			),
		),
	)
}

func carModelFilterFromInput(filterInput model.CarModelFilterInput, ignoreNonUnique bool) model.CarModelFilter {
	filter := model.CarModelFilter{
		ID: filterInput.ID,
	}

	if ignoreNonUnique {
		return filter
	}

	if filterInput.FuelType != nil {
		fuelType, _ := model.ParseCarFuelType(*filterInput.FuelType)
		filter.FuelType = &fuelType
	}
	if filterInput.Transmission != nil {
		transmission, _ := model.ParseCarTransmission(*filterInput.Transmission)
		filter.Transmission = &transmission
	}
	if filterInput.BodyType != nil {
		bodyType, _ := model.ParseCarBodyType(*filterInput.BodyType)
		filter.BodyType = &bodyType
	}
	if filterInput.Class != nil {
		class, _ := model.ParseCarClass(*filterInput.Class)
		filter.Class = &class
	}
	if filterInput.PaginationInput.Limit == nil {
		filter.Limit = new(defaultPaginationLimit)
	}
	if filterInput.PaginationInput.Offset == nil {
		filter.Offset = new(defaultPaginationOffset)
	}

	return filter
}

func carFilterFromInput(filterInput model.CarFilterInput, ignoreNonUnique bool) model.CarFilter {
	filter := model.CarFilter{
		ID: filterInput.ID,
	}

	if ignoreNonUnique {
		return filter
	}

	if filterInput.Status != nil {
		status, _ := model.ParseCarStatus(*filterInput.Status)
		filter.Status = &status
	}

	if filterInput.ModelFilter != nil {
		filter.ModelFilter = new(carModelFilterFromInput(*filterInput.ModelFilter, false))
	}

	if filterInput.LocationFilter != nil {
		filter.LocationFilter = &model.LocationFilter{
			Location: model.Location{
				Latitude:  filterInput.LocationFilter.Location.Latitude,
				Longitude: filterInput.LocationFilter.Location.Longitude,
			},
			RadiusKM: filterInput.LocationFilter.RadiusKM,
		}
	}

	if filterInput.PaginationInput.Limit == nil {
		filter.Limit = new(defaultPaginationLimit)
	} else {
		filter.Limit = filterInput.PaginationInput.Limit
	}
	if filterInput.PaginationInput.Offset == nil {
		filter.Offset = new(defaultPaginationOffset)
	} else {
		filter.Offset = filterInput.PaginationInput.Offset
	}

	return filter
}
