package service

import (
	"context"
	"log/slog"
	"time"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type CarInsuranceService struct {
	insuranceRepo CarInsuranceRepository
	objectStorage ObjectStorage

	validate *validator.Validate
	log      *slog.Logger
}

func NewCarInsuranceService(
	insuranceRepo CarInsuranceRepository,
	objectStorage ObjectStorage,
	validate *validator.Validate,
	log *slog.Logger,
) *CarInsuranceService {
	s := &CarInsuranceService{
		insuranceRepo: insuranceRepo,
		objectStorage: objectStorage,
		validate:      validate,
	}

	s.log = pkglog.WithComponent(log, "service.CarInsuranceService")

	return s
}

func (s *CarInsuranceService) Create(ctx context.Context, createInput validation.CarInsuranceCreate) (string, error) {
	const method = "Create"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, createInput)
	if err != nil {
		return "", handleError(logger, err)
	}

	insType, _ := model.ParseInsuranceType(createInput.Type)
	now := time.Now()

	insurance := model.CarInsurance{
		CarID:     createInput.CarID,
		Type:      insType,
		Provider:  createInput.Provider,
		PolicyNum: createInput.PolicyNum,
		StartsAt:  createInput.StartsAt,
		ExpiresAt: createInput.ExpiresAt,
		CostTenge: createInput.CostTenge,
		Status:    model.InsuranceStatusActive,
		Notes:     createInput.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := s.insuranceRepo.Insert(ctx, insurance)
	if err != nil {
		return "", handleError(logger, err)
	}

	return id, nil
}

func (s *CarInsuranceService) Get(ctx context.Context, id string) (model.CarInsurance, error) {
	const method = "Get"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	insurance, err := s.insuranceRepo.FindByID(ctx, id)
	if err != nil {
		return model.CarInsurance{}, handleError(logger, err)
	}

	for i := range insurance.Images {
		url, err := s.objectStorage.GetPresignedURL(ctx, insurance.Images[i].Key)
		if err != nil {
			return model.CarInsurance{}, handleError(logger, err)
		}
		insurance.Images[i].URL = url
	}

	return insurance, nil
}

func (s *CarInsuranceService) GetAll(ctx context.Context, filterInput validation.CarInsuranceFilter) ([]model.CarInsurance, error) {
	const method = "GetAll"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return nil, handleError(logger, err)
	}
	filter := insuranceFilterFromInput(filterInput)

	insurances, err := s.insuranceRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	for i := range insurances {
		for j := range insurances[i].Images {
			url, err := s.objectStorage.GetPresignedURL(ctx, insurances[i].Images[j].Key)
			if err != nil {
				return nil, handleError(logger, err)
			}
			insurances[i].Images[j].URL = url
		}
	}

	return insurances, nil
}

func (s *CarInsuranceService) Update(ctx context.Context, id string, updateInput validation.CarInsuranceUpdate) error {
	const method = "Update"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, updateInput); err != nil {
		return handleError(logger, err)
	}

	update := model.CarInsuranceUpdate{
		Provider:  updateInput.Provider,
		PolicyNum: updateInput.PolicyNum,
		StartsAt:  updateInput.StartsAt,
		ExpiresAt: updateInput.ExpiresAt,
		CostTenge: updateInput.CostTenge,
		Notes:     updateInput.Notes,
		ImageKeys: updateInput.ImageKeys,
		UpdatedAt: time.Now(),
	}

	if updateInput.Status != nil {
		status, _ := model.ParseInsuranceStatus(*updateInput.Status)
		update.Status = &status
	}

	err := s.insuranceRepo.Update(ctx, id, update)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarInsuranceService) Delete(ctx context.Context, id string) error {
	const method = "Delete"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := s.insuranceRepo.Delete(ctx, id); err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarInsuranceService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	const method = "GetImageUploadData"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	data, err := s.objectStorage.GetInsuranceImageUploadData(ctx)
	if err != nil {
		return sharedmodel.ImageUploadData{}, handleError(logger, err)
	}

	return data, nil
}
