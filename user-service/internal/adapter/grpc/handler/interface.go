package handler

import (
	"context"

	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"
)

type DocumentAnalyzedSubscriber interface {
	SubscribeStream(userID *string, passed *bool) (<-chan model.DocumentAnalyzedEvent, func())
}

type UserService interface {
	Create(ctx context.Context, data validation.UserCreate) (string, error)
	Get(ctx context.Context, id string) (model.User, error)
	List(ctx context.Context, filter validation.UserFilter) ([]model.User, error)
	Update(ctx context.Context, id string, data validation.UserUpdate) error
	Delete(ctx context.Context, id string) error

	Register(ctx context.Context, data validation.UserCreate) (string, error)
	SignIn(ctx context.Context, creds validation.Credentials) (string, error)

	SendActivationCode(ctx context.Context) error
	CheckActivationCode(ctx context.Context, code string) error

	GetUserProfileImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error)

	CreateDocument(ctx context.Context, data validation.DocumentCreate) (string, error)
	GetDocumentImageUploadData(ctx context.Context, imageType string) (sharedmodel.ImageUploadData, error)
	ListDocuments(ctx context.Context, filter validation.DocumentFilter) ([]model.Document, error)
	CheckDocument(ctx context.Context, docID string, data validation.DocumentUpdate) error
}
