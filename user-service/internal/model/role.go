package model

type Role int32

const (
	RoleUser Role = iota
	RoleAdmin
	RoleTechSupport
	RoleFinanceManager
	RoleMaintenanceSpecialist
)

var roleName = map[Role]string{
	RoleUser:                  "user",
	RoleAdmin:                 "admin",
	RoleTechSupport:           "tech_support",
	RoleFinanceManager:        "finance_manager",
	RoleMaintenanceSpecialist: "maintenance_specialist",
}

func (role Role) String() string {
	return roleName[role]
}
