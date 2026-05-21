package handler

import (
	"context"

	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"
)

type UserService interface {
	Create(ctx context.Context, data validation.UserCreate) (string, error)
	Get(ctx context.Context, id string) (model.User, error)
	List(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	Update(ctx context.Context, id string, data validation.UserUpdate) error
	Delete(ctx context.Context, id string) error

	Register(ctx context.Context, data validation.UserCreate) (string, error)
	SignIn(ctx context.Context, creds model.Credentials) (string, error)

	SendActivationCode(ctx context.Context) error
	CheckActivationCode(ctx context.Context, code string) error

	GetUserProfileImageUploadData(ctx context.Context) (model.ImageUploadData, error)

	CreateDocument(ctx context.Context, data validation.DocumentCreate) (string, error)
	GetDocumentImageUploadData(ctx context.Context, imageType string) (model.ImageUploadData, error)
	GetProcessedDocumentsForUser(ctx context.Context, userID string) ([]model.Document, error)
	CheckDocument(ctx context.Context, docID string, data validation.DocumentUpdate) error
}
