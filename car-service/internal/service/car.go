package service

import (
	"car-rental-car-service/internal/model"
	"car-rental-car-service/internal/validation"
	"context"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
)

type CarService struct {
	carRepo      CarRepository
	carModelRepo CarModelRepository

	validate *validator.Validate
	log      *slog.Logger
}

func NewCarService(
	carRepo CarRepository,
	carModelRepo CarModelRepository,
	log *slog.Logger,
) *CarService {
	s := &CarService{
		carRepo:      carRepo,
		carModelRepo: carModelRepo,
	}

	s.log = log.With(
		slog.Group("src",
			slog.String("component", "CarService"),
		),
	)

	return s
}

func (s *CarService) CreateCar(ctx context.Context, createInput model.CreateCarInput) (string, error) {
	const method = "CreateCar"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return "", handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, createInput)
	if err != nil {
		return "", handleError(logger, err)
	}

	now := time.Now()

	car := model.Car{
		ModelID:          createInput.ModelID,
		VIN:              createInput.VIN,
		LicensePlate:     createInput.LicensePlate,
		Color:            createInput.Color,
		YearManufactured: createInput.YearManufactured,
		MileageKM:        createInput.MileageKM,
		FuelLevel:        createInput.FuelLevel,
		BatteryLevel:     createInput.BatteryLevel,
		Notes:            createInput.Notes,
		Status:           model.CarStatusAvailable,
		LastSeenAt:       now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	id, err := s.carRepo.Insert(ctx, car)
	if err != nil {
		return "", handleError(logger, err)
	}

	return id, nil
}

func (s *CarService) GetCar(ctx context.Context, filterInput model.CarFilterInput) (model.Car, error) {
	const method = "GetCar"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return model.Car{}, handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return model.Car{}, handleError(logger, err)
	}
	filter := carFilterFromInput(filterInput, true)

	car, err := s.carRepo.FindOne(ctx, filter)
	if err != nil {
		return model.Car{}, handleError(logger, err)
	}

	return car, nil
}

func (s *CarService) GetCars(ctx context.Context, filterInput model.CarFilterInput) ([]model.Car, error) {
	const method = "GetCars"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return nil, handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return nil, handleError(logger, err)
	}
	filter := carFilterFromInput(filterInput, false)

	cars, err := s.carRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return cars, nil
}

func (s *CarService) UpdateCar(ctx context.Context, filterInput model.CarFilterInput, updateInput model.UpdateCarInput) error {
	const method = "UpdateCar"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return handleError(logger, err)
	}
	filter := carFilterFromInput(filterInput, true)

	err = validation.ValidateInput(s.validate, updateInput)
	if err != nil {
		return handleError(logger, err)
	}

	now := time.Now()

	update := model.CarUpdate{
		ModelID:      updateInput.ModelID,
		LicensePlate: updateInput.LicensePlate,
		Color:        updateInput.Color,
		Notes:        updateInput.Notes,
		UpdatedAt:    now,
	}

	err = s.carRepo.Update(ctx, filter, update)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarService) UpdateCarStatus(ctx context.Context, filterInput model.CarFilterInput, statusInput model.UpdateCarStatusInput) error {
	const method = "UpdateCarStatus"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return handleError(logger, err)
	}

	if err = validation.ValidateInput(s.validate, filterInput); err != nil {
		return handleError(logger, err)
	}
	filter := carFilterFromInput(filterInput, true)

	if err = validation.ValidateInput(s.validate, statusInput); err != nil {
		return handleError(logger, err)
	}

	current, err := s.carRepo.FindOne(ctx, filter)
	if err != nil {
		return handleError(logger, err)
	}

	err = transitionCarStatus(current.Status, statusInput.Status)
	if err != nil {
		return handleError(logger, err)
	}

	now := time.Now()

	update := model.CarUpdate{
		Status:    new(statusInput.Status),
		UpdatedAt: now,
	}

	err = s.carRepo.Update(ctx, filter, update)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarService) GetAvailableCars(ctx context.Context, filterInput model.CarFilterInput) ([]model.Car, error) {
	const method = "GetAvailableCars"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return nil, handleError(logger, err)
	}

	if err = validation.ValidateInput(s.validate, filterInput); err != nil {
		return nil, handleError(logger, err)
	}
	filter := carFilterFromInput(filterInput, false)

	filter.Status = new(model.CarStatusAvailable)

	cars, err := s.carRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return cars, nil
}

func (s *CarService) DeleteCar(ctx context.Context, filterInput model.CarFilterInput) error {
	const method = "DeleteCar"

	md, err := metadataFromCtx(ctx, method)
	logger := loggerWithMetadata(s.log, md)
	if err != nil {
		return handleError(logger, err)
	}

	err = validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return handleError(logger, err)
	}
	filter := carFilterFromInput(filterInput, true)

	err = s.carRepo.Delete(ctx, filter)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}
