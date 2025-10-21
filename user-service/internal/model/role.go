package model

type Role struct {
	ID   uint32
	Name RoleName
}

type RoleName int

const (
	RoleUser RoleName = iota
	RoleAdmin
	RoleTechSupport
	RoleFinanceManager
	RoleMaintenanceSpecialist
)
