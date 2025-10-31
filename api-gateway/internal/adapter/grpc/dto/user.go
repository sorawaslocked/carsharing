package dto

import (
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	"github.com/sorawaslocked/car-rental-protos/gen/base"
)

func FromProto(user *base.User) model.User {
	return model.User{
		ID:           user.ID,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		BirthDate:    user.BirthDate,
		PasswordHash: user.PasswordHash,
		Roles:        user.Roles,
		CreatedAt:    user.CreatedAt.AsTime(),
		UpdatedAt:    user.UpdatedAt.AsTime(),
		IsActive:     user.IsActive,
		IsConfirmed:  user.IsConfirmed,
	}
}
