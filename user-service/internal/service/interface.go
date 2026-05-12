package service

import (
	"context"

	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

type UserRepository interface {
	Insert(ctx context.Context, user model.User) (string, error)
	FindByID(ctx context.Context, id string) (model.User, error)
	FindOne(ctx context.Context, filter model.UserFilter) (model.User, error)
	Find(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	Update(ctx context.Context, id string, update model.UserRepoUpdate) error
	Delete(ctx context.Context, id string) error
}

type Publisher interface {
	PublishUserCreated(ctx context.Context, id string) error
	PublishUserUpdated(ctx context.Context, id string, isSecurityUpdate bool) error
	PublishUserDeleted(ctx context.Context, id string) error
}

type ActivationCodeStorage interface {
	Save(ctx context.Context, userID string) (string, error)
	Get(ctx context.Context, userID string) ([]byte, error)
}

type Mailer interface {
	SendActivationCode(ctx context.Context, receiver, code string) error
}
