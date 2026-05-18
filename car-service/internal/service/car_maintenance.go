package service

import (
	"context"
	"log/slog"
	"time"

	"carsharing/car-service/internal/model"
	pkglog "carsharing/car-service/internal/pkg/log"
	"carsharing/car-service/internal/pkg/utils"
	"carsharing/car-service/internal/validation"
	"github.com/go-playground/validator/v10"
)

type CarMaintenanceService struct {
	templateRepo     CarMaintenanceTemplateRepository
	recordRepo       CarMaintenanceRecordRepository
	serviceStateRepo CarServiceStateRepository
	carRepo          CarRepository
	carService       *CarService
	objectStorage    ObjectStorage

	validate *validator.Validate
	log      *slog.Logger
}

func NewCarMaintenanceService(
	templateRepo CarMaintenanceTemplateRepository,
	recordRepo CarMaintenanceRecordRepository,
	serviceStateRepo CarServiceStateRepository,
	carRepo CarRepository,
	carService *CarService,
	objectStorage ObjectStorage,
	validate *validator.Validate,
	log *slog.Logger,
) *CarMaintenanceService {
	s := &CarMaintenanceService{
		templateRepo:     templateRepo,
		recordRepo:       recordRepo,
		serviceStateRepo: serviceStateRepo,
		carRepo:          carRepo,
		carService:       carService,
		objectStorage:    objectStorage,
		validate:         validate,
	}

	s.log = pkglog.WithComponent(log, "service.CarMaintenanceService")

	return s
}

func (s *CarMaintenanceService) CreateTemplate(ctx context.Context, createInput model.CarMaintenanceTemplateCreateInput) (string, error) {
	const method = "CreateTemplate"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, createInput); err != nil {
		return "", handleError(logger, err)
	}

	now := time.Now()
	template := model.CarMaintenanceTemplate{
		Name:        createInput.Name,
		KmInterval:  createInput.KmInterval,
		DayInterval: createInput.DayInterval,
		IsMandatory: createInput.IsMandatory,
		WarnPct:     createInput.WarnPct,
		PullPct:     createInput.PullPct,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	id, err := s.templateRepo.Insert(ctx, template)
	if err != nil {
		return "", handleError(logger, err)
	}

	return id, nil
}

func (s *CarMaintenanceService) GetTemplate(ctx context.Context, id string) (model.CarMaintenanceTemplate, error) {
	const method = "GetTemplate"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	template, err := s.templateRepo.FindByID(ctx, id)
	if err != nil {
		return model.CarMaintenanceTemplate{}, handleError(logger, err)
	}

	return template, nil
}

func (s *CarMaintenanceService) GetAllTemplates(ctx context.Context, filterInput model.CarMaintenanceTemplateFilterInput) ([]model.CarMaintenanceTemplate, error) {
	const method = "GetAllTemplates"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return nil, handleError(logger, err)
	}
	filter := maintenanceTemplateFilterFromInput(filterInput)

	templates, err := s.templateRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	return templates, nil
}

func (s *CarMaintenanceService) UpdateTemplate(ctx context.Context, id string, updateInput model.CarMaintenanceTemplateUpdateInput) error {
	const method = "UpdateTemplate"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, updateInput)
	if err != nil {
		return handleError(logger, err)
	}

	err = s.templateRepo.Update(ctx, id, model.CarMaintenanceTemplateUpdate{
		Name:        updateInput.Name,
		KmInterval:  updateInput.KmInterval,
		DayInterval: updateInput.DayInterval,
		IsMandatory: updateInput.IsMandatory,
		WarnPct:     updateInput.WarnPct,
		PullPct:     updateInput.PullPct,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarMaintenanceService) DeleteTemplate(ctx context.Context, id string) error {
	const method = "DeleteTemplate"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	err := s.templateRepo.Delete(ctx, id)
	if err != nil {
		return handleError(logger, err)
	}

	return nil
}

func (s *CarMaintenanceService) GetRecord(ctx context.Context, id string) (model.CarMaintenanceRecord, error) {
	const method = "GetRecord"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	record, err := s.recordRepo.FindByID(ctx, id)
	if err != nil {
		return model.CarMaintenanceRecord{}, handleError(logger, err)
	}

	for i := range record.ReceiptImages {
		url, err := s.objectStorage.GetPresignedURL(ctx, *record.ReceiptImages[i].Key)
		if err != nil {
			return model.CarMaintenanceRecord{}, handleError(logger, err)
		}
		record.ReceiptImages[i].URL = &url
	}

	return record, nil
}

func (s *CarMaintenanceService) GetRecords(ctx context.Context, filterInput model.CarMaintenanceRecordFilterInput) ([]model.CarMaintenanceRecord, error) {
	const method = "GetRecords"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	err := validation.ValidateInput(s.validate, filterInput)
	if err != nil {
		return nil, handleError(logger, err)
	}
	filter := maintenanceRecordFilterFromInput(filterInput)

	records, err := s.recordRepo.Find(ctx, filter)
	if err != nil {
		return nil, handleError(logger, err)
	}

	for i := range records {
		for j := range records[i].ReceiptImages {
			url, err := s.objectStorage.GetPresignedURL(ctx, *records[i].ReceiptImages[j].Key)
			if err != nil {
				return nil, handleError(logger, err)
			}
			records[i].ReceiptImages[j].URL = &url
		}
	}

	return records, nil
}

func (s *CarMaintenanceService) CompleteRecord(ctx context.Context, id string, completeInput model.CarMaintenanceRecordCompleteInput) error {
	const method = "CompleteRecord"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	if err := validation.ValidateInput(s.validate, completeInput); err != nil {
		return handleError(logger, err)
	}

	record, err := s.recordRepo.FindByID(ctx, id)
	if err != nil {
		return handleError(logger, err)
	}

	template, err := s.templateRepo.FindByID(ctx, record.TemplateID)
	if err != nil {
		return handleError(logger, err)
	}

	now := time.Now()
	completedStatus := model.MaintenanceRecordStatusCompleted
	if err = s.recordRepo.Update(ctx, id, model.CarMaintenanceRecordUpdate{
		Status:           &completedStatus,
		CompletedKM:      &completeInput.CompletedKM,
		CostTenge:        &completeInput.CostTenge,
		CompletedAt:      &now,
		Notes:            completeInput.Notes,
		ReceiptImageKeys: completeInput.ReceiptImageKeys,
		UpdatedAt:        now,
	}); err != nil {
		return handleError(logger, err)
	}

	// Reset the service-state clock for this car-template pair.
	state := model.CarServiceState{
		CarID:      record.CarID,
		TemplateID: record.TemplateID,
		LastKM:     completeInput.CompletedKM,
		LastDate:   &now,
	}
	if template.KmInterval != nil {
		nextKM := completeInput.CompletedKM + *template.KmInterval
		state.NextDueKM = &nextKM
	}
	if template.DayInterval != nil {
		nextDate := now.AddDate(0, 0, int(*template.DayInterval))
		state.NextDueDate = &nextDate
	}

	if err = s.serviceStateRepo.Upsert(ctx, state); err != nil {
		return handleError(logger, err)
	}

	if err = s.carService.UpdateCarStatus(
		ctx, record.CarID,
		model.CarStatusUpdateInput{Status: string(model.CarStatusAvailable)},
	); err != nil {
		return handleError(logger, err)
	}

	logger.Info("maintenance record completed",
		slog.String("recordID", record.ID),
		slog.String("carID", record.CarID),
		slog.Int("completedKM", int(completeInput.CompletedKM)),
	)

	return nil
}

func (s *CarMaintenanceService) GetReceiptImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	const method = "GetReceiptImageUploadData"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)

	data, err := s.objectStorage.GetMaintenanceReceiptImageUploadData(ctx)
	if err != nil {
		return model.ImageUploadData{}, handleError(logger, err)
	}

	return data, nil
}

func (s *CarMaintenanceService) EvaluateCarMaintenance(ctx context.Context, carID string) error {
	const method = "EvaluateCarMaintenance"
	logger := pkglog.WithMethod(s.log, method)

	md := utils.MetadataFromCtx(ctx)
	logger = pkglog.WithMetadata(logger, md)
	logger = logger.With(slog.String("carID", carID))

	car, err := s.carRepo.FindByID(ctx, carID)
	if err != nil {
		return handleError(logger, err)
	}

	states, err := s.serviceStateRepo.FindAll(ctx, model.CarServiceStateFilter{CarID: &carID})
	if err != nil {
		return handleError(logger, err)
	}

	for _, state := range states {
		template, err := s.templateRepo.FindByID(ctx, state.TemplateID)
		if err != nil {
			logger.Error("failed to load template",
				slog.String("templateID", state.TemplateID),
				slog.String("error", err.Error()),
			)
			continue
		}

		pct := maintenancePct(car.MileageKM, state, template)

		switch {
		case pct >= template.PullPct:
			if err = s.createWorkOrder(ctx, car, template, "urgent"); err != nil {
				logger.Error("failed to create urgent work order",
					slog.String("templateName", template.Name),
					slog.String("error", err.Error()),
				)
				continue
			}

			if err = s.carService.UpdateCarStatus(
				ctx, carID,
				model.CarStatusUpdateInput{Status: string(model.CarStatusMaintenance)},
			); err != nil {
				logger.Error("failed to transition car to maintenance",
					slog.String("templateName", template.Name),
					slog.String("error", err.Error()),
				)
			}

		case pct >= template.WarnPct:
			if err = s.createWorkOrder(ctx, car, template, "scheduled"); err != nil {
				logger.Error("failed to create scheduled work order",
					slog.String("templateName", template.Name),
					slog.String("error", err.Error()),
				)
			}
		}
	}

	return nil
}

func (s *CarMaintenanceService) createWorkOrder(ctx context.Context, car model.Car, template model.CarMaintenanceTemplate, priority string) error {
	var dueBy *time.Time
	if template.DayInterval != nil {
		t := time.Now().AddDate(0, 0, int(*template.DayInterval)/10)
		dueBy = &t
	}

	now := time.Now()
	_, err := s.recordRepo.Insert(ctx, model.CarMaintenanceRecord{
		CarID:      car.ID,
		TemplateID: template.ID,
		Status:     model.MaintenanceRecordStatusPending,
		OdometerAt: int32(car.MileageKM),
		DueBy:      dueBy,
		CreatedAt:  now,
		UpdatedAt:  now,
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
