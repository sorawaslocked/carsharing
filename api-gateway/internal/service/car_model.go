package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type CarModelService struct {
	presenter CarModelPresenter
	log       *slog.Logger
}

func NewCarModelService(presenter CarModelPresenter, log *slog.Logger) *CarModelService {
	return &CarModelService{
		presenter: presenter,
		log:       pkglog.WithComponent(log, "service.CarModelService"),
	}
}

func (s *CarModelService) Create(ctx context.Context, data model.CarModelCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	id, err := s.presenter.Create(ctx, data)
	if err != nil {
		log.Warn("creating car model", pkglog.Err(err))

		return "", err
	}

	return id, nil
}

func (s *CarModelService) Get(ctx context.Context, id string) (model.CarModel, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	carModel, err := s.presenter.Get(ctx, id)
	if err != nil {
		log.Warn("getting car model", pkglog.Err(err))

		return model.CarModel{}, err
	}

	return carModel, nil
}

func (s *CarModelService) List(ctx context.Context, filter model.CarModelFilter) ([]model.CarModel, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	carModels, err := s.presenter.List(ctx, filter)
	if err != nil {
		log.Warn("listing car models", pkglog.Err(err))

		return nil, err
	}

	return carModels, nil
}

func (s *CarModelService) Update(ctx context.Context, id string, data model.CarModelUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.Update(ctx, id, data); err != nil {
		log.Warn("updating car model", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarModelService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.Delete(ctx, id); err != nil {
		log.Warn("deleting car model", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarModelService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.presenter.GetImageUploadData(ctx)
	if err != nil {
		log.Warn("getting car model image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}
