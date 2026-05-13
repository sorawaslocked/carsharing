package dto

import (
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	baseuser "github.com/sorawaslocked/car-rental-protos/gen/base/user"
)

func UserFromProto(u *baseuser.User) model.User {
	return model.User{
		ID:          u.GetId(),
		Email:       u.GetEmail(),
		PhoneNumber: u.PhoneNumber,
		FirstName:   u.GetFirstName(),
		LastName:    u.GetLastName(),
		BirthDate:   u.GetBirthDate(),
		Password: model.Password{
			Hash: u.GetPasswordHash(),
		},
		ProfileImageURL:    u.ProfileImageUrl,
		Roles:              u.GetRoles(),
		IsDocumentVerified: u.GetIsDocumentVerified(),
		IsEmailVerified:    u.GetIsEmailVerified(),
		IsSuspended:        u.GetIsSuspended(),
		CreatedAt:          u.GetCreatedAt().AsTime(),
		UpdatedAt:          u.GetUpdatedAt().AsTime(),
	}
}
