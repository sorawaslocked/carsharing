package auth

import (
	"carsharing/car-service/internal/model"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
)

var fleetRoles = []model.Role{model.RoleAdmin, model.RoleFleetManager}

func buildPolicies() map[string]methodPolicy {
	return map[string]methodPolicy{
		// HealthService — public.
		carsvc.HealthService_Health_FullMethodName: {public: true},

		// CarModelService — reads open to any authenticated caller; writes restricted to fleet roles.
		carsvc.CarModelService_CreateCarModel_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarModelService_GetCarModel_FullMethodName:                {},
		carsvc.CarModelService_ListCarModels_FullMethodName:              {},
		carsvc.CarModelService_UpdateCarModel_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarModelService_DeleteCarModel_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarModelService_GetCarModelImageUploadData_FullMethodName: {allowedRoles: fleetRoles},

		// CarService — reads open to any authenticated caller; writes restricted to fleet roles.
		carsvc.CarService_CreateCar_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarService_GetCar_FullMethodName:                {},
		carsvc.CarService_ListCars_FullMethodName:              {},
		carsvc.CarService_UpdateCar_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarService_UpdateCarTelemetry_FullMethodName:    {allowedRoles: fleetRoles},
		carsvc.CarService_UpdateCarStatus_FullMethodName:       {allowedRoles: fleetRoles},
		carsvc.CarService_DeleteCar_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarService_GetCarImageUploadData_FullMethodName: {allowedRoles: fleetRoles},
		carsvc.CarService_GetCarStatusHistory_FullMethodName:   {allowedRoles: fleetRoles},
		carsvc.CarService_GetCarFuelHistory_FullMethodName:     {allowedRoles: fleetRoles},
		carsvc.CarService_GetCarLocationHistory_FullMethodName: {allowedRoles: fleetRoles},
		carsvc.CarService_GetCarBatteryHistory_FullMethodName:  {allowedRoles: fleetRoles},
		carsvc.CarService_GetCarMileageHistory_FullMethodName:  {allowedRoles: fleetRoles},

		// ZoneService — restricted to fleet roles.
		carsvc.ZoneService_CreateZone_FullMethodName: {allowedRoles: fleetRoles},
		carsvc.ZoneService_GetZone_FullMethodName:    {},
		carsvc.ZoneService_ListZones_FullMethodName:  {},
		carsvc.ZoneService_UpdateZone_FullMethodName: {allowedRoles: fleetRoles},
		carsvc.ZoneService_DeleteZone_FullMethodName: {allowedRoles: fleetRoles},

		// CarInsuranceService — restricted to fleet roles.
		carsvc.CarInsuranceService_CreateCarInsurance_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarInsuranceService_GetCarInsurance_FullMethodName:                {},
		carsvc.CarInsuranceService_ListCarInsurances_FullMethodName:              {},
		carsvc.CarInsuranceService_UpdateCarInsurance_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarInsuranceService_DeleteCarInsurance_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarInsuranceService_GetCarInsuranceImageUploadData_FullMethodName: {allowedRoles: fleetRoles},

		// CarMaintenanceService — restricted to fleet roles.
		carsvc.CarMaintenanceService_CreateMaintenanceTemplate_FullMethodName:            {allowedRoles: fleetRoles},
		carsvc.CarMaintenanceService_GetMaintenanceTemplate_FullMethodName:               {allowedRoles: fleetRoles},
		carsvc.CarMaintenanceService_ListMaintenanceTemplates_FullMethodName:             {allowedRoles: fleetRoles},
		carsvc.CarMaintenanceService_UpdateMaintenanceTemplate_FullMethodName:            {allowedRoles: fleetRoles},
		carsvc.CarMaintenanceService_DeleteMaintenanceTemplate_FullMethodName:            {allowedRoles: fleetRoles},
		carsvc.CarMaintenanceService_ListMaintenanceRecords_FullMethodName:               {allowedRoles: fleetRoles},
		carsvc.CarMaintenanceService_CompleteMaintenanceRecord_FullMethodName:            {allowedRoles: fleetRoles},
		carsvc.CarMaintenanceService_GetMaintenanceReceiptImageUploadData_FullMethodName: {allowedRoles: fleetRoles},
	}
}
