package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/utils"
	"github.com/sorawaslocked/car-rental-car-service/internal/validation"

	"github.com/go-playground/validator/v10"
)

type CarModelService struct {
	carModelRepo  CarModelRepository
	objectStorage ObjectStorage

	validate *validator.Validate
	log      *slog.Logger
}

func NewCarModelService(
	carModelRepo CarModelRepository,
	objectStorage ObjectStorage,
	validate *validator.Validate,
	log *slog.Logger,
) *CarModelService {
	s := &CarModelService{
		carModelRepo:  carModelRepo,
		objectStorage: objectStorage,
		validate:      validate,
	}

	s.log = pkglog.WithComponent(log, "service.CarModelService")

	return s
}

func (s *CarModelService) Create(ctx context.Context, createInput model.CarModelCreateInput) (string, error) {
	const method = "Create"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, createInput)
	if err != nil {
		return "", handleError(logger, err)
	}

	fuelType, _ := model.ParseCarFuelType(createInput.FuelType)
	transmission, _ := model.ParseCarTransmission(createInput.Transmission)
	bodyType, _ := model.ParseCarBodyType(createInput.BodyType)
	class, _ := model.ParseCarClass(createInput.Class)
	now := time.Now()

	cm := model.CarModel{
		Brand:        createInput.Brand,
		Model:        createInput.Model,
		Year:         createInput.Year,
		FuelType:     fuelType,
		Transmission: transmission,
		BodyType:     bodyType,
		Class:        class,
		Seats:        createInput.Seats,
		EngineVolume: createInput.EngineVolume,
		RangeKM:      createInput.RangeKM,
		Features:     createInput.Features,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	id, err := s.carModelRepo.Insert(ctx, cm)
	if err != nil {
		return "", handleError(logger, err)
	}

	return id, nil
}

func (s *CarModelService) Get(ctx context.Context, id string) (model.CarModel, error) {
	const method = "Get"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	carModel, err := s.carModelRepo.FindByID(ctx, id)
	if err != nil {
		return model.CarModel{}, handleError(logger, err)
	}

	return carModel, nil
}

func (s *CarModelService) GetAll(ctx context.Context, filterInput model.CarModelFilterInput) ([]model.CarModel, error) {
	const method = "GetAll"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, filterInput); err != nil {
		return nil, handleError(logger, err)
	}
	filter := carModelFilterFromInput(filterInput)

	carModels, err := s.carModelRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return carModels, nil
}

func (s *CarModelService) Update(ctx context.Context, id string, updateInput model.CarModelUpdateInput) error {
	const method = "Update"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, updateInput); err != nil {
		return handleError(logger, err)
	}

	update := model.CarModelUpdate{
		Brand:        updateInput.Brand,
		Model:        updateInput.Model,
		Year:         updateInput.Year,
		Seats:        updateInput.Seats,
		EngineVolume: updateInput.EngineVolume,
		RangeKM:      updateInput.RangeKM,
		Features:     updateInput.Features,
		ImageKeys:    updateInput.ImageKeys,
		UpdatedAt:    time.Now(),
	}

	if updateInput.FuelType != nil {
		fuelType, _ := model.ParseCarFuelType(*updateInput.FuelType)
		update.FuelType = &fuelType
	}
	if updateInput.Transmission != nil {
		transmission, _ := model.ParseCarTransmission(*updateInput.Transmission)
		update.Transmission = &transmission
	}
	if updateInput.BodyType != nil {
		bodyType, _ := model.ParseCarBodyType(*updateInput.BodyType)
		update.BodyType = &bodyType
	}
	if updateInput.Class != nil {
		class, _ := model.ParseCarClass(*updateInput.Class)
		update.Class = &class
	}

	if err := s.carModelRepo.Update(ctx, id, update); err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarModelService) Delete(ctx context.Context, id string) error {
	const method = "Delete"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := s.carModelRepo.Delete(ctx, id); err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarModelService) GetImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	const method = "GetImageUploadData"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if s.objectStorage == nil {
		logger.Error("object storage not configured")
		return model.ImageUploadData{}, model.ErrInternalServerError
	}

	data, err := s.objectStorage.GetImageUploadData(ctx, storageKeyPrefixCarModel)
	if err != nil {
		return model.ImageUploadData{}, handleError(logger, err)
	}

	return data, nil
}

func (s *CarModelService) GetImageURLs(ctx context.Context, id string) ([]string, error) {
	const method = "GetImageURLs"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if s.objectStorage == nil {
		return nil, nil
	}

	carModel, err := s.carModelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, handleError(logger, err)
	}

	urls := make([]string, 0, len(carModel.ImageKeys))
	for _, k := range carModel.ImageKeys {
		url, err := s.objectStorage.GetPresignedURL(ctx, k)
		if err != nil {
			return nil, handleError(logger, err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}
