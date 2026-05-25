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

type CarMaintenanceService struct {
	log      *slog.Logger
	validate *validator.Validate

	templateRepo     CarMaintenanceTemplateRepository
	recordRepo       CarMaintenanceRecordRepository
	serviceStateRepo CarServiceStateRepository
	carRepo          CarRepository
	carService       *CarService
	objectStorage    ObjectStorage
}

func NewCarMaintenanceService(
	log *slog.Logger,
	validate *validator.Validate,
	templateRepo CarMaintenanceTemplateRepository,
	recordRepo CarMaintenanceRecordRepository,
	serviceStateRepo CarServiceStateRepository,
	carRepo CarRepository,
	carService *CarService,
	objectStorage ObjectStorage,
) *CarMaintenanceService {
	return &CarMaintenanceService{
		log:              pkglog.WithComponent(log, "service.CarMaintenanceService"),
		validate:         validate,
		templateRepo:     templateRepo,
		recordRepo:       recordRepo,
		serviceStateRepo: serviceStateRepo,
		carRepo:          carRepo,
		carService:       carService,
		objectStorage:    objectStorage,
	}
}

func (s *CarMaintenanceService) CreateTemplate(ctx context.Context, data validation.CarMaintenanceTemplateCreate) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CreateTemplate"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return "", err
	}

	now := time.Now()
	id, err := s.templateRepo.Insert(ctx, model.CarMaintenanceTemplate{
		Name:        data.Name,
		KmInterval:  data.KmInterval,
		DayInterval: data.DayInterval,
		IsMandatory: data.IsMandatory,
		WarnPct:     data.WarnPct,
		PullPct:     data.PullPct,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		log.Error("repo: inserting maintenance template", pkglog.Err(err))
		return "", err
	}

	return id, nil
}

func (s *CarMaintenanceService) GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetTemplate"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.CarMaintenanceTemplate{}, err
	}

	template, err := s.templateRepo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrCarMaintenanceTemplateNotFound) {
			log.Error("repo: finding maintenance template by id", pkglog.Err(err))
		}
		return model.CarMaintenanceTemplate{}, err
	}

	return template, nil
}

func (s *CarMaintenanceService) ListTemplates(ctx context.Context, filter validation.CarMaintenanceTemplateFilter) ([]model.CarMaintenanceTemplate, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "ListTemplates"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	templates, err := s.templateRepo.Find(ctx, maintenanceTemplateFilter(filter))
	if err != nil {
		log.Error("repo: listing maintenance templates", pkglog.Err(err))
		return nil, err
	}

	return templates, nil
}

func (s *CarMaintenanceService) UpdateTemplate(ctx context.Context, id string, data validation.CarMaintenanceTemplateUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "UpdateTemplate"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	if err := s.templateRepo.Update(ctx, id, model.CarMaintenanceTemplateUpdate{
		Name:        data.Name,
		KmInterval:  data.KmInterval,
		DayInterval: data.DayInterval,
		IsMandatory: data.IsMandatory,
		WarnPct:     data.WarnPct,
		PullPct:     data.PullPct,
		UpdatedAt:   time.Now(),
	}); err != nil {
		if !errors.Is(err, model.ErrCarMaintenanceTemplateNotFound) {
			log.Error("repo: updating maintenance template", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarMaintenanceService) DeleteTemplate(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "DeleteTemplate"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := s.templateRepo.Delete(ctx, id); err != nil {
		if !errors.Is(err, model.ErrCarMaintenanceTemplateNotFound) {
			log.Error("repo: deleting maintenance template", pkglog.Err(err))
		}
		return err
	}

	return nil
}

func (s *CarMaintenanceService) GetRecord(ctx context.Context, id string) (model.CarMaintenanceRecord, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetRecord"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return model.CarMaintenanceRecord{}, err
	}

	record, err := s.recordRepo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrCarMaintenanceRecordNotFound) {
			log.Error("repo: finding maintenance record by id", pkglog.Err(err))
		}
		return model.CarMaintenanceRecord{}, err
	}

	for i := range record.ReceiptImages {
		url, err := s.objectStorage.GetPresignedURL(ctx, record.ReceiptImages[i].Key)
		if err != nil {
			log.Error("object storage: getting presigned url", pkglog.Err(err))
			return model.CarMaintenanceRecord{}, err
		}
		record.ReceiptImages[i].URL = url
	}

	return record, nil
}

func (s *CarMaintenanceService) ListRecords(ctx context.Context, filter validation.CarMaintenanceRecordFilter) ([]model.CarMaintenanceRecord, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "ListRecords"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateInput(s.validate, filter); err != nil {
		return nil, err
	}

	records, err := s.recordRepo.Find(ctx, maintenanceRecordFilter(filter))
	if err != nil {
		log.Error("repo: listing maintenance records", pkglog.Err(err))
		return nil, err
	}

	for i := range records {
		for j := range records[i].ReceiptImages {
			url, err := s.objectStorage.GetPresignedURL(ctx, records[i].ReceiptImages[j].Key)
			if err != nil {
				log.Warn("object storage: getting presigned url", pkglog.Err(err))
				continue
			}
			records[i].ReceiptImages[j].URL = url
		}
	}

	return records, nil
}

func (s *CarMaintenanceService) CompleteRecord(ctx context.Context, id string, data validation.CarMaintenanceRecordComplete) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CompleteRecord"), utils.MetadataFromCtx(ctx))

	if err := validation.ValidateID(s.validate, id); err != nil {
		return err
	}

	if err := validation.ValidateInput(s.validate, data); err != nil {
		return err
	}

	record, err := s.recordRepo.FindByID(ctx, id)
	if err != nil {
		if !errors.Is(err, model.ErrCarMaintenanceRecordNotFound) {
			log.Error("repo: finding maintenance record by id", pkglog.Err(err))
		}
		return err
	}

	template, err := s.templateRepo.FindByID(ctx, record.TemplateID)
	if err != nil {
		if !errors.Is(err, model.ErrCarMaintenanceTemplateNotFound) {
			log.Error("repo: finding maintenance template by id", pkglog.Err(err))
		}
		return err
	}

	now := time.Now()
	completedStatus := model.MaintenanceRecordStatusCompleted
	recordUpdate := model.CarMaintenanceRecordUpdate{
		Status:           &completedStatus,
		CompletedKM:      &data.CompletedKM,
		CostTenge:        &data.CostTenge,
		CompletedAt:      &now,
		Notes:            data.Notes,
		ReceiptImageKeys: data.ReceiptImageKeys,
		UpdatedAt:        now,
	}

	state := model.CarServiceState{
		CarID:      record.CarID,
		TemplateID: record.TemplateID,
		LastKM:     data.CompletedKM,
		LastDate:   &now,
	}
	if template.KmInterval != nil {
		nextKM := data.CompletedKM + *template.KmInterval
		state.NextDueKM = &nextKM
	}
	if template.DayInterval != nil {
		nextDate := now.AddDate(0, 0, int(*template.DayInterval))
		state.NextDueDate = &nextDate
	}

	if err = s.recordRepo.UpdateWithServiceState(ctx, id, recordUpdate, state); err != nil {
		if !errors.Is(err, model.ErrCarMaintenanceRecordNotFound) {
			log.Error("repo: completing maintenance record", pkglog.Err(err))
		}
		return err
	}

	if err = s.carService.UpdateCarStatus(
		ctx, record.CarID,
		validation.CarStatusUpdate{Status: string(model.CarStatusAvailable)},
	); err != nil {
		return err
	}

	log.Info("maintenance record completed",
		slog.String("recordID", record.ID),
		slog.String("carID", record.CarID),
		slog.Int("completedKM", int(data.CompletedKM)),
	)

	return nil
}

func (s *CarMaintenanceService) GetReceiptImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "GetReceiptImageUploadData"), utils.MetadataFromCtx(ctx))

	data, err := s.objectStorage.GetMaintenanceReceiptImageUploadData(ctx)
	if err != nil {
		log.Error("object storage: getting upload data", pkglog.Err(err))
		return sharedmodel.ImageUploadData{}, err
	}

	return data, nil
}

func (s *CarMaintenanceService) EvaluateCarMaintenance(ctx context.Context, carID string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "EvaluateCarMaintenance"), utils.MetadataFromCtx(ctx))
	log = log.With(slog.String("carID", carID))

	if err := validation.ValidateID(s.validate, carID); err != nil {
		return err
	}

	car, err := s.carRepo.FindByID(ctx, carID)
	if err != nil {
		if !errors.Is(err, model.ErrCarNotFound) {
			log.Error("repo: finding car by id", pkglog.Err(err))
		}
		return err
	}

	states, err := s.serviceStateRepo.FindAll(ctx, model.CarServiceStateFilter{CarID: &carID})
	if err != nil {
		log.Error("repo: finding car service states", pkglog.Err(err))
		return err
	}

	for _, state := range states {
		template, err := s.templateRepo.FindByID(ctx, state.TemplateID)
		if err != nil {
			log.Error("repo: finding maintenance template",
				slog.String("templateID", state.TemplateID),
				pkglog.Err(err),
			)
			continue
		}

		pct := maintenancePct(car.MileageKM, state, template)

		switch {
		case pct >= template.PullPct:
			if err = s.createWorkOrder(ctx, car, template, "urgent"); err != nil {
				log.Error("failed to create urgent work order",
					slog.String("templateName", template.Name),
					pkglog.Err(err),
				)
				continue
			}

			if err = s.carService.UpdateCarStatus(
				ctx, carID,
				validation.CarStatusUpdate{Status: string(model.CarStatusMaintenance)},
			); err != nil {
				log.Error("failed to transition car to maintenance",
					slog.String("templateName", template.Name),
					pkglog.Err(err),
				)
			}

		case pct >= template.WarnPct:
			if err = s.createWorkOrder(ctx, car, template, "scheduled"); err != nil {
				log.Error("failed to create scheduled work order",
					slog.String("templateName", template.Name),
					pkglog.Err(err),
				)
			}
		}
	}

	return nil
}

func (s *CarMaintenanceService) createWorkOrder(ctx context.Context, car model.Car, template model.CarMaintenanceTemplate, priority string) error {
	var dueBy *time.Time
	if template.DayInterval != nil {
		t := time.Now().AddDate(0, 0, int(*template.DayInterval))
		dueBy = &t
	}

	now := time.Now()
	_, err := s.recordRepo.Insert(ctx, model.CarMaintenanceRecord{
		CarID:              car.ID,
		TemplateID:         template.ID,
		Status:             model.MaintenanceRecordStatusPending,
		MileageAtWarningKM: int32(car.MileageKM),
		DueBy:              dueBy,
		CreatedAt:          now,
		UpdatedAt:          now,
	})
	if err != nil {
		return err
	}

	s.log.Info("work order created",
		slog.String("carID", car.ID),
		slog.String("template", template.Name),
		slog.String("priority", priority),
	)

	return nil
}

func maintenanceTemplateFilter(filter validation.CarMaintenanceTemplateFilter) model.CarMaintenanceTemplateFilter {
	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	return model.CarMaintenanceTemplateFilter{
		IsMandatory: filter.IsMandatory,
		Pagination:  &sharedmodel.Pagination{Limit: filter.Pagination.Limit, Offset: filter.Pagination.Offset},
	}
}

func maintenanceRecordFilter(filter validation.CarMaintenanceRecordFilter) model.CarMaintenanceRecordFilter {
	repoFilter := model.CarMaintenanceRecordFilter{
		CarID:      filter.CarID,
		TemplateID: filter.TemplateID,
	}
	if filter.Status != nil {
		st, _ := model.MaintenanceRecordStatusFromString(*filter.Status)
		repoFilter.Status = &st
	}
	if filter.Pagination == nil {
		filter.Pagination = sharedvalidation.DefaultPagination()
	}
	repoFilter.Pagination = &sharedmodel.Pagination{Limit: filter.Pagination.Limit, Offset: filter.Pagination.Offset}
	return repoFilter
}

func maintenancePct(currentMileageKM int64, state model.CarServiceState, template model.CarMaintenanceTemplate) float64 {
	var kmPct, dayPct float64

	if template.KmInterval != nil && *template.KmInterval > 0 {
		kmPct = float64(int32(currentMileageKM)-state.LastKM) / float64(*template.KmInterval)
	}

	if template.DayInterval != nil && *template.DayInterval > 0 && state.LastDate != nil {
		dayPct = time.Since(*state.LastDate).Hours() / 24 / float64(*template.DayInterval)
	}

	if kmPct > dayPct {
		return kmPct
	}

	return dayPct
}
