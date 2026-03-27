package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/utils"
	"github.com/sorawaslocked/car-rental-car-service/internal/validation"

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

func (s *CarService) Create(ctx context.Context, createInput model.CarCreateInput) (string, error) {
	const method = "Create"
	logger := defaultLogger(s.log, method)

	md, ok := utils.MetadataFromCtx(ctx)
	if !ok {
		return "", handleError(logger, model.ErrMissingMetadata)
	}
	logger = loggerWithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, createInput)
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

func (s *CarService) Get(ctx context.Context, filterInput model.CarFilterInput) (model.Car, error) {
	const method = "Get"
	logger := defaultLogger(s.log, method)

	md, ok := utils.MetadataFromCtx(ctx)
	if !ok {
		return model.Car{}, handleError(logger, model.ErrMissingMetadata)
	}
	logger = loggerWithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, filterInput)
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

func (s *CarService) GetAll(ctx context.Context, filterInput model.CarFilterInput) ([]model.Car, error) {
	const method = "GetAll"
	logger := defaultLogger(s.log, method)

	md, ok := utils.MetadataFromCtx(ctx)
	if !ok {
		return []model.Car{}, handleError(logger, model.ErrMissingMetadata)
	}
	logger = loggerWithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, filterInput)
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

func (s *CarService) GetAvailableCars(ctx context.Context, filterInput model.CarFilterInput) ([]model.Car, error) {
	const method = "GetAvailableCars"
	logger := defaultLogger(s.log, method)

	md, ok := utils.MetadataFromCtx(ctx)
	if !ok {
		return []model.Car{}, handleError(logger, model.ErrMissingMetadata)
	}
	logger = loggerWithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, filterInput); err != nil {
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

func (s *CarService) Update(ctx context.Context, filterInput model.CarFilterInput, updateInput model.CarUpdateInput) error {
	const method = "Update"
	logger := defaultLogger(s.log, method)

	md, ok := utils.MetadataFromCtx(ctx)
	if !ok {
		return handleError(logger, model.ErrMissingMetadata)
	}
	logger = loggerWithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, filterInput)
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

func (s *CarService) UpdateCarStatus(ctx context.Context, filterInput model.CarFilterInput, statusInput model.CarStatusUpdateInput) error {
	const method = "UpdateCarStatus"
	logger := defaultLogger(s.log, method)

	md, ok := utils.MetadataFromCtx(ctx)
	if !ok {
		return handleError(logger, model.ErrMissingMetadata)
	}
	logger = loggerWithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return handleError(logger, err)
	}
	filter := carFilterFromInput(filterInput, true)

	err = validation.ValidateInput(s.validate, statusInput)
	if err != nil {
		return handleError(logger, err)
	}
	status, _ := model.ParseCarStatus(statusInput.Status)

	current, err := s.carRepo.FindOne(ctx, filter)
	if err != nil {
		return handleError(logger, err)
	}

	err = transitionCarStatus(current.Status, status)
	if err != nil {
		return handleError(logger, err)
	}

	now := time.Now()

	update := model.CarUpdate{
		Status:    &status,
		UpdatedAt: now,
	}

	err = s.carRepo.Update(ctx, filter, update)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarService) Delete(ctx context.Context, filterInput model.CarFilterInput) error {
	const method = "Delete"
	logger := defaultLogger(s.log, method)

	md, ok := utils.MetadataFromCtx(ctx)
	if !ok {
		return handleError(logger, model.ErrMissingMetadata)
	}
	logger = loggerWithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, filterInput)
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
