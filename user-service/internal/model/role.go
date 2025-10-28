package model

type Role int32

const (
	RoleUser                  Role = 1
	RoleAdmin                 Role = 2
	RoleTechSupport           Role = 3
	RoleFinanceManager        Role = 4
	RoleMaintenanceSpecialist Role = 5
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
