package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	sharedvalidation "carsharing/shared/validation"
	"github.com/go-playground/validator/v10"
)

type ZoneService struct {
	log      *slog.Logger
	validate *validator.Validate

	zoneRepo ZoneRepository
}

func NewZoneService(
	log *slog.Logger,
	validate *validator.Validate,
	zoneRepo ZoneRepository,
) *ZoneService {
	return &ZoneService{
		log:      pkglog.WithComponent(log, "service.ZoneService"),
		validate: validate,
		zoneRepo: zoneRepo,
	}
}

func (s *ZoneService) Create(ctx context.Context, data validation.ZoneCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	zoneType, _ := model.ZoneTypeFromString(data.Type)
	now := time.Now()

	id, err := s.zoneRepo.Insert(ctx, model.Zone{
		Name:            data.Name,
		Type:            zoneType,
		BoundaryGeoJSON: data.BoundaryGeoJSON,
		FeeAdjustment:   data.FeeAdjustment,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		log.Error("repo: inserting zone", pkglog.Err(err))
		return "", err
	}

	return id, nil
}

func (s *ZoneService) Get(ctx context.Context, id string) (model.Zone, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.Zone{}, err
	}

	zone, err := s.zoneRepo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrZoneNotFound) {
			log.Error("repo: finding zone by id", pkglog.Err(err))
		}
		return model.Zone{}, err
	}

	return zone, nil
}

func (s *ZoneService) List(ctx context.Context, filter validation.ZoneFilter) ([]model.Zone, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	zones, err := s.zoneRepo.Find(ctx, zoneFilter(filter))
	if err != nil {
		log.Error("repo: listing zones", pkglog.Err(err))
		return nil, err
	}

	return zones, nil
}

func (s *ZoneService) Update(ctx context.Context, id string, data validation.ZoneUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	update := model.ZoneUpdate{
		Name:            data.Name,
		BoundaryGeoJSON: data.BoundaryGeoJSON,
		FeeAdjustment:   data.FeeAdjustment,
		IsActive:        data.IsActive,
		UpdatedAt:       time.Now(),
	}

	if data.Type != nil {
		zoneType, _ := model.ZoneTypeFromString(*data.Type)
		update.Type = &zoneType
	}

	if err := s.zoneRepo.Update(ctx, id, update); err != nil {
		if !errors.Is(err, model.ErrZoneNotFound) {
			log.Error("repo: updating zone", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *ZoneService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := s.zoneRepo.Delete(ctx, id); err != nil {
		if !errors.Is(err, model.ErrZoneNotFound) {
			log.Error("repo: deleting zone", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func zoneFilter(filter validation.ZoneFilter) model.ZoneFilter {
	repoFilter := model.ZoneFilter{IsActive: filter.IsActive}
	if filter.Type != nil {
		zt, _ := model.ZoneTypeFromString(*filter.Type)
		repoFilter.Type = &zt
	}
	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	repoFilter.Pagination = &sharedmodel.Pagination{Limit: filter.Pagination.Limit, Offset: filter.Pagination.Offset}
	return repoFilter
}
