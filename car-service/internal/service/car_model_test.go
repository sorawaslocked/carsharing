package service

import (
	"context"
	"testing"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/service/mocks"
	"carsharing/car-service/internal/validation"
	sharedmodel "carsharing/shared/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestCarModelService(t *testing.T, carModelRepo CarModelRepository, objectStorage ObjectStorage) *CarModelService {
	t.Helper()
	return NewCarModelService(discardLogger(), newTestValidator(t), carModelRepo, objectStorage)
}

func TestCarModelServiceCreate(t *testing.T) {
	ctx := context.Background()

	validInput := validation.CarModelCreate{
		Brand:        "Toyota",
		Model:        "Camry",
		Year:         2020,
		FuelType:     string(model.CarFuelTypePetrol),
		Transmission: string(model.CarTransmissionAuto),
		BodyType:     string(model.CarBodyTypeSedan),
		Class:        string(model.CarClassComfort),
		Seats:        5,
	}

	t.Run("happy path returns inserted id", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		svc := newTestCarModelService(t, repo, nil)

		repo.EXPECT().Insert(ctx, mock.Anything).Return("model-123", nil)

		id, err := svc.Create(ctx, validInput)
		assert.NoError(t, err)
		assert.Equal(t, "model-123", id)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		svc := newTestCarModelService(t, repo, nil)

		repo.EXPECT().Insert(ctx, mock.Anything).Return("", model.ErrSql)

		_, err := svc.Create(ctx, validInput)
		assert.Error(t, err)
	})

	t.Run("validation rejects missing required fields", func(t *testing.T) {
		svc := newTestCarModelService(t, nil, nil)

		_, err := svc.Create(ctx, validation.CarModelCreate{})
		assert.Error(t, err)
	})
}

func TestCarModelServiceGet(t *testing.T) {
	ctx := context.Background()
	modelID := "a0000000-0000-4000-8000-000000000001"

	t.Run("returns model with no images", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		repo.EXPECT().FindByID(ctx, modelID).Return(model.CarModel{ID: modelID, Brand: "BMW"}, nil)

		got, err := svc.Get(ctx, modelID)
		assert.NoError(t, err)
		assert.Equal(t, modelID, got.ID)
	})

	t.Run("populates presigned URL for each image", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		key1, key2 := "models/a.jpg", "models/b.jpg"
		url1, url2 := "https://cdn/a", "https://cdn/b"
		repo.EXPECT().FindByID(ctx, modelID).Return(model.CarModel{
			ID:     modelID,
			Images: []sharedmodel.Image{{Key: key1}, {Key: key2}},
		}, nil)
		storage.EXPECT().GetPresignedURL(ctx, key1).Return(url1, nil)
		storage.EXPECT().GetPresignedURL(ctx, key2).Return(url2, nil)

		got, err := svc.Get(ctx, modelID)
		assert.NoError(t, err)
		assert.Equal(t, url1, got.Images[0].URL)
		assert.Equal(t, url2, got.Images[1].URL)
	})

	t.Run("not found returns ErrCarModelNotFound", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		repo.EXPECT().FindByID(ctx, modelID).Return(model.CarModel{}, model.ErrCarModelNotFound)

		_, err := svc.Get(ctx, modelID)
		assert.ErrorIs(t, err, model.ErrCarModelNotFound)
	})

	t.Run("presigned URL error is propagated", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		key := "models/photo.jpg"
		repo.EXPECT().FindByID(ctx, modelID).Return(model.CarModel{
			Images: []sharedmodel.Image{{Key: key}},
		}, nil)
		storage.EXPECT().GetPresignedURL(ctx, key).Return("", model.ErrSql)

		_, err := svc.Get(ctx, modelID)
		assert.Error(t, err)
	})
}

func TestCarModelServiceGetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty list", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, nil)

		got, err := svc.List(ctx, validation.CarModelFilter{})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("populates presigned URLs", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		key := "models/photo.jpg"
		presigned := "https://cdn/photo"
		repo.EXPECT().Find(ctx, mock.Anything).Return(
			[]model.CarModel{{Images: []sharedmodel.Image{{Key: key}}}}, nil,
		)
		storage.EXPECT().GetPresignedURL(ctx, key).Return(presigned, nil)

		got, err := svc.List(ctx, validation.CarModelFilter{})
		assert.NoError(t, err)
		assert.Equal(t, presigned, got[0].Images[0].URL)
	})

	t.Run("repo error is propagated", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		repo.EXPECT().Find(ctx, mock.Anything).Return(nil, model.ErrSql)

		_, err := svc.List(ctx, validation.CarModelFilter{})
		assert.Error(t, err)
	})
}

func TestCarModelServiceUpdate(t *testing.T) {
	ctx := context.Background()
	modelID := "a0000000-0000-4000-8000-000000000001"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		svc := newTestCarModelService(t, repo, nil)

		repo.EXPECT().Update(ctx, modelID, mock.Anything).Return(nil)

		assert.NoError(t, svc.Update(ctx, modelID, validation.CarModelUpdate{}))
	})

	t.Run("not found returns ErrCarModelNotFound", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		svc := newTestCarModelService(t, repo, nil)

		repo.EXPECT().Update(ctx, modelID, mock.Anything).Return(model.ErrCarModelNotFound)

		assert.ErrorIs(t, svc.Update(ctx, modelID, validation.CarModelUpdate{}), model.ErrCarModelNotFound)
	})
}

func TestCarModelServiceDelete(t *testing.T) {
	ctx := context.Background()
	modelID := "a0000000-0000-4000-8000-000000000001"

	t.Run("happy path delegates to repo", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		svc := newTestCarModelService(t, repo, nil)

		repo.EXPECT().Delete(ctx, modelID).Return(nil)

		assert.NoError(t, svc.Delete(ctx, modelID))
	})

	t.Run("not found returns ErrCarModelNotFound", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		svc := newTestCarModelService(t, repo, nil)

		repo.EXPECT().Delete(ctx, modelID).Return(model.ErrCarModelNotFound)

		assert.ErrorIs(t, svc.Delete(ctx, modelID), model.ErrCarModelNotFound)
	})
}

func TestCarModelServiceGetImageUploadData(t *testing.T) {
	ctx := context.Background()

	t.Run("returns upload data from object storage", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		want := sharedmodel.ImageUploadData{ObjectKey: "models/xyz.jpg", PresignedPutURL: "https://upload.example.com"}
		storage.EXPECT().GetCarModelImageUploadData(ctx).Return(want, nil)

		got, err := svc.GetImageUploadData(ctx)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("storage error is propagated", func(t *testing.T) {
		repo := mocks.NewMockCarModelRepository(t)
		storage := mocks.NewMockObjectStorage(t)
		svc := newTestCarModelService(t, repo, storage)

		storage.EXPECT().GetCarModelImageUploadData(ctx).Return(sharedmodel.ImageUploadData{}, model.ErrSql)

		_, err := svc.GetImageUploadData(ctx)
		assert.Error(t, err)
	})
}
