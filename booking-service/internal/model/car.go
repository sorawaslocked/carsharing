package model

type CarStatus string

const (
	CarStatusAvailable    CarStatus = "available"
	CarStatusReserved     CarStatus = "reserved"
	CarStatusInUse        CarStatus = "in_use"
	CarStatusMaintenance  CarStatus = "maintenance"
	CarStatusOutOfService CarStatus = "out_of_service"
)

var validCarStatuses = map[CarStatus]struct{}{
	CarStatusAvailable:    {},
	CarStatusReserved:     {},
	CarStatusInUse:        {},
	CarStatusMaintenance:  {},
	CarStatusOutOfService: {},
}

func CarStatusFromString(s string) (CarStatus, bool) {
	cs := CarStatus(s)
	if _, ok := validCarStatuses[cs]; !ok {
		return "", false
	}
	return cs, true
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
