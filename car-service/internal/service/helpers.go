package service

import (
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
)

const (
	defaultPaginationLimit  int64 = 20
	defaultPaginationOffset int64 = 0
)

func paginationFromInput(input model.PaginationInput) model.Pagination {
	defLimit := defaultPaginationLimit
	defOffset := defaultPaginationOffset

	p := model.Pagination{}

	if input.Limit != nil {
		p.Limit = input.Limit
	} else {
		p.Limit = &defLimit
	}

	if input.Offset != nil {
		p.Offset = input.Offset
	} else {
		p.Offset = &defOffset
	}

	return p
}

func carModelFilterFromInput(filterInput model.CarModelFilterInput) model.CarModelFilter {
	filter := model.CarModelFilter{
		Brand:    filterInput.Brand,
		Model:    filterInput.Model,
		MinSeats: filterInput.MinSeats,
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

	filter.Pagination = paginationFromInput(filterInput.PaginationInput)

	return filter
}

func carFilterFromInput(filterInput model.CarFilterInput) model.CarFilter {
	filter := model.CarFilter{}

	if filterInput.Status != nil {
		status, _ := model.ParseCarStatus(*filterInput.Status)
		filter.Status = &status
	}

	if filterInput.ModelFilter != nil {
		mf := carModelFilterFromInput(*filterInput.ModelFilter)
		filter.ModelFilter = &mf
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

	filter.Pagination = paginationFromInput(filterInput.PaginationInput)

	return filter
}

func zoneFilterFromInput(filterInput model.ZoneFilterInput) model.ZoneFilter {
	filter := model.ZoneFilter{
		IsActive: filterInput.IsActive,
	}

	if filterInput.Type != nil {
		zoneType, _ := model.ParseZoneType(*filterInput.Type)
		filter.Type = &zoneType
	}

	filter.Pagination = paginationFromInput(filterInput.PaginationInput)

	return filter
}

func insuranceFilterFromInput(filterInput model.CarInsuranceFilterInput) model.CarInsuranceFilter {
	filter := model.CarInsuranceFilter{
		CarID:              filterInput.CarID,
		ExpiringWithinDays: filterInput.ExpiringWithinDays,
	}

	if filterInput.Type != nil {
		insType, _ := model.ParseInsuranceType(*filterInput.Type)
		filter.Type = &insType
	}

	if filterInput.Status != nil {
		status, _ := model.ParseInsuranceStatus(*filterInput.Status)
		filter.Status = &status
	}

	filter.Pagination = paginationFromInput(filterInput.PaginationInput)

	return filter
}

func maintenanceTemplateFilterFromInput(filterInput model.CarMaintenanceTemplateFilterInput) model.CarMaintenanceTemplateFilter {
	filter := model.CarMaintenanceTemplateFilter{
		IsMandatory: filterInput.IsMandatory,
	}

	filter.Pagination = paginationFromInput(filterInput.PaginationInput)

	return filter
}

func maintenanceRecordFilterFromInput(filterInput model.CarMaintenanceRecordFilterInput) model.CarMaintenanceRecordFilter {
	filter := model.CarMaintenanceRecordFilter{
		CarID:      filterInput.CarID,
		TemplateID: filterInput.TemplateID,
	}

	if filterInput.Status != nil {
		status, _ := model.ParseMaintenanceRecordStatus(*filterInput.Status)
		filter.Status = &status
	}

	filter.Pagination = paginationFromInput(filterInput.PaginationInput)

	return filter
}
