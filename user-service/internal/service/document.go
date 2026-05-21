package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"
)

func (s *UserService) CreateDocument(ctx context.Context, data validation.DocumentCreate) (string, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CreateDocument"), utils.MetadataFromCtx(ctx))

	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return "", err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	doc := model.Document{
		UserID:    userID,
		ImageType: model.ImageType(data.ImageType),
		Status:    model.DocumentStatusPending,
		Image:     &model.Image{Key: data.ObjectKey},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.docRepo.Insert(ctx, doc)
	if err != nil {
		logger.Error("repo: inserting document", pkglog.Err(err))
		return "", err
	}

	s.documentAnalyzer.Analyze(ctx, id, data.ObjectKey)

	return id, nil
}

func (s *UserService) GetDocumentImageUploadData(ctx context.Context, imageType string) (model.ImageUploadData, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetDocumentImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.objectStorage.GetDocumentImageUploadData(ctx, imageType)
	if err != nil {
		logger.Error("object storage: getting upload data", pkglog.Err(err))
		return model.ImageUploadData{}, err
	}

	return data, nil
}

func (s *UserService) GetProcessedDocumentsForUser(ctx context.Context, userID string) ([]model.Document, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetProcessedDocumentsForUser"), utils.MetadataFromCtx(ctx))

	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			logger.Error("repo: finding user", pkglog.Err(err))
		}
		return nil, err
	}

	pending := model.DocumentStatusPending
	docs, err := s.docRepo.Find(ctx, model.DocumentFilter{
		UserID:        &userID,
		ExcludeStatus: &pending,
		LatestPerType: true,
	})
	if err != nil {
		logger.Error("repo: listing documents", pkglog.Err(err))
		return nil, err
	}

	for i, doc := range docs {
		if doc.Image == nil || doc.Image.Key == "" {
			continue
		}
		url, err := s.objectStorage.GetImageURL(ctx, doc.Image.Key)
		if err != nil {
			logger.Error("object storage: resolving document image url", pkglog.Err(err))
			return nil, err
		}
		docs[i].Image.URL = url
	}

	return docs, nil
}

func (s *UserService) CheckDocument(ctx context.Context, docID string, data validation.DocumentUpdate) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CheckDocument"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	status := model.DocumentStatus(data.Status)

	doc, err := s.docRepo.FindByID(ctx, docID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			logger.Error("repo: finding document", pkglog.Err(err))
		}
		return err
	}

	if err := s.docRepo.Update(ctx, docID, model.DocumentUpdate{
		Status:    &status,
		Error:     data.Error,
		UpdatedAt: time.Now(),
	}); err != nil {
		logger.Error("repo: updating document", pkglog.Err(err))
		return err
	}

	if status == model.DocumentStatusApproved {
		if err := s.checkAndFlagDocumentVerified(ctx, logger, doc.UserID); err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) HandleDocumentAnalyzed(ctx context.Context, event model.DocumentAnalyzedEvent) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "HandleDocumentAnalyzed"), utils.MetadataFromCtx(ctx))

	doc, err := s.docRepo.FindByID(ctx, event.DocumentID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			logger.Error("repo: finding document", pkglog.Err(err))
		}
		return err
	}

	status := model.DocumentStatusProcessed
	var docError *string
	if !event.Passed {
		status = model.DocumentStatusRejected
		if len(event.Defects) > 0 {
			msg := buildDefectsDescription(event.Defects)
			docError = &msg
		}
	}

	if err := s.docRepo.Update(ctx, doc.ID, model.DocumentUpdate{
		Status:    &status,
		Error:     docError,
		UpdatedAt: time.Now(),
	}); err != nil {
		logger.Error("repo: updating document status", pkglog.Err(err))
		return err
	}

	return nil
}

func (s *UserService) checkAndFlagDocumentVerified(ctx context.Context, logger *slog.Logger, userID string) error {
	latestDocs, err := s.docRepo.Find(ctx, model.DocumentFilter{
		UserID:        &userID,
		LatestPerType: true,
	})
	if err != nil {
		logger.Error("repo: finding latest documents per type", pkglog.Err(err))
		return err
	}

	if len(latestDocs) < len(model.AllImageTypes()) {
		return nil
	}

	for _, doc := range latestDocs {
		if doc.Status != model.DocumentStatusApproved {
			return nil
		}
	}

	isDocVerified := true
	if err := s.userRepo.Update(ctx, userID, model.UserUpdate{
		IsDocumentVerified: &isDocVerified,
		UpdatedAt:          time.Now(),
	}); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			logger.Error("repo: flagging user as document verified", pkglog.Err(err))
		}
		return err
	}

	if err := s.publisher.PublishUserUpdated(ctx, userID, false); err != nil {
		logger.Error("nats: publishing user.updated", pkglog.Err(err))
	}

	return nil
}

func buildDefectsDescription(defects []model.Defect) string {
	parts := make([]string, len(defects))
	for i, d := range defects {
		parts[i] = d.Type + ": " + d.Description
	}
	return strings.Join(parts, "; ")
}
