package handler

import (
	"car-rental-car-service/internal/model"
	"context"
)

type CarModelService interface {
	Create(ctx context.Context, createInput model.CarModelCreateInput) (string, error)
	Get(ctx context.Context, filterInput model.CarModelFilterInput) (model.CarModel, error)
	GetAll(ctx context.Context, filterInput model.CarModelFilterInput) ([]model.CarModel, error)
	Update(ctx context.Context, filterInput model.CarModelFilterInput, updateInput model.CarModelUpdateInput) error
	Delete(ctx context.Context, filterInput model.CarModelFilterInput) error
}

type CarService interface {
	Create(ctx context.Context, createInput model.CarCreateInput) (string, error)
	Get(ctx context.Context, filterInput model.CarFilterInput) (model.Car, error)
	GetAll(ctx context.Context, filterInput model.CarFilterInput) ([]model.Car, error)
	GetAvailableCars(ctx context.Context, filterInput model.CarFilterInput) ([]model.Car, error)
	Update(ctx context.Context, filterInput model.CarFilterInput, updateInput model.CarUpdateInput) error
	UpdateCarStatus(ctx context.Context, filterInput model.CarFilterInput, statusInput model.CarStatusUpdateInput) error
	Delete(ctx context.Context, filterInput model.CarFilterInput) error
}
