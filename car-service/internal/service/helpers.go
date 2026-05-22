package service

import (
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
)

func paginationFromInput(p *sharedmodel.Pagination) *sharedmodel.Pagination {
	if p != nil {
		return p
	}
	return &sharedmodel.Pagination{
		Limit:  sharedmodel.DefaultPaginationLimit,
		Offset: sharedmodel.DefaultPaginationOffset,
	}
}

func carModelFilterFromInput(filterInput validation.CarModelFilter) model.CarModelFilter {
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

	filter.Pagination = paginationFromInput(filterInput.Pagination)

	return filter
}

func carFilterFromInput(filterInput validation.CarFilter) model.CarFilter {
	filter := model.CarFilter{}

	if filterInput.Status != nil {
		status, _ := model.ParseCarStatus(*filterInput.Status)
		filter.Status = &status
	}

	if filterInput.ModelFilter != nil {
		mf := carModelFilterFromInput(*filterInput.ModelFilter)
		filter.ModelFilter = &mf
	}

	filter.LocationFilter = filterInput.LocationFilter

	filter.Pagination = paginationFromInput(filterInput.Pagination)

	return filter
}

func zoneFilterFromInput(filterInput validation.ZoneFilter) model.ZoneFilter {
	filter := model.ZoneFilter{
		IsActive: filterInput.IsActive,
	}

	if filterInput.Type != nil {
		zoneType, _ := model.ParseZoneType(*filterInput.Type)
		filter.Type = &zoneType
	}

	filter.Pagination = paginationFromInput(filterInput.Pagination)

	return filter
}

func insuranceFilterFromInput(filterInput validation.CarInsuranceFilter) model.CarInsuranceFilter {
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

	filter.Pagination = paginationFromInput(filterInput.Pagination)

	return filter
}

func maintenanceTemplateFilterFromInput(filterInput validation.CarMaintenanceTemplateFilter) model.CarMaintenanceTemplateFilter {
	filter := model.CarMaintenanceTemplateFilter{
		IsMandatory: filterInput.IsMandatory,
	}

	filter.Pagination = paginationFromInput(filterInput.Pagination)

	return filter
}

func maintenanceRecordFilterFromInput(filterInput validation.CarMaintenanceRecordFilter) model.CarMaintenanceRecordFilter {
	filter := model.CarMaintenanceRecordFilter{
		CarID:      filterInput.CarID,
		TemplateID: filterInput.TemplateID,
	}

	if filterInput.Status != nil {
		status, _ := model.ParseMaintenanceRecordStatus(*filterInput.Status)
		filter.Status = &status
	}

	filter.Pagination = paginationFromInput(filterInput.Pagination)

	return filter
}
