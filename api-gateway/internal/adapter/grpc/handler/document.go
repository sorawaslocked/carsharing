package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/utils"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *UserHandler) CreateDocument(ctx context.Context, objectKey, imageType string) (string, error) {
	logger := pkglog.WithMethod(h.log, "CreateDocument")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.CreateDocument(ctx, &usersvc.CreateDocumentRequest{
		ObjectKey: objectKey,
		ImageType: imageType,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *UserHandler) GetDocumentImageUploadData(ctx context.Context, imageType string) (model.ImageUploadData, error) {
	logger := pkglog.WithMethod(h.log, "GetDocumentImageUploadData")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetUploadDocumentData(ctx, &usersvc.GetUploadDocumentDataRequest{ImageType: imageType})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}

func (h *UserHandler) GetProfileImageUploadData(ctx context.Context) (model.ImageUploadData, error) {
	logger := pkglog.WithMethod(h.log, "GetProfileImageUploadData")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetProfileImageUploadData(ctx, &emptypb.Empty{})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}

func (h *UserHandler) GetProcessedDocumentsForUser(ctx context.Context, userID string) ([]model.Document, error) {
	logger := pkglog.WithMethod(h.log, "GetProcessedDocumentsForUser")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetProcessedDocumentsForUser(ctx, &usersvc.GetProcessedDocumentsForUserRequest{UserId: userID})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	docs := make([]model.Document, len(res.GetDocuments()))
	for i, d := range res.GetDocuments() {
		docs[i] = dto.DocumentFromProto(d)
	}

	return docs, nil
}

func (h *UserHandler) CheckDocument(ctx context.Context, docID, status string, reason *string) error {
	logger := pkglog.WithMethod(h.log, "CheckDocument")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.CheckDocument(ctx, &usersvc.CheckDocumentRequest{
		DocId:  docID,
		Status: status,
		Error:  reason,
	})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}
