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

type CarModelService struct {
	log      *slog.Logger
	validate *validator.Validate

	carModelRepo  CarModelRepository
	objectStorage ObjectStorage
}

func NewCarModelService(
	log *slog.Logger,
	validate *validator.Validate,
	carModelRepo CarModelRepository,
	objectStorage ObjectStorage,
) *CarModelService {
	return &CarModelService{
		log:           pkglog.WithComponent(log, "service.CarModelService"),
		validate:      validate,
		carModelRepo:  carModelRepo,
		objectStorage: objectStorage,
	}
}

func (s *CarModelService) Create(ctx context.Context, data validation.CarModelCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	fuelType, _ := model.CarFuelTypeFromString(data.FuelType)
	transmission, _ := model.CarTransmissionFromString(data.Transmission)
	bodyType, _ := model.CarBodyTypeFromString(data.BodyType)
	class, _ := model.CarClassFromString(data.Class)
	now := time.Now()

	id, err := s.carModelRepo.Insert(ctx, model.CarModel{
		Brand:        data.Brand,
		Model:        data.Model,
		Year:         data.Year,
		FuelType:     fuelType,
		Transmission: transmission,
		BodyType:     bodyType,
		Class:        class,
		Seats:        data.Seats,
		EngineVolume: data.EngineVolume,
		RangeKM:      data.RangeKM,
		Features:     data.Features,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		log.Error("repo: inserting car model", pkglog.Err(err))
		return "", err
	}

	return id, nil
}

func (s *CarModelService) Get(ctx context.Context, id string) (model.CarModel, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.CarModel{}, err
	}

	carModel, err := s.carModelRepo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrCarModelNotFound) {
			log.Error("repo: finding car model by id", pkglog.Err(err))
		}
		return model.CarModel{}, err
	}

	for i := range carModel.Images {
		url, err := s.objectStorage.GetPresignedURL(ctx, carModel.Images[i].Key)
		if err != nil {
			log.Error("object storage: getting presigned url", pkglog.Err(err))
			return model.CarModel{}, err
		}
		carModel.Images[i].URL = url
	}

	return carModel, nil
}

func (s *CarModelService) List(ctx context.Context, filter validation.CarModelFilter) ([]model.CarModel, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	carModels, err := s.carModelRepo.Find(ctx, carModelFilter(filter))
	if err != nil {
		log.Error("repo: listing car models", pkglog.Err(err))
		return nil, err
	}

	for i := range carModels {
		for j := range carModels[i].Images {
			url, err := s.objectStorage.GetPresignedURL(ctx, carModels[i].Images[j].Key)
			if err != nil {
				log.Warn("object storage: getting presigned url", pkglog.Err(err))
				continue
			}
			carModels[i].Images[j].URL = url
		}
	}

	return carModels, nil
}

func (s *CarModelService) Update(ctx context.Context, id string, data validation.CarModelUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	update := model.CarModelUpdate{
		Brand:        data.Brand,
		Model:        data.Model,
		Year:         data.Year,
		Seats:        data.Seats,
		EngineVolume: data.EngineVolume,
		RangeKM:      data.RangeKM,
		Features:     data.Features,
		ImageKeys:    data.ImageKeys,
		UpdatedAt:    time.Now(),
	}

	if data.FuelType != nil {
		fuelType, _ := model.CarFuelTypeFromString(*data.FuelType)
		update.FuelType = &fuelType
	}
	if data.Transmission != nil {
		transmission, _ := model.CarTransmissionFromString(*data.Transmission)
		update.Transmission = &transmission
	}
	if data.BodyType != nil {
		bodyType, _ := model.CarBodyTypeFromString(*data.BodyType)
		update.BodyType = &bodyType
	}
	if data.Class != nil {
		class, _ := model.CarClassFromString(*data.Class)
		update.Class = &class
	}

	if err := s.carModelRepo.Update(ctx, id, update); err != nil {
		if !errors.Is(err, model.ErrCarModelNotFound) {
			log.Error("repo: updating car model", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarModelService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := s.carModelRepo.Delete(ctx, id); err != nil {
		if !errors.Is(err, model.ErrCarModelNotFound) {
			log.Error("repo: deleting car model", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarModelService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.objectStorage.GetCarModelImageUploadData(ctx)
	if err != nil {
		log.Error("object storage: getting upload data", pkglog.Err(err))
		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}

func carModelFilter(filter validation.CarModelFilter) model.CarModelFilter {
	repoFilter := model.CarModelFilter{
		Brand:    filter.Brand,
		Model:    filter.Model,
		MinSeats: filter.MinSeats,
	}
	if filter.FuelType != nil {
		ft, _ := model.CarFuelTypeFromString(*filter.FuelType)
		repoFilter.FuelType = &ft
	}
	if filter.Transmission != nil {
		tr, _ := model.CarTransmissionFromString(*filter.Transmission)
		repoFilter.Transmission = &tr
	}
	if filter.BodyType != nil {
		bt, _ := model.CarBodyTypeFromString(*filter.BodyType)
		repoFilter.BodyType = &bt
	}
	if filter.Class != nil {
		cl, _ := model.CarClassFromString(*filter.Class)
		repoFilter.Class = &cl
	}
	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	repoFilter.Pagination = &sharedmodel.Pagination{Limit: filter.Pagination.Limit, Offset: filter.Pagination.Offset}
	return repoFilter
}
