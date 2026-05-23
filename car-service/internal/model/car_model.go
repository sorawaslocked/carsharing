package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarFuelType string

const (
	CarFuelTypePetrol   CarFuelType = "petrol"
	CarFuelTypeDiesel   CarFuelType = "diesel"
	CarFuelTypeElectric CarFuelType = "electric"
	CarFuelTypeHybrid   CarFuelType = "hybrid"
)

var validCarFuelTypes = map[CarFuelType]struct{}{
	CarFuelTypePetrol:   {},
	CarFuelTypeDiesel:   {},
	CarFuelTypeElectric: {},
	CarFuelTypeHybrid:   {},
}

func CarFuelTypeFromString(s string) (CarFuelType, bool) {
	ft := CarFuelType(s)
	if _, ok := validCarFuelTypes[ft]; !ok {
		return "", false
	}
	return ft, true
}

func (t CarFuelType) String() string {
	return string(t)
}

type CarTransmission string

const (
	CarTransmissionManual CarTransmission = "manual"
	CarTransmissionAuto   CarTransmission = "auto"
)

var validCarTransmissions = map[CarTransmission]struct{}{
	CarTransmissionManual: {},
	CarTransmissionAuto:   {},
}

func CarTransmissionFromString(s string) (CarTransmission, bool) {
	tr := CarTransmission(s)
	if _, ok := validCarTransmissions[tr]; !ok {
		return "", false
	}
	return tr, true
}

func (t CarTransmission) String() string {
	return string(t)
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

var validCarBodyTypes = map[CarBodyType]struct{}{
	CarBodyTypeSedan:       {},
	CarBodyTypeHatchback:   {},
	CarBodyTypeSUV:         {},
	CarBodyTypeCrossover:   {},
	CarBodyTypeMinivan:     {},
	CarBodyTypeCoupe:       {},
	CarBodyTypeConvertible: {},
	CarBodyTypePickup:      {},
}

func CarBodyTypeFromString(s string) (CarBodyType, bool) {
	bt := CarBodyType(s)
	if _, ok := validCarBodyTypes[bt]; !ok {
		return "", false
	}
	return bt, true
}

func (t CarBodyType) String() string {
	return string(t)
}

type CarClass string

const (
	CarClassEconomy  CarClass = "economy"
	CarClassCompact  CarClass = "compact"
	CarClassComfort  CarClass = "comfort"
	CarClassBusiness CarClass = "business"
	CarClassLuxury   CarClass = "luxury"
)

var validCarClasses = map[CarClass]struct{}{
	CarClassEconomy:  {},
	CarClassCompact:  {},
	CarClassComfort:  {},
	CarClassBusiness: {},
	CarClassLuxury:   {},
}

func CarClassFromString(s string) (CarClass, bool) {
	c := CarClass(s)
	if _, ok := validCarClasses[c]; !ok {
		return "", false
	}
	return c, true
}

func (c CarClass) String() string {
	return string(c)
}

type CarModel struct {
	ID           string
	Brand        string
	Model        string
	Year         int16
	FuelType     CarFuelType
	Transmission CarTransmission
	BodyType     CarBodyType
	Class        CarClass
	Seats        int8
	EngineVolume *float32
	RangeKM      int32
	Features     []string
	Images       []sharedmodel.Image

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CarModelFilter struct {
	ID           *string
	Brand        *string
	Model        *string
	FuelType     *CarFuelType
	Transmission *CarTransmission
	BodyType     *CarBodyType
	Class        *CarClass
	MinSeats     *int8

	Pagination *sharedmodel.Pagination
}

type CarModelUpdate struct {
	Brand        *string
	Model        *string
	Year         *int16
	FuelType     *CarFuelType
	Transmission *CarTransmission
	BodyType     *CarBodyType
	Class        *CarClass
	Seats        *int8
	EngineVolume *float32
	RangeKM      *int32
	Features     []string
	ImageKeys    []string
	UpdatedAt    time.Time
}
