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

type CarInsuranceService struct {
	log      *slog.Logger
	validate *validator.Validate

	insuranceRepo CarInsuranceRepository
	objectStorage ObjectStorage
}

func NewCarInsuranceService(
	log *slog.Logger,
	validate *validator.Validate,
	insuranceRepo CarInsuranceRepository,
	objectStorage ObjectStorage,
) *CarInsuranceService {
	return &CarInsuranceService{
		log:           pkglog.WithComponent(log, "service.CarInsuranceService"),
		validate:      validate,
		insuranceRepo: insuranceRepo,
		objectStorage: objectStorage,
	}
}

func (s *CarInsuranceService) Create(ctx context.Context, data validation.CarInsuranceCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	insType, _ := model.InsuranceTypeFromString(data.Type)
	now := time.Now()

	id, err := s.insuranceRepo.Insert(ctx, model.CarInsurance{
		CarID:     data.CarID,
		Type:      insType,
		Provider:  data.Provider,
		PolicyNum: data.PolicyNum,
		StartsAt:  data.StartsAt,
		ExpiresAt: data.ExpiresAt,
		CostTenge: data.CostTenge,
		Status:    model.InsuranceStatusActive,
		Notes:     data.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		log.Error("repo: inserting car insurance", pkglog.Err(err))
		return "", err
	}

	return id, nil
}

func (s *CarInsuranceService) Get(ctx context.Context, id string) (model.CarInsurance, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.CarInsurance{}, err
	}

	insurance, err := s.insuranceRepo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrCarInsuranceNotFound) {
			log.Error("repo: finding car insurance by id", pkglog.Err(err))
		}
		return model.CarInsurance{}, err
	}

	for i := range insurance.Images {
		url, err := s.objectStorage.GetPresignedURL(ctx, insurance.Images[i].Key)
		if err != nil {
			log.Error("object storage: getting presigned url", pkglog.Err(err))
			return model.CarInsurance{}, err
		}
		insurance.Images[i].URL = url
	}

	return insurance, nil
}

func (s *CarInsuranceService) List(ctx context.Context, filter validation.CarInsuranceFilter) ([]model.CarInsurance, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	insurances, err := s.insuranceRepo.Find(ctx, carInsuranceFilter(filter))
	if err != nil {
		log.Error("repo: listing car insurances", pkglog.Err(err))
		return nil, err
	}

	for i := range insurances {
		for j := range insurances[i].Images {
			url, err := s.objectStorage.GetPresignedURL(ctx, insurances[i].Images[j].Key)
			if err != nil {
				log.Warn("object storage: getting presigned url", pkglog.Err(err))
				continue
			}
			insurances[i].Images[j].URL = url
		}
	}

	return insurances, nil
}

func (s *CarInsuranceService) Update(ctx context.Context, id string, data validation.CarInsuranceUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	update := model.CarInsuranceUpdate{
		Provider:  data.Provider,
		PolicyNum: data.PolicyNum,
		StartsAt:  data.StartsAt,
		ExpiresAt: data.ExpiresAt,
		CostTenge: data.CostTenge,
		Notes:     data.Notes,
		ImageKeys: data.ImageKeys,
		UpdatedAt: time.Now(),
	}

	if data.Status != nil {
		status, _ := model.InsuranceStatusFromString(*data.Status)
		update.Status = &status
	}

	if err := s.insuranceRepo.Update(ctx, id, update); err != nil {
		if !errors.Is(err, model.ErrCarInsuranceNotFound) {
			log.Error("repo: updating car insurance", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarInsuranceService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := s.insuranceRepo.Delete(ctx, id); err != nil {
		if !errors.Is(err, model.ErrCarInsuranceNotFound) {
			log.Error("repo: deleting car insurance", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarInsuranceService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.objectStorage.GetInsuranceImageUploadData(ctx)
	if err != nil {
		log.Error("object storage: getting upload data", pkglog.Err(err))
		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}

func carInsuranceFilter(filter validation.CarInsuranceFilter) model.CarInsuranceFilter {
	repoFilter := model.CarInsuranceFilter{
		CarID:              filter.CarID,
		ExpiringWithinDays: filter.ExpiringWithinDays,
	}
	if filter.Type != nil {
		it, _ := model.InsuranceTypeFromString(*filter.Type)
		repoFilter.Type = &it
	}
	if filter.Status != nil {
		st, _ := model.InsuranceStatusFromString(*filter.Status)
		repoFilter.Status = &st
	}
	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	repoFilter.Pagination = &sharedmodel.Pagination{Limit: filter.Pagination.Limit, Offset: filter.Pagination.Offset}
	return repoFilter
}
