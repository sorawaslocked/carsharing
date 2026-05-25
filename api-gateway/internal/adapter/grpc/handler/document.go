package handler

import (
	"context"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	usersvc "carsharing/protos/gen/service/user"
	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *UserHandler) CreateDocument(ctx context.Context, objectKey, imageType string) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CreateDocument"), utils.MetadataFromCtx(ctx))

	res, err := h.client.CreateDocument(ctx, &usersvc.CreateDocumentRequest{
		ObjectKey: objectKey,
		ImageType: imageType,
	})
	if err != nil {
		log.Warn("creating document", pkglog.Err(err))

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *UserHandler) GetDocumentImageUploadData(ctx context.Context, imageType string) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetDocumentImageUploadData"), utils.MetadataFromCtx(ctx))

	res, err := h.client.GetUploadDocumentData(ctx, &usersvc.GetUploadDocumentDataRequest{ImageType: imageType})
	if err != nil {
		log.Warn("getting document image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}

func (h *UserHandler) GetProfileImageUploadData(ctx context.Context) (sharedmodel.ImageUploadData, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetProfileImageUploadData"), utils.MetadataFromCtx(ctx))

	res, err := h.client.GetProfileImageUploadData(ctx, &emptypb.Empty{})
	if err != nil {
		log.Warn("getting profile image upload data", pkglog.Err(err))

		return sharedmodel.ImageUploadData{}, dto.FromGrpcErr(err)
	}

	return dto.ImageUploadDataFromProto(res.GetUploadData()), nil
}

func (h *UserHandler) ListDocuments(ctx context.Context, filter model.DocumentFilter) ([]model.Document, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "ListDocuments"), utils.MetadataFromCtx(ctx))

	req := &usersvc.ListDocumentsRequest{
		UserId:    filter.UserID,
		Status:    filter.Status,
		ImageType: filter.ImageType,
		Sort:      filter.Sort,
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListDocuments(ctx, req)
	if err != nil {
		log.Warn("listing documents", pkglog.Err(err))

		return nil, dto.FromGrpcErr(err)
	}

	docs := make([]model.Document, len(res.GetDocuments()))
	for i, d := range res.GetDocuments() {
		docs[i] = dto.DocumentFromProto(d)
	}

	return docs, nil
}

func (h *UserHandler) CheckDocument(ctx context.Context, docID, status string, reason *string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CheckDocument"), utils.MetadataFromCtx(ctx))

	_, err := h.client.CheckDocument(ctx, &usersvc.CheckDocumentRequest{
		DocId:  docID,
		Status: status,
		Error:  reason,
	})
	if err != nil {
		log.Warn("checking document", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}
