package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type CarMaintenanceService struct {
	presenter CarMaintenancePresenter
	log       *slog.Logger
}

func NewCarMaintenanceService(presenter CarMaintenancePresenter, log *slog.Logger) *CarMaintenanceService {
	return &CarMaintenanceService{
		presenter: presenter,
		log:       pkglog.WithComponent(log, "service.CarMaintenanceService"),
	}
}

func (s *CarMaintenanceService) CreateTemplate(ctx context.Context, data model.CarMaintenanceTemplateCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CreateTemplate"), utils.MetadataFromCtx(ctx))

	id, err := s.presenter.CreateTemplate(ctx, data)
	if err != nil {
		log.Warn("creating maintenance template", pkglog.Err(err))

		return "", err
	}

	return id, nil
}

func (s *CarMaintenanceService) GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetTemplate"), utils.MetadataFromCtx(ctx))

	template, err := s.presenter.GetTemplate(ctx, id)
	if err != nil {
		log.Warn("getting maintenance template", pkglog.Err(err))

		return model.CarMaintenanceTemplate{}, err
	}

	return template, nil
}

func (s *CarMaintenanceService) ListTemplates(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "ListTemplates"), utils.MetadataFromCtx(ctx))

	templates, err := s.presenter.ListTemplates(ctx, filter)
	if err != nil {
		log.Warn("listing maintenance templates", pkglog.Err(err))

		return nil, err
	}

	return templates, nil
}

func (s *CarMaintenanceService) UpdateTemplate(ctx context.Context, id string, data model.CarMaintenanceTemplateUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "UpdateTemplate"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.UpdateTemplate(ctx, id, data); err != nil {
		log.Warn("updating maintenance template", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarMaintenanceService) DeleteTemplate(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "DeleteTemplate"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.DeleteTemplate(ctx, id); err != nil {
		log.Warn("deleting maintenance template", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarMaintenanceService) ListRecords(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "ListRecords"), utils.MetadataFromCtx(ctx))

	records, err := s.presenter.ListRecords(ctx, filter)
	if err != nil {
		log.Warn("listing maintenance records", pkglog.Err(err))

		return nil, err
	}

	return records, nil
}

func (s *CarMaintenanceService) CompleteRecord(ctx context.Context, recordID string, data model.CarMaintenanceRecordComplete) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CompleteRecord"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.CompleteRecord(ctx, recordID, data); err != nil {
		log.Warn("completing maintenance record", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarMaintenanceService) GetReceiptImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetReceiptImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.presenter.GetReceiptImageUploadData(ctx)
	if err != nil {
		log.Warn("getting maintenance receipt image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}
