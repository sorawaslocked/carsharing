package model

type CarFuelType string

const (
	CarFuelTypePetrol   CarFuelType = "petrol"
	CarFuelTypeDiesel   CarFuelType = "diesel"
	CarFuelTypeElectric CarFuelType = "electric"
	CarFuelTypeHybrid   CarFuelType = "hybrid"
)

func ParseCarFuelType(s string) (CarFuelType, bool) {
	switch CarFuelType(s) {
	case CarFuelTypePetrol, CarFuelTypeDiesel, CarFuelTypeElectric, CarFuelTypeHybrid:
		return CarFuelType(s), true
	default:
		return "", false
	}
}

type CarTransmission string

const (
	CarTransmissionManual CarTransmission = "manual"
	CarTransmissionAuto   CarTransmission = "auto"
)

func ParseCarTransmission(s string) (CarTransmission, bool) {
	switch CarTransmission(s) {
	case CarTransmissionManual, CarTransmissionAuto:
		return CarTransmission(s), true
	default:
		return "", false
	}
}

type CarBodyType string

const (
	CarBodyTypeSedan       CarBodyType = "sedan"
	CarBodyTypeHatchback   CarBodyType = "hatchback"
	CarBodyTypeSUV         CarBodyType = "suv"
	CarBodyTypeCrossover   CarBodyType = "crossover"
	CarBodyTypeMinivan     CarBodyType = "minivan"
	CarBodyTypeCoupe       CarBodyType = "coupe"
	CarBodyTypeConvertible CarBodyType = "convertible"
	CarBodyTypePickup      CarBodyType = "pickup"
)

func ParseCarBodyType(s string) (CarBodyType, bool) {
	switch CarBodyType(s) {
	case CarBodyTypeSedan, CarBodyTypeHatchback, CarBodyTypeSUV, CarBodyTypeCrossover, CarBodyTypeMinivan, CarBodyTypeCoupe, CarBodyTypeConvertible, CarBodyTypePickup:
		return CarBodyType(s), true
	default:
		return "", false
	}
}

type CarClass string

const (
	CarClassEconomy  CarClass = "economy"
	CarClassCompact  CarClass = "compact"
	CarClassComfort  CarClass = "comfort"
	CarClassBusiness CarClass = "business"
	CarClassLuxury   CarClass = "luxury"
)

func ParseCarClass(s string) (CarClass, bool) {
	switch CarClass(s) {
	case CarClassEconomy, CarClassCompact, CarClassComfort, CarClassBusiness, CarClassLuxury:
		return CarClass(s), true
	default:
		return "", false
	}
}

type CarStatus string

const (
	CarStatusAvailable    CarStatus = "available"
	CarStatusReserved     CarStatus = "reserved"
	CarStatusInUse        CarStatus = "in_use"
	CarStatusMaintenance  CarStatus = "maintenance"
	CarStatusOutOfService CarStatus = "out_of_service"
)

func ParseCarStatus(s string) (CarStatus, bool) {
	switch CarStatus(s) {
	case CarStatusAvailable, CarStatusReserved, CarStatusInUse, CarStatusMaintenance, CarStatusOutOfService:
		return CarStatus(s), true
	default:
		return "", false
	}
}

type ZoneType string

const (
	ZoneTypeOperating ZoneType = "operating"
	ZoneTypeNoDrop    ZoneType = "no_drop"
	ZoneParkingHub    ZoneType = "parking_hub"
	ZoneTypeSurcharge ZoneType = "surcharge"
)

func ParseZoneType(s string) (ZoneType, bool) {
	switch ZoneType(s) {
	case ZoneTypeOperating, ZoneTypeNoDrop, ZoneParkingHub, ZoneTypeSurcharge:
		return ZoneType(s), true
	default:
		return "", false
	}
}

type InsuranceType string

const (
	InsuranceTypeOSAGO InsuranceType = "osago"
	InsuranceTypeKASKO InsuranceType = "kasko"
)

func ParseInsuranceType(s string) (InsuranceType, bool) {
	switch InsuranceType(s) {
	case InsuranceTypeOSAGO, InsuranceTypeKASKO:
		return InsuranceType(s), true
	default:
		return "", false
	}
}

type InsuranceStatus string

const (
	InsuranceStatusActive    InsuranceStatus = "active"
	InsuranceStatusExpired   InsuranceStatus = "expired"
	InsuranceStatusCancelled InsuranceStatus = "cancelled"
)

func ParseInsuranceStatus(s string) (InsuranceStatus, bool) {
	switch InsuranceStatus(s) {
	case InsuranceStatusActive, InsuranceStatusExpired, InsuranceStatusCancelled:
		return InsuranceStatus(s), true
	default:
		return "", false
	}
}

type MaintenanceRecordStatus string

const (
	MaintenanceRecordStatusPending    MaintenanceRecordStatus = "pending"
	MaintenanceRecordStatusInProgress MaintenanceRecordStatus = "in_progress"
	MaintenanceRecordStatusCompleted  MaintenanceRecordStatus = "completed"
)

func ParseMaintenanceRecordStatus(s string) (MaintenanceRecordStatus, bool) {
	switch MaintenanceRecordStatus(s) {
	case MaintenanceRecordStatusPending, MaintenanceRecordStatusInProgress, MaintenanceRecordStatusCompleted:
		return MaintenanceRecordStatus(s), true
	default:
		return "", false
	}
}

type CarStatusActor string

const (
	CarStatusActorUser       CarStatusActor = "user"
	CarStatusActorOps        CarStatusActor = "ops"
	CarStatusActorSystem     CarStatusActor = "system"
	CarStatusActorScheduler  CarStatusActor = "scheduler"
	CarStatusActorTelematics CarStatusActor = "telemetry"
)
