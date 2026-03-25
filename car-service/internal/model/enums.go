package model

type CarFuelType string

const (
	CarFuelTypePetrol   CarFuelType = "petrol"
	CarFuelTypeDiesel   CarFuelType = "diesel"
	CarFuelTypeElectric CarFuelType = "electric"
	CarFuelTypeHybrid   CarFuelType = "hybrid"
)

type CarTransmission string

const (
	CarTransmissionManual CarTransmission = "manual"
	CarTransmissionAuto   CarTransmission = "auto"
)

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

type CarClass string

const (
	CarClassEconomy  CarClass = "economy"
	CarClassCompact  CarClass = "compact"
	CarClassComfort  CarClass = "comfort"
	CarClassBusiness CarClass = "business"
	CarClassLuxury   CarClass = "luxury"
)

type CarStatus string

const (
	CarStatusAvailable    CarStatus = "available"
	CarStatusReserved     CarStatus = "reserved"
	CarStatusInUse        CarStatus = "in_use"
	CarStatusMaintenance  CarStatus = "maintenance"
	CarStatusOutOfService CarStatus = "out_of_service"
)

func ParseCarFuelType(s string) (CarFuelType, bool) {
	switch CarFuelType(s) {
	case CarFuelTypePetrol, CarFuelTypeDiesel, CarFuelTypeElectric, CarFuelTypeHybrid:
		return CarFuelType(s), true
	default:
		return "", false
	}
}

func ParseCarTransmission(s string) (CarTransmission, bool) {
	switch CarTransmission(s) {
	case CarTransmissionManual, CarTransmissionAuto:
		return CarTransmission(s), true
	default:
		return "", false
	}
}

func ParseCarBodyType(s string) (CarBodyType, bool) {
	switch CarBodyType(s) {
	case CarBodyTypeSedan, CarBodyTypeHatchback, CarBodyTypeSUV, CarBodyTypeCrossover, CarBodyTypeMinivan, CarBodyTypeCoupe, CarBodyTypeConvertible, CarBodyTypePickup:
		return CarBodyType(s), true
	default:
		return "", false
	}
}

func ParseCarClass(s string) (CarClass, bool) {
	switch CarClass(s) {
	case CarClassEconomy, CarClassCompact, CarClassComfort, CarClassBusiness, CarClassLuxury:
		return CarClass(s), true
	default:
		return "", false
	}
}

func ParseCarStatus(s string) (CarStatus, bool) {
	switch CarStatus(s) {
	case CarStatusAvailable, CarStatusReserved, CarStatusInUse, CarStatusMaintenance, CarStatusOutOfService:
		return CarStatus(s), true
	default:
		return "", false
	}
}
