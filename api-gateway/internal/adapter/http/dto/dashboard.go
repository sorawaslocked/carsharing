package dto

type DashboardResponse struct {
	Users       UserStats        `json:"users"`
	Fleet       FleetStats       `json:"fleet"`
	Bookings    BookingStats     `json:"bookings"`
	Trips       TripStats        `json:"trips"`
	Insurance   InsuranceStats   `json:"insurance"`
	Maintenance MaintenanceStats `json:"maintenance"`
}

type UserStats struct {
	Total               int `json:"total"`
	Active              int `json:"active"`
	Suspended           int `json:"suspended"`
	FullyVerified       int `json:"fullyVerified"`
	PendingVerification int `json:"pendingVerification"`
}

type FleetStats struct {
	Total        int `json:"total"`
	Available    int `json:"available"`
	Reserved     int `json:"reserved"`
	InUse        int `json:"inUse"`
	Maintenance  int `json:"maintenance"`
	OutOfService int `json:"outOfService"`
	Retired      int `json:"retired"`
}

type BookingStats struct {
	Active int `json:"active"`
}

type TripStats struct {
	Active int `json:"active"`
}

type InsuranceStats struct {
	Active           int `json:"active"`
	ExpiringIn30Days int `json:"expiringIn30Days"`
}

type MaintenanceStats struct {
	Pending    int `json:"pending"`
	InProgress int `json:"inProgress"`
	Overdue    int `json:"overdue"`
}
