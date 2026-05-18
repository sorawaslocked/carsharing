package service

import (
	"context"
	"log/slog"
	"time"

	"carsharing/car-service/internal/model"
	pkglog "carsharing/car-service/internal/pkg/log"
	"carsharing/car-service/internal/pkg/utils"
	"carsharing/car-service/internal/validation"
	"github.com/go-playground/validator/v10"
)

type ZoneService struct {
	zoneRepo ZoneRepository

	validate *validator.Validate
	log      *slog.Logger
}

func NewZoneService(
	zoneRepo ZoneRepository,
	validate *validator.Validate,
	log *slog.Logger,
) *ZoneService {
	s := &ZoneService{
		zoneRepo: zoneRepo,
		validate: validate,
	}

	s.log = pkglog.WithComponent(log, "service.ZoneService")

	return s
}

func (s *ZoneService) Create(ctx context.Context, createInput model.ZoneCreateInput) (string, error) {
	const method = "Create"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, createInput); err != nil {
		return "", handleError(logger, err)
	}

	zoneType, _ := model.ParseZoneType(createInput.Type)
	now := time.Now()

	zone := model.Zone{
		Name:            createInput.Name,
		Type:            zoneType,
		BoundaryGeoJSON: createInput.BoundaryGeoJSON,
		FeeAdjustment:   createInput.FeeAdjustment,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	id, err := s.zoneRepo.Insert(ctx, zone)
	if err != nil {
		return "", handleError(logger, err)
	}

	return id, nil
}

func (s *ZoneService) Get(ctx context.Context, id string) (model.Zone, error) {
	const method = "Get"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	zone, err := s.zoneRepo.FindByID(ctx, id)
	if err != nil {
		return model.Zone{}, handleError(logger, err)
	}

	return zone, nil
}

func (s *ZoneService) GetAll(ctx context.Context, filterInput model.ZoneFilterInput) ([]model.Zone, error) {
	const method = "GetAll"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, filterInput); err != nil {
		return nil, handleError(logger, err)
	}
	filter := zoneFilterFromInput(filterInput)

	zones, err := s.zoneRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return zones, nil
}

func (s *ZoneService) Update(ctx context.Context, id string, updateInput model.ZoneUpdateInput) error {
	const method = "Update"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, updateInput); err != nil {
		return handleError(logger, err)
	}

	update := model.ZoneUpdate{
		Name:            updateInput.Name,
		BoundaryGeoJSON: updateInput.BoundaryGeoJSON,
		FeeAdjustment:   updateInput.FeeAdjustment,
		IsActive:        updateInput.IsActive,
		UpdatedAt:       time.Now(),
	}

	if updateInput.Type != nil {
		zoneType, _ := model.ParseZoneType(*updateInput.Type)
		update.Type = &zoneType
	}

	if err := s.zoneRepo.Update(ctx, id, update); err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *ZoneService) Delete(ctx context.Context, id string) error {
	const method = "Delete"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := s.zoneRepo.Delete(ctx, id); err != nil {
		return handleError(logger, err)
	}

	return nil
}
