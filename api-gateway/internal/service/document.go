package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

func (s *UserService) CreateDocument(ctx context.Context, objectKey, imageType string) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CreateDocument"), utils.MetadataFromCtx(ctx))
	log.Debug("creating document")

	id, err := s.presenter.CreateDocument(ctx, objectKey, imageType)
	if err != nil {
		log.Warn("creating document", pkglog.Err(err))

		return "", err
	}

	log.Debug("document created", slog.String("id", id))

	return id, nil
}

func (s *UserService) GetUploadDocumentData(ctx context.Context, imageType string) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetUploadDocumentData"), utils.MetadataFromCtx(ctx))
	log.Debug("getting document image upload data")

	data, err := s.presenter.GetDocumentImageUploadData(ctx, imageType)
	if err != nil {
		log.Warn("getting document upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}

func (s *UserService) GetProfileImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetProfileImageUploadData"), utils.MetadataFromCtx(ctx))
	log.Debug("getting profile image upload data")

	data, err := s.presenter.GetProfileImageUploadData(ctx)
	if err != nil {
		log.Warn("getting profile image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}

func (s *UserService) ListDocuments(ctx context.Context, filter model.DocumentFilter) ([]model.Document, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "ListDocuments"), utils.MetadataFromCtx(ctx))
	log.Debug("listing documents")

	docs, err := s.presenter.ListDocuments(ctx, filter)
	if err != nil {
		log.Warn("listing documents", pkglog.Err(err))

		return nil, err
	}

	return docs, nil
}

func (s *UserService) CheckDocument(ctx context.Context, docID, status string, documentError *string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CheckDocument"), utils.MetadataFromCtx(ctx))
	log.Debug("checking document")

	if err := s.presenter.CheckDocument(ctx, docID, status, documentError); err != nil {
		log.Warn("checking document", pkglog.Err(err))

		return err
	}

	return nil
}
