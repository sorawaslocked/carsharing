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

func newTestCarInsuranceService(t *testing.T, insuranceRepo CarInsuranceRepository, objectStorage ObjectStorage) *CarInsuranceService {
	t.Helper()
	return NewCarInsuranceService(discardLogger(), newTestValidator(t), insuranceRepo, objectStorage)
}

var (
	testCarID    = "00000000-0000-0000-0000-000000000001"
	testStartsAt = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	testEndsAt   = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
)

func validInsuranceCreateInput() validation.CarInsuranceCreate {
	return validation.CarInsuranceCreate{
		CarID:     testCarID,
		Type:      string(model.InsuranceTypeOSAGO),
		Provider:  "InsureCo",
		PolicyNum: "POL-001",
		StartsAt:  testStartsAt,
		ExpiresAt: testEndsAt,
	}
}

func TestCarInsuranceServiceCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("happy path returns inserted id", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		svc := newTestCarInsuranceService(t, repo, nil)

		repo.EXPECT().Insert(ctx, mock.MatchedBy(func(ins model.CarInsurance) bool {
			return ins.CarID == testCarID &&
				ins.Type == model.InsuranceTypeOSAGO &&
				ins.Status == model.InsuranceStatusActive
		})).Return("ins-123", nil)

		id, err := svc.Create(ctx, validInsuranceCreateInput())
		assert.NoError(t, err)
		assert.Equal(t, "ins-123", id)
	})

	t.Run("new insurance defaults to active status", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		svc := newTestCarInsuranceService(t, repo, nil)

		repo.EXPECT().Insert(ctx, mock.MatchedBy(func(ins model.CarInsurance) bool {
			return ins.Status == model.InsuranceStatusActive
		})).Return("ins-456", nil)

		_, err := svc.Create(ctx, validInsuranceCreateInput())
		assert.NoError(t, err)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		svc := newTestCarInsuranceService(t, repo, nil)

		repo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrSql)

		_, err := svc.Create(ctx, validInsuranceCreateInput())
		assert.Error(t, err)
	})

	t.Run("validation rejects invalid car id", func(t *testing.T) {
		svc := newTestCarInsuranceService(t, nil, nil)

		input := validInsuranceCreateInput()
		input.CarID = "not-a-uuid"

		_, err := svc.Create(ctx, input)
		assert.Error(t, err)
	})
}

func TestCarInsuranceServiceGet(t *testing.T) {
	ctx := context.Background()
	insID := "d0000000-0000-4000-8000-000000000001"

	t.Run("returns insurance with no images", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarInsuranceService(t, repo, storage)

		repo.EXPECT().FindByID(ctx, insID).Return(model.CarInsurance{ID: insID}, nil)

		got, err := svc.Get(ctx, insID)
		assert.NoError(t, err)
		assert.Equal(t, insID, got.ID)
	})

	t.Run("populates presigned URL for each image", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarInsuranceService(t, repo, storage)

		key := "insurance/doc.pdf"
		presigned := "https://cdn/doc"
		repo.EXPECT().FindByID(ctx, insID).Return(model.CarInsurance{
			ID:     insID,
			Images: []sharedmodel.Image{{Key: key}},
		}, nil)
		storage.EXPECT().GetPresignedURL(ctx, key).Return(presigned, nil)

		got, err := svc.Get(ctx, insID)
		assert.NoError(t, err)
		assert.Equal(t, presigned, got.Images[0].URL)
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarInsuranceService(t, repo, storage)

		repo.EXPECT().FindByID(ctx, insID).Return(model.CarInsurance{}, model.ErrNotFound)

		_, err := svc.Get(ctx, insID)
		assert.ErrorIs(t, err, model.ErrNotFound)
	})
}

func TestCarInsuranceServiceGetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty list", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarInsuranceService(t, repo, storage)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, nil)

		got, err := svc.List(ctx, validation.CarInsuranceFilter{})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("forwards car id filter to repo", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarInsuranceService(t, repo, storage)

		repo.EXPECT().Find(ctx, mock.MatchedBy(func(f model.CarInsuranceFilter) bool {
			return f.CarID != nil && *f.CarID == testCarID
		})).Return(nil, nil)

		got, err := svc.List(ctx, validation.CarInsuranceFilter{CarID: &testCarID})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarInsuranceService(t, repo, storage)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, model.ErrSql)

		_, err := svc.List(ctx, validation.CarInsuranceFilter{})
		assert.Error(t, err)
	})
}

func TestCarInsuranceServiceUpdate(t *testing.T) {
	ctx := context.Background()
	insID := "d0000000-0000-4000-8000-000000000001"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		svc := newTestCarInsuranceService(t, repo, nil)

		repo.EXPECT().Update(ctx, insID, mock.Anything).Return(nil)

		assert.NoError(t, svc.Update(ctx, insID, validation.CarInsuranceUpdate{}))
	})

	t.Run("status string is parsed to enum", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		svc := newTestCarInsuranceService(t, repo, nil)

		statusStr := string(model.InsuranceStatusExpired)
		repo.EXPECT().Update(ctx, insID, mock.MatchedBy(func(u model.CarInsuranceUpdate) bool {
			return u.Status != nil && *u.Status == model.InsuranceStatusExpired
		})).Return(nil)

		assert.NoError(t, svc.Update(ctx, insID, validation.CarInsuranceUpdate{Status: &statusStr}))
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		svc := newTestCarInsuranceService(t, repo, nil)

		repo.EXPECT().Update(ctx, insID, mock.Anything).Return(model.ErrNotFound)

		assert.ErrorIs(t, svc.Update(ctx, insID, validation.CarInsuranceUpdate{}), model.ErrNotFound)
	})
}

func TestCarInsuranceServiceDelete(t *testing.T) {
	ctx := context.Background()
	insID := "d0000000-0000-4000-8000-000000000001"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		svc := newTestCarInsuranceService(t, repo, nil)

		repo.EXPECT().Delete(ctx, insID).Return(nil)

		assert.NoError(t, svc.Delete(ctx, insID))
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		svc := newTestCarInsuranceService(t, repo, nil)

		repo.EXPECT().Delete(ctx, insID).Return(model.ErrNotFound)

		assert.ErrorIs(t, svc.Delete(ctx, insID), model.ErrNotFound)
	})
}

func TestCarInsuranceServiceGetImageUploadData(t *testing.T) {
	ctx := context.Background()

	t.Run("returns upload data from object storage", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarInsuranceService(t, repo, storage)

		want := sharedmodel.ImageUploadData{ObjectKey: "insurance/doc.pdf", PresignedPutURL: "https://upload.example.com"}
		storage.EXPECT().GetInsuranceImageUploadData(ctx).Return(want, nil)

		got, err := svc.GetImageUploadData(ctx)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("storage error is propagated", func(t *testing.T) {
		repo := mocks.NewMockCarInsuranceRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarInsuranceService(t, repo, storage)

		storage.EXPECT().GetInsuranceImageUploadData(ctx).Return(sharedmodel.ImageUploadData{}, model.ErrSql)

		_, err := svc.GetImageUploadData(ctx)
		assert.Error(t, err)
	})
}
