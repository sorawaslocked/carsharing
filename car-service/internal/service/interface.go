package service

import (
	"car-rental-car-service/internal/model"
	"context"
)

type CarModelRepository interface {
	Insert(ctx context.Context, carModel model.CarModel) (string, error)
	FindOne(ctx context.Context, filter model.CarModelFilter) (model.CarModel, error)
	Find(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error)
	Update(ctx context.Context, filter model.CarModelFilter, update model.CarModelUpdate) error
	Delete(ctx context.Context, filter model.CarModelFilter) error
}

type CarRepository interface {
	Insert(ctx context.Context, car model.Car) (string, error)
	FindOne(ctx context.Context, filter model.CarFilter) (model.Car, error)
	Find(ctx context.Context, filter model.CarFilter) ([]model.Car, error)
	Update(ctx context.Context, filter model.CarFilter, update model.CarUpdate) error
	Delete(ctx context.Context, filter model.CarFilter) error
}

type TelematicsRepository interface {
	InsertEvent(ctx context.Context, event model.CarTelematicsEvent) error
}

type TelematicsQueue interface {
	Pop(ctx context.Context) (model.TelematicsUpdate, AckFunc, error)
}

type AckFunc func(err error) error
