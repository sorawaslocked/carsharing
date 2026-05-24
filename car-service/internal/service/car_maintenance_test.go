package service

import (
	"context"
	"testing"
	"time"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/service/mocks"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestCarMaintenanceService(
	t *testing.T,
	templateRepo CarMaintenanceTemplateRepository,
	recordRepo CarMaintenanceRecordRepository,
	serviceStateRepo CarServiceStateRepository,
	carRepo CarRepository,
	carService *CarService,
	objectStorage ObjectStorage,
) *CarMaintenanceService {
	t.Helper()
	return NewCarMaintenanceService(
		discardLogger(), newTestValidator(t),
		templateRepo, recordRepo, serviceStateRepo, carRepo, carService, objectStorage,
	)
}

// newTestCarServiceForMaintenance builds a CarService backed by the given carRepo
// with no status-log repo and no event publisher (both are nil-safe in UpdateCarStatus).
func newTestCarServiceForMaintenance(t *testing.T, carRepo CarRepository) *CarService {
	t.Helper()
	return NewCarService(discardLogger(), newTestValidator(t), nil, carRepo, nil, nil, nil, nil, nil)
}

// ── maintenancePct ────────────────────────────────────────────────────────────

func TestMaintenancePct(t *testing.T) {
	kmInterval := int32(10_000)
	dayInterval := int32(365)

	tests := []struct {
		name        string
		mileageKM   int64
		state       model.CarServiceState
		template    model.CarMaintenanceTemplate
		expectAbove float64
		expectBelow float64
	}{
		{
			name:        "no intervals → always 0",
			mileageKM:   99_999,
			state:       model.CarServiceState{},
			template:    model.CarMaintenanceTemplate{},
			expectAbove: 0,
			expectBelow: 0.001,
		},
		{
			name:        "km 50% done",
			mileageKM:   5_000,
			state:       model.CarServiceState{LastKM: 0},
			template:    model.CarMaintenanceTemplate{KmInterval: &kmInterval},
			expectAbove: 0.499,
			expectBelow: 0.501,
		},
		{
			name:        "km 100% done",
			mileageKM:   10_000,
			state:       model.CarServiceState{LastKM: 0},
			template:    model.CarMaintenanceTemplate{KmInterval: &kmInterval},
			expectAbove: 0.999,
			expectBelow: 1.001,
		},
		{
			name:        "km overdue (>100%)",
			mileageKM:   15_000,
			state:       model.CarServiceState{LastKM: 0},
			template:    model.CarMaintenanceTemplate{KmInterval: &kmInterval},
			expectAbove: 1.4,
			expectBelow: 1.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pct := maintenancePct(tt.mileageKM, tt.state, tt.template)
			assert.Greater(t, pct, tt.expectAbove-0.001, "pct should be > %f", tt.expectAbove)
			assert.Less(t, pct, tt.expectBelow+0.001, "pct should be < %f", tt.expectBelow)
		})
	}

	t.Run("day pct wins when higher than km pct", func(t *testing.T) {
		last := time.Now().Add(-400 * 24 * time.Hour) // ~400 days ago
		state := model.CarServiceState{LastKM: 0, LastDate: &last}
		tmpl := model.CarMaintenanceTemplate{
			KmInterval:  &kmInterval,  // 0% km progress
			DayInterval: &dayInterval, // > 100% day progress
		}
		pct := maintenancePct(0, state, tmpl)
		assert.Greater(t, pct, 1.0)
	})
}

// ── CreateTemplate ────────────────────────────────────────────────────────────

func TestCarMaintenanceServiceCreateTemplate(t *testing.T) {
	ctx := context.Background()

	validInput := validation.CarMaintenanceTemplateCreate{
		Name:        "Oil Change",
		IsMandatory: true,
		WarnPct:     0.8,
		PullPct:     1.0,
	}

	t.Run("happy path returns inserted id", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(tmpl model.CarMaintenanceTemplate) bool {
			return tmpl.Name == "Oil Change" && tmpl.IsMandatory
		})).Return("tmpl-123", nil)

		id, err := svc.CreateTemplate(ctx, validInput)
		assert.NoError(t, err)
		assert.Equal(t, "tmpl-123", id)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrSql)

		_, err := svc.CreateTemplate(ctx, validInput)
		assert.Error(t, err)
	})

	t.Run("validation rejects missing name", func(t *testing.T) {
		svc := newTestCarMaintenanceService(t, nil, nil, nil, nil, nil, nil)

		_, err := svc.CreateTemplate(ctx, validation.CarMaintenanceTemplateCreate{PullPct: 1.0})
		assert.Error(t, err)
	})
}

// ── GetTemplate ───────────────────────────────────────────────────────────────

func TestCarMaintenanceServiceGetTemplate(t *testing.T) {
	ctx := context.Background()
	tmplID := "e0000000-0000-4000-8000-000000000001"

	t.Run("returns template", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().FindByID(ctx, tmplID).Return(
			model.CarMaintenanceTemplate{ID: tmplID, Name: "Oil Change"}, nil,
		)

		got, err := svc.GetTemplate(ctx, tmplID)
		assert.NoError(t, err)
		assert.Equal(t, tmplID, got.ID)
	})

	t.Run("not found returns ErrCarMaintenanceTemplateNotFound", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().FindByID(ctx, tmplID).Return(model.CarMaintenanceTemplate{}, model.ErrCarMaintenanceTemplateNotFound)

		_, err := svc.GetTemplate(ctx, tmplID)
		assert.ErrorIs(t, err, model.ErrCarMaintenanceTemplateNotFound)
	})
}

// ── GetAllTemplates ───────────────────────────────────────────────────────────

func TestCarMaintenanceServiceGetAllTemplates(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty list", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().Find(ctx, mock.Anything).Return(nil, nil)

		got, err := svc.ListTemplates(ctx, validation.CarMaintenanceTemplateFilter{})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("mandatory filter is forwarded to repo", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		mandatory := true
		templateRepo.EXPECT().Find(ctx, mock.MatchedBy(func(f model.CarMaintenanceTemplateFilter) bool {
			return f.IsMandatory != nil && *f.IsMandatory
		})).Return(nil, nil)

		_, err := svc.ListTemplates(ctx, validation.CarMaintenanceTemplateFilter{IsMandatory: &mandatory})
		assert.NoError(t, err)
	})
}

// ── UpdateTemplate ────────────────────────────────────────────────────────────

func TestCarMaintenanceServiceUpdateTemplate(t *testing.T) {
	ctx := context.Background()
	tmplID := "e0000000-0000-4000-8000-000000000001"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().Update(ctx, tmplID, mock.Anything).Return(nil)

		assert.NoError(t, svc.UpdateTemplate(ctx, tmplID, validation.CarMaintenanceTemplateUpdate{}))
	})

	t.Run("not found returns ErrCarMaintenanceTemplateNotFound", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().Update(ctx, tmplID, mock.Anything).Return(model.ErrCarMaintenanceTemplateNotFound)

		err := svc.UpdateTemplate(ctx, tmplID, validation.CarMaintenanceTemplateUpdate{})
		assert.ErrorIs(t, err, model.ErrCarMaintenanceTemplateNotFound)
	})
}

// ── DeleteTemplate ────────────────────────────────────────────────────────────

func TestCarMaintenanceServiceDeleteTemplate(t *testing.T) {
	ctx := context.Background()
	tmplID := "e0000000-0000-4000-8000-000000000001"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().Delete(ctx, tmplID).Return(nil)

		assert.NoError(t, svc.DeleteTemplate(ctx, tmplID))
	})

	t.Run("not found returns ErrCarMaintenanceTemplateNotFound", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, nil, nil, nil, nil)

		templateRepo.EXPECT().Delete(ctx, tmplID).Return(model.ErrCarMaintenanceTemplateNotFound)

		err := svc.DeleteTemplate(ctx, tmplID)
		assert.ErrorIs(t, err, model.ErrCarMaintenanceTemplateNotFound)
	})
}

// ── GetRecord ─────────────────────────────────────────────────────────────────

func TestCarMaintenanceServiceGetRecord(t *testing.T) {
	ctx := context.Background()
	recordID := "f0000000-0000-4000-8000-000000000001"

	t.Run("returns record with no receipt images", func(t *testing.T) {
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarMaintenanceService(t, nil, recordRepo, nil, nil, nil, storage)

		recordRepo.EXPECT().FindByID(ctx, recordID).Return(
			model.CarMaintenanceRecord{ID: recordID}, nil,
		)

		got, err := svc.GetRecord(ctx, recordID)
		assert.NoError(t, err)
		assert.Equal(t, recordID, got.ID)
	})

	t.Run("populates presigned URL for each receipt image", func(t *testing.T) {
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarMaintenanceService(t, nil, recordRepo, nil, nil, nil, storage)

		key := "receipts/invoice.pdf"
		presigned := "https://cdn/invoice"
		recordRepo.EXPECT().FindByID(ctx, recordID).Return(model.CarMaintenanceRecord{
			ID:            recordID,
			ReceiptImages: []sharedmodel.Image{{Key: key}},
		}, nil)
		storage.EXPECT().GetPresignedURL(ctx, key).Return(presigned, nil)

		got, err := svc.GetRecord(ctx, recordID)
		assert.NoError(t, err)
		assert.Equal(t, presigned, got.ReceiptImages[0].URL)
	})

	t.Run("not found returns ErrCarMaintenanceRecordNotFound", func(t *testing.T) {
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarMaintenanceService(t, nil, recordRepo, nil, nil, nil, storage)

		recordRepo.EXPECT().FindByID(ctx, recordID).Return(model.CarMaintenanceRecord{}, model.ErrCarMaintenanceRecordNotFound)

		_, err := svc.GetRecord(ctx, recordID)
		assert.ErrorIs(t, err, model.ErrCarMaintenanceRecordNotFound)
	})
}

// ── GetRecords ────────────────────────────────────────────────────────────────

func TestCarMaintenanceServiceGetRecords(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty list", func(t *testing.T) {
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarMaintenanceService(t, nil, recordRepo, nil, nil, nil, storage)

		recordRepo.EXPECT().Find(ctx, mock.Anything).Return(nil, nil)

		got, err := svc.ListRecords(ctx, validation.CarMaintenanceRecordFilter{})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarMaintenanceService(t, nil, recordRepo, nil, nil, nil, storage)

		recordRepo.EXPECT().Find(ctx, mock.Anything).Return(nil, model.ErrSql)

		_, err := svc.ListRecords(ctx, validation.CarMaintenanceRecordFilter{})
		assert.Error(t, err)
	})
}

// ── CompleteRecord ────────────────────────────────────────────────────────────

func TestCarMaintenanceServiceCompleteRecord(t *testing.T) {
	ctx := context.Background()

	const (
		recordID   = "f0000000-0000-4000-8000-000000000001"
		carID      = "c0000000-0000-4000-8000-000000000002"
		templateID = "tmpl-789"
	)

	t.Run("happy path updates record, resets service state, releases car", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		stateRepo := mocks.NewMockCarServiceStateRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		carService := newTestCarServiceForMaintenance(t, carRepo)
		svc := newTestCarMaintenanceService(t, templateRepo, recordRepo, stateRepo, nil, carService, nil)

		kmInterval := int32(5_000)
		dayInterval := int32(180)

		recordRepo.EXPECT().FindByID(ctx, recordID).Return(model.CarMaintenanceRecord{
			ID:         recordID,
			CarID:      carID,
			TemplateID: templateID,
		}, nil)
		templateRepo.EXPECT().FindByID(ctx, templateID).Return(model.CarMaintenanceTemplate{
			ID:          templateID,
			KmInterval:  &kmInterval,
			DayInterval: &dayInterval,
		}, nil)
		recordRepo.EXPECT().Update(ctx, recordID, mock.MatchedBy(func(u model.CarMaintenanceRecordUpdate) bool {
			return u.Status != nil && *u.Status == model.MaintenanceRecordStatusCompleted &&
				u.CompletedKM != nil && *u.CompletedKM == 50_000
		})).Return(nil)
		stateRepo.EXPECT().Upsert(ctx, mock.MatchedBy(func(s model.CarServiceState) bool {
			return s.CarID == carID && s.TemplateID == templateID &&
				s.LastKM == 50_000 && s.NextDueKM != nil && s.NextDueDate != nil
		})).Return(nil)
		// carService.UpdateCarStatus → FindByID + Update
		carRepo.EXPECT().FindByID(ctx, carID).Return(
			model.Car{ID: carID, Status: model.CarStatusMaintenance}, nil,
		)
		carRepo.EXPECT().Update(ctx, carID, mock.MatchedBy(func(u model.CarUpdate) bool {
			return u.Status != nil && *u.Status == model.CarStatusAvailable
		})).Return(nil)

		err := svc.CompleteRecord(ctx, recordID, validation.CarMaintenanceRecordComplete{
			CompletedKM: 50_000,
			CostTenge:   15_000,
		})
		assert.NoError(t, err)
	})

	t.Run("km-only template sets NextDueKM but not NextDueDate", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		stateRepo := mocks.NewMockCarServiceStateRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		carService := newTestCarServiceForMaintenance(t, carRepo)
		svc := newTestCarMaintenanceService(t, templateRepo, recordRepo, stateRepo, nil, carService, nil)

		kmInterval := int32(5_000)

		recordRepo.EXPECT().FindByID(ctx, recordID).Return(model.CarMaintenanceRecord{
			ID: recordID, CarID: carID, TemplateID: templateID,
		}, nil)
		templateRepo.EXPECT().FindByID(ctx, templateID).Return(model.CarMaintenanceTemplate{
			ID:         templateID,
			KmInterval: &kmInterval,
		}, nil)
		recordRepo.EXPECT().Update(ctx, recordID, mock.Anything).Return(nil)
		stateRepo.EXPECT().Upsert(ctx, mock.MatchedBy(func(s model.CarServiceState) bool {
			return s.NextDueKM != nil && s.NextDueDate == nil
		})).Return(nil)
		carRepo.EXPECT().FindByID(ctx, carID).Return(
			model.Car{ID: carID, Status: model.CarStatusMaintenance}, nil,
		)
		carRepo.EXPECT().Update(ctx, carID, mock.Anything).Return(nil)

		err := svc.CompleteRecord(ctx, recordID, validation.CarMaintenanceRecordComplete{CompletedKM: 10_000})
		assert.NoError(t, err)
	})

	t.Run("record not found returns ErrCarMaintenanceRecordNotFound", func(t *testing.T) {
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		svc := newTestCarMaintenanceService(t, nil, recordRepo, nil, nil, nil, nil)

		recordRepo.EXPECT().FindByID(ctx, recordID).Return(model.CarMaintenanceRecord{}, model.ErrCarMaintenanceRecordNotFound)

		err := svc.CompleteRecord(ctx, recordID, validation.CarMaintenanceRecordComplete{CompletedKM: 1000})
		assert.ErrorIs(t, err, model.ErrCarMaintenanceRecordNotFound)
	})

	t.Run("template not found returns ErrCarMaintenanceTemplateNotFound", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		svc := newTestCarMaintenanceService(t, templateRepo, recordRepo, nil, nil, nil, nil)

		recordRepo.EXPECT().FindByID(ctx, recordID).Return(model.CarMaintenanceRecord{
			ID: recordID, CarID: carID, TemplateID: templateID,
		}, nil)
		templateRepo.EXPECT().FindByID(ctx, templateID).Return(model.CarMaintenanceTemplate{}, model.ErrCarMaintenanceTemplateNotFound)

		err := svc.CompleteRecord(ctx, recordID, validation.CarMaintenanceRecordComplete{CompletedKM: 1000})
		assert.ErrorIs(t, err, model.ErrCarMaintenanceTemplateNotFound)
	})
}

// ── GetReceiptImageUploadData ─────────────────────────────────────────────────

func TestCarMaintenanceServiceGetReceiptImageUploadData(t *testing.T) {
	ctx := context.Background()

	t.Run("returns upload data from object storage", func(t *testing.T) {
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarMaintenanceService(t, nil, nil, nil, nil, nil, storage)

		want := sharedmodel.ImageUploadData{ObjectKey: "receipts/abc.pdf", PresignedPutURL: "https://upload.example.com"}
		storage.EXPECT().GetMaintenanceReceiptImageUploadData(ctx).Return(want, nil)

		got, err := svc.GetReceiptImageUploadData(ctx)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("storage error is propagated", func(t *testing.T) {
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarMaintenanceService(t, nil, nil, nil, nil, nil, storage)

		storage.EXPECT().GetMaintenanceReceiptImageUploadData(ctx).Return(sharedmodel.ImageUploadData{}, model.ErrSql)

		_, err := svc.GetReceiptImageUploadData(ctx)
		assert.Error(t, err)
	})
}

// ── EvaluateCarMaintenance ────────────────────────────────────────────────────

func TestEvaluateCarMaintenance(t *testing.T) {
	ctx := context.Background()
	carID := "c0000000-0000-4000-8000-000000000001"
	templateID := "tmpl-abc"

	kmInterval := int32(10_000)
	warnPct := 0.8
	pullPct := 1.0

	template := model.CarMaintenanceTemplate{
		ID:         templateID,
		KmInterval: &kmInterval,
		WarnPct:    warnPct,
		PullPct:    pullPct,
	}
	state := model.CarServiceState{
		CarID:      carID,
		TemplateID: templateID,
		LastKM:     0,
	}

	t.Run("urgent path creates work order and transitions car to maintenance", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		stateRepo := mocks.NewMockCarServiceStateRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		carService := newTestCarServiceForMaintenance(t, carRepo)
		svc := newTestCarMaintenanceService(t, templateRepo, recordRepo, stateRepo, carRepo, carService, nil)

		// mileage 10 000 km since lastKM=0, interval=10 000 → pct=1.0 >= pullPct=1.0
		carRepo.EXPECT().FindByID(ctx, carID).
			Return(model.Car{ID: carID, MileageKM: 10_000, Status: model.CarStatusAvailable}, nil)
		stateRepo.EXPECT().FindAll(ctx, mock.Anything).Return([]model.CarServiceState{state}, nil)
		templateRepo.EXPECT().FindByID(ctx, templateID).Return(template, nil)
		recordRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(r model.CarMaintenanceRecord) bool {
			return r.CarID == carID && r.TemplateID == templateID &&
				r.Status == model.MaintenanceRecordStatusPending
		})).Return("rec-new", nil)
		// carService.UpdateCarStatus → FindByID again + Update
		carRepo.EXPECT().FindByID(ctx, carID).
			Return(model.Car{ID: carID, Status: model.CarStatusAvailable}, nil)
		carRepo.EXPECT().Update(ctx, carID, mock.MatchedBy(func(u model.CarUpdate) bool {
			return u.Status != nil && *u.Status == model.CarStatusMaintenance
		})).Return(nil)

		err := svc.EvaluateCarMaintenance(ctx, carID)
		assert.NoError(t, err)
	})

	t.Run("warn path creates work order but does not change car status", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		recordRepo := mocks.NewMockCarMaintenanceRecordRepository(t)
		stateRepo := mocks.NewMockCarServiceStateRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		carService := newTestCarServiceForMaintenance(t, carRepo)
		svc := newTestCarMaintenanceService(t, templateRepo, recordRepo, stateRepo, carRepo, carService, nil)

		// mileage 8 500 → pct=0.85 >= warnPct=0.8 but < pullPct=1.0
		carRepo.EXPECT().FindByID(ctx, carID).
			Return(model.Car{ID: carID, MileageKM: 8_500, Status: model.CarStatusAvailable}, nil)
		stateRepo.EXPECT().FindAll(ctx, mock.Anything).Return([]model.CarServiceState{state}, nil)
		templateRepo.EXPECT().FindByID(ctx, templateID).Return(template, nil)
		recordRepo.EXPECT().Insert(ctx, mock.Anything).Return("rec-warn", nil)
		// no carRepo.Update expected: status must NOT change

		err := svc.EvaluateCarMaintenance(ctx, carID)
		assert.NoError(t, err)
	})

	t.Run("below warn threshold creates no work order", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		stateRepo := mocks.NewMockCarServiceStateRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		carService := newTestCarServiceForMaintenance(t, carRepo)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, stateRepo, carRepo, carService, nil)

		// mileage 5 000 → pct=0.5 < warnPct=0.8
		carRepo.EXPECT().FindByID(ctx, carID).
			Return(model.Car{ID: carID, MileageKM: 5_000, Status: model.CarStatusAvailable}, nil)
		stateRepo.EXPECT().FindAll(ctx, mock.Anything).Return([]model.CarServiceState{state}, nil)
		templateRepo.EXPECT().FindByID(ctx, templateID).Return(template, nil)
		// no recordRepo.Insert and no carRepo.Update expected

		err := svc.EvaluateCarMaintenance(ctx, carID)
		assert.NoError(t, err)
	})

	t.Run("no service states returns no error and does nothing", func(t *testing.T) {
		stateRepo := mocks.NewMockCarServiceStateRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		carService := newTestCarServiceForMaintenance(t, carRepo)
		svc := newTestCarMaintenanceService(t, nil, nil, stateRepo, carRepo, carService, nil)

		carRepo.EXPECT().FindByID(ctx, carID).
			Return(model.Car{ID: carID, MileageKM: 50_000}, nil)
		stateRepo.EXPECT().FindAll(ctx, mock.Anything).Return(nil, nil)

		err := svc.EvaluateCarMaintenance(ctx, carID)
		assert.NoError(t, err)
	})

	t.Run("car not found returns ErrCarNotFound", func(t *testing.T) {
		carRepo := mocks.NewMockCarRepository(t)
		carService := newTestCarServiceForMaintenance(t, carRepo)
		svc := newTestCarMaintenanceService(t, nil, nil, nil, carRepo, carService, nil)

		carRepo.EXPECT().FindByID(ctx, carID).Return(model.Car{}, model.ErrCarNotFound)

		err := svc.EvaluateCarMaintenance(ctx, carID)
		assert.ErrorIs(t, err, model.ErrCarNotFound)
	})

	t.Run("template load failure is logged and skipped, not fatal", func(t *testing.T) {
		templateRepo := mocks.NewMockCarMaintenanceTemplateRepository(t)
		stateRepo := mocks.NewMockCarServiceStateRepository(t)
		carRepo := mocks.NewMockCarRepository(t)
		carService := newTestCarServiceForMaintenance(t, carRepo)
		svc := newTestCarMaintenanceService(t, templateRepo, nil, stateRepo, carRepo, carService, nil)

		carRepo.EXPECT().FindByID(ctx, carID).
			Return(model.Car{ID: carID, MileageKM: 50_000}, nil)
		stateRepo.EXPECT().FindAll(ctx, mock.Anything).Return([]model.CarServiceState{state}, nil)
		templateRepo.EXPECT().FindByID(ctx, templateID).Return(model.CarMaintenanceTemplate{}, model.ErrCarMaintenanceTemplateNotFound)

		// Must not return an error: template failures are non-fatal per implementation.
		err := svc.EvaluateCarMaintenance(ctx, carID)
		assert.NoError(t, err)
	})
}
