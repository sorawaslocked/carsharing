package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
)

func (s *UserService) CreateDocument(ctx context.Context, objectKey, imageType string) (string, error) {
	return s.presenter.CreateDocument(ctx, objectKey, imageType)
}

func (s *UserService) GetUploadDocumentData(ctx context.Context, imageType string) (model.ImageUploadData, error) {
	return s.presenter.GetDocumentImageUploadData(ctx, imageType)
}

func (s *UserService) GetProfileImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.presenter.GetProfileImageUploadData(ctx)
}

func (s *UserService) GetProcessedDocumentsForUser(ctx context.Context, userID string) ([]model.Document, error) {
	return s.presenter.GetProcessedDocumentsForUser(ctx, userID)
}

func (s *UserService) CheckDocument(ctx context.Context, docID, status string, documentError *string) error {
	return s.presenter.CheckDocument(ctx, docID, status, documentError)
}
