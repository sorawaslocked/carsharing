package service

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/model"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

type CarInsuranceService struct {
	presenter CarInsurancePresenter
	log       *slog.Logger
}

func NewCarInsuranceService(presenter CarInsurancePresenter, log *slog.Logger) *CarInsuranceService {
	return &CarInsuranceService{
		presenter: presenter,
		log:       pkglog.WithComponent(log, "service.CarInsuranceService"),
	}
}

func (s *CarInsuranceService) Create(ctx context.Context, data model.CarInsuranceCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Create"), utils.MetadataFromCtx(ctx))

	id, err := s.presenter.Create(ctx, data)
	if err != nil {
		log.Warn("creating car insurance", pkglog.Err(err))

		return "", err
	}

	return id, nil
}

func (s *CarInsuranceService) Get(ctx context.Context, id string) (model.CarInsurance, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Get"), utils.MetadataFromCtx(ctx))

	insurance, err := s.presenter.Get(ctx, id)
	if err != nil {
		log.Warn("getting car insurance", pkglog.Err(err))

		return model.CarInsurance{}, err
	}

	return insurance, nil
}

func (s *CarInsuranceService) List(ctx context.Context, filter model.CarInsuranceFilter) ([]model.CarInsurance, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "List"), utils.MetadataFromCtx(ctx))

	insurances, err := s.presenter.List(ctx, filter)
	if err != nil {
		log.Warn("listing car insurances", pkglog.Err(err))

		return nil, err
	}

	return insurances, nil
}

func (s *CarInsuranceService) Update(ctx context.Context, id string, data model.CarInsuranceUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Update"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.Update(ctx, id, data); err != nil {
		log.Warn("updating car insurance", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarInsuranceService) Delete(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "Delete"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.Delete(ctx, id); err != nil {
		log.Warn("deleting car insurance", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *CarInsuranceService) GetImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.presenter.GetImageUploadData(ctx)
	if err != nil {
		log.Warn("getting car insurance image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}
