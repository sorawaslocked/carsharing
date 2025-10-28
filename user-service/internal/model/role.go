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

var nameRole = map[string]Role{
	"user":                   RoleUser,
	"admin":                  RoleAdmin,
	"tech_support":           RoleTechSupport,
	"finance_manager":        RoleFinanceManager,
	"maintenance_specialist": RoleMaintenanceSpecialist,
}

func (role Role) String() string {
	return roleName[role]
}

func FromStringToRole(s string) (Role, error) {
	role, ok := nameRole[s]
	if !ok {
		return 0, ErrInvalidRole
	}

	return role, nil
}
