package dto

import (
	"github.com/sorawaslocked/car-rental-protos/gen/base"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func FromCreateUserRequest(req *usersvc.CreateRequest) (model.UserCreateData, error) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return model.UserCreateData{}, model.ValidationErrors{
			"email": model.ErrInvalidDateFormat,
		}
	}

	data := model.UserCreateData{
		Email:                req.Email,
		PhoneNumber:          *req.PhoneNumber,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		BirthDate:            birthDate,
		IsActive:             &req.IsActive,
		IsConfirmed:          &req.IsConfirmed,
	}

	if len(req.Roles) > 0 {
		roles := make([]model.Role, len(req.Roles))

		for i, roleStr := range req.Roles {
			role, err := model.FromStringToRole(roleStr)
			if err != nil {
				return model.UserCreateData{}, model.ValidationErrors{
					"role": model.ErrInvalidRole,
				}
			}
			roles[i] = role
		}

		data.Roles = &roles
	}

	return data, nil
}

func ToUserProto(user model.User) *base.User {
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.String()
	}

	return &base.User{
		ID:           0,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		BirthDate:    user.BirthDate.Format("2006-01-02"),
		PasswordHash: user.PasswordHash,
		Roles:        roles,
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
		IsActive:     user.IsActive,
		IsConfirmed:  user.IsConfirmed,
	}
}

func FromUpdateUserRequest(req *usersvc.UpdateRequest) (model.UserUpdateData, error) {
	data := model.UserUpdateData{
		Email:                req.NewEmail,
		PhoneNumber:          req.PhoneNumber,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
		IsActive:             req.IsActive,
		IsConfirmed:          req.IsConfirmed,
	}

	if len(req.Roles) > 0 {
		roles := make([]model.Role, len(req.Roles))

		for i, roleStr := range req.Roles {
			role, err := model.FromStringToRole(roleStr)
			if err != nil {
				return model.UserUpdateData{}, model.ValidationErrors{
					"role": model.ErrInvalidRole,
				}
			}
			roles[i] = role
		}

		data.Roles = &roles
	}

	if req.BirthDate != nil {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return model.UserUpdateData{}, model.ValidationErrors{
				"email": model.ErrInvalidDateFormat,
			}
		}

		data.BirthDate = &birthDate
	}

	return data, nil
}
