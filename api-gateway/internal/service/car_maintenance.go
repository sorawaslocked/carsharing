package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
)

type CarMaintenanceService struct {
	presenter CarMaintenancePresenter
}

func NewCarMaintenanceService(presenter CarMaintenancePresenter) *CarMaintenanceService {
	return &CarMaintenanceService{presenter: presenter}
}

func (s *CarMaintenanceService) CreateTemplate(ctx context.Context, data model.CarMaintenanceTemplateCreate) (string, error) {
	return s.presenter.CreateTemplate(ctx, data)
}

func (s *CarMaintenanceService) GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error) {
	return s.presenter.GetTemplate(ctx, id)
}

func (s *CarMaintenanceService) ListTemplates(ctx context.Context, filter model.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error) {
	return s.presenter.ListTemplates(ctx, filter)
}

func (s *CarMaintenanceService) UpdateTemplate(ctx context.Context, id string, data model.CarMaintenanceTemplateUpdate) error {
	return s.presenter.UpdateTemplate(ctx, id, data)
}

func (s *CarMaintenanceService) DeleteTemplate(ctx context.Context, id string) error {
	return s.presenter.DeleteTemplate(ctx, id)
}

func (s *CarMaintenanceService) ListRecords(ctx context.Context, filter model.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error) {
	return s.presenter.ListRecords(ctx, filter)
}

func (s *CarMaintenanceService) CompleteRecord(ctx context.Context, recordID string, data model.CarMaintenanceRecordComplete) error {
	return s.presenter.CompleteRecord(ctx, recordID, data)
}

func (s *CarMaintenanceService) GetReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	return s.presenter.GetReceiptImageUploadData(ctx)
}
