package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"

	basepb "carsharing/protos/gen/base"
	baseuserpb "carsharing/protos/gen/base/user"
	usersvc "carsharing/protos/gen/service/user"
)

func FromCreateDocumentRequest(req *usersvc.CreateDocumentRequest) validation.DocumentCreate {
	return validation.DocumentCreate{
		ObjectKey: req.GetObjectKey(),
		ImageType: req.GetImageType(),
	}
}

func FromGetUploadDocumentDataRequest(req *usersvc.GetUploadDocumentDataRequest) string {
	return req.GetImageType()
}

func FromCheckDocumentRequest(req *usersvc.CheckDocumentRequest) (string, validation.DocumentUpdate) {
	data := validation.DocumentUpdate{
		Status: req.GetStatus(),
		Error:  req.Error,
	}

	return req.GetDocId(), data
}

func DocumentToProto(doc model.Document) *baseuserpb.Document {
	d := &baseuserpb.Document{
		Id:        doc.ID,
		UserId:    doc.UserID,
		ImageType: doc.ImageType.String(),
		Status:    doc.Status.String(),
		CreatedAt: timestamppb.New(doc.CreatedAt),
		UpdatedAt: timestamppb.New(doc.UpdatedAt),
	}

	if doc.Error != nil {
		d.Error = doc.Error
	}
	if doc.Image.URL != "" {
		d.ImageUrl = doc.Image.URL
	}

	return d
}

func ImageUploadDataToProto(data sharedmodel.ImageUploadData) *basepb.ImageUploadData {
	return &basepb.ImageUploadData{
		PresignedPutUrl: data.PresignedPutURL,
		ObjectKey:       data.ObjectKey,
	}
}
