package dto

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"
	baseuserpb "github.com/sorawaslocked/car-rental-protos/gen/base/user"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
)

func FromCreateUserRequest(req *usersvc.CreateUserRequest) (validation.UserCreate, error) {
	birthDate, err := time.Parse("2006-01-02", req.GetBirthDate())
	if err != nil {
		return validation.UserCreate{}, validation.Errors{
			"birth_date": validation.ErrInvalidDateFormat,
		}
	}

	return validation.UserCreate{
		Email:                req.GetEmail(),
		PhoneNumber:          req.PhoneNumber,
		FirstName:            req.GetFirstName(),
		LastName:             req.GetLastName(),
		BirthDate:            birthDate,
		Password:             req.GetPassword(),
		PasswordConfirmation: req.GetPasswordConfirmation(),
	}, nil
}

func FromRegisterRequest(req *usersvc.RegisterRequest) (validation.UserCreate, error) {
	birthDate, err := time.Parse("2006-01-02", req.GetBirthDate())
	if err != nil {
		return validation.UserCreate{}, validation.Errors{
			"birth_date": validation.ErrInvalidDateFormat,
		}
	}

	return validation.UserCreate{
		Email:                req.GetEmail(),
		PhoneNumber:          req.PhoneNumber,
		FirstName:            req.GetFirstName(),
		LastName:             req.GetLastName(),
		BirthDate:            birthDate,
		Password:             req.GetPassword(),
		PasswordConfirmation: req.GetPasswordConfirmation(),
	}, nil
}

func FromListUsersRequest(req *usersvc.ListUsersRequest) model.UserFilter {
	filter := model.UserFilter{
		Email:              req.Email,
		PhoneNumber:        req.PhoneNumber,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		IsDocumentVerified: req.IsDocumentVerified,
		IsEmailVerified:    req.IsEmailVerified,
		IsSuspended:        req.IsSuspended,
	}

	if req.Pagination != nil {
		filter.Pagination = &model.Pagination{
			Limit:  req.Pagination.Limit,
			Offset: req.Pagination.Offset,
		}
	}

	return filter
}

func FromUpdateUserRequest(req *usersvc.UpdateUserRequest) (validation.UserUpdate, error) {
	update := validation.UserUpdate{
		Email:                req.Email,
		PhoneNumber:          req.PhoneNumber,
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
		ProfileImageKey:      req.ProfileImageKey,
		IsDocumentVerified:   req.IsDocumentVerified,
		IsEmailVerified:      req.IsEmailVerified,
		IsSuspended:          req.IsSuspended,
	}

	if req.BirthDate != nil {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return validation.UserUpdate{}, validation.Errors{
				"birth_date": validation.ErrInvalidDateFormat,
			}
		}
		update.BirthDate = &birthDate
	}

	if len(req.Roles) > 0 {
		roles := make([]sharedmodel.Role, len(req.Roles))
		for i, r := range req.Roles {
			role, err := sharedmodel.RoleFromString(r)
			if err != nil {
				return validation.UserUpdate{}, validation.Errors{
					"roles": validation.ErrInvalidRole,
				}
			}
			roles[i] = role
		}
		update.Roles = roles
	}

	return update, nil
}

func FromSignInRequest(req *usersvc.SignInRequest) model.Credentials {
	return model.Credentials{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.GetPassword(),
	}
}

func UserToProto(user model.User) *baseuserpb.User {
	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.String()
	}

	u := &baseuserpb.User{
		Id:                 user.ID,
		Email:              user.Email,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		BirthDate:          user.BirthDate.Format("2006-01-02"),
		PasswordHash:       user.PasswordHash,
		Roles:              roles,
		IsDocumentVerified: user.IsDocumentVerified,
		IsEmailVerified:    user.IsEmailVerified,
		IsSuspended:        user.IsSuspended,
		CreatedAt:          timestamppb.New(user.CreatedAt),
		UpdatedAt:          timestamppb.New(user.UpdatedAt),
	}

	if user.PhoneNumber != nil {
		u.PhoneNumber = user.PhoneNumber
	}
	if user.ProfileImage != nil && user.ProfileImage.URL != "" {
		u.ProfileImageUrl = &user.ProfileImage.URL
	}

	return u
}
