package interceptor

import "github.com/sorawaslocked/car-rental-user-service/internal/model"

const (
	UserServiceCreate = "/service.user.UserService/Create"
	UserServiceGet    = "/service.user.UserService/Get"
	UserServiceGetAll = "/service.user.UserService/GetAll"
	UserServiceUpdate = "/service.user.UserService/Update"
	UserServiceDelete = "/service.user.UserService/Delete"
)

func createPermittedRoles() map[string]map[model.Role]bool {
	permittedRoles := make(map[string]map[model.Role]bool)

	permittedRoles[UserServiceCreate] = map[model.Role]bool{
		model.RoleAdmin: true,
	}
	permittedRoles[UserServiceGet] = map[model.Role]bool{
		model.RoleUser:  true,
		model.RoleAdmin: true,
	}
	permittedRoles[UserServiceGetAll] = map[model.Role]bool{
		model.RoleAdmin: true,
	}
	permittedRoles[UserServiceUpdate] = map[model.Role]bool{
		model.RoleUser:  true,
		model.RoleAdmin: true,
	}
	permittedRoles[UserServiceDelete] = map[model.Role]bool{
		model.RoleAdmin: true,
	}

	return permittedRoles
}
