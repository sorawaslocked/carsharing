package model

type Role string

const (
	RoleUser           Role = "user"
	RoleAdmin          Role = "admin"
	RoleFleetManager   Role = "fleet_manager"
	RoleUserManager    Role = "user_manager"
	RoleBookingManager Role = "booking_manager"
)

var validRoles = map[Role]struct{}{
	RoleUser:           {},
	RoleAdmin:          {},
	RoleFleetManager:   {},
	RoleUserManager:    {},
	RoleBookingManager: {},
}

func RoleFromString(s string) (Role, error) {
	r := Role(s)
	if _, ok := validRoles[r]; !ok {
		return "", ErrInvalidRole
	}
	return r, nil
}

func (r Role) String() string {
	return string(r)
}
