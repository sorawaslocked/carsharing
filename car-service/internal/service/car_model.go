package service

import (
	"car-rental-car-service/internal/model"
	"car-rental-car-service/internal/validation"
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
)

type CarModelService struct {
	carModelRepo CarModelRepository

	validate *validator.Validate
	log      *slog.Logger
}

func NewCarModelService(
	carModelRepo CarModelRepository,
	log *slog.Logger,
) *CarModelService {
	s := &CarModelService{
		carModelRepo: carModelRepo,
	}

	s.log = log.With(
		slog.Group("src",
			slog.String("component", "CarModelService"),
		),
	)

	return s
}

func (s *CarModelService) Create(ctx context.Context, createInput model.CarModelCreateInput) (string, error) {
	const method = "Create"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return "", handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, createInput)
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

func (s *CarModelService) Get(ctx context.Context, filterInput model.CarModelFilterInput) (model.CarModel, error) {
	const method = "Get"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return model.CarModel{}, handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return model.CarModel{}, handleError(logger, err)
	}
	filter := carModelFilterFromInput(filterInput, true)

	carModel, err := s.carModelRepo.FindOne(ctx, filter)
	if err != nil {
		return model.CarModel{}, handleError(logger, err)
	}

	return carModel, nil
}

func (s *CarModelService) GetAll(ctx context.Context, filterInput model.CarModelFilterInput) ([]model.CarModel, error) {
	const method = "GetAll"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return nil, handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return nil, handleError(logger, err)
	}
	filter := carModelFilterFromInput(filterInput, false)

	carModels, err := s.carModelRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return carModels, nil
}

func (s *CarModelService) Update(ctx context.Context, filterInput model.CarModelFilterInput, updateInput model.CarModelUpdateInput) error {
	const method = "Update"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return handleError(logger, err)
	}
	filter := carModelFilterFromInput(filterInput, true)

	err = validation.ValidateInput(s.validate, updateInput)
	if err != nil {
		return handleError(logger, err)
	}

	now := time.Now()
	update := model.CarModelUpdate{
		Brand:        updateInput.Brand,
		Model:        updateInput.Model,
		Year:         updateInput.Year,
		Seats:        updateInput.Seats,
		EngineVolume: updateInput.EngineVolume,
		RangeKM:      updateInput.RangeKM,
		Features:     updateInput.Features,
		UpdatedAt:    now,
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

	err = s.carModelRepo.Update(ctx, filter, update)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarModelService) Delete(ctx context.Context, filterInput model.CarModelFilterInput) error {
	const method = "Delete"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return handleError(logger, err)
	}
	filter := carModelFilterFromInput(filterInput, true)

	err = s.carModelRepo.Delete(ctx, filter)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}
