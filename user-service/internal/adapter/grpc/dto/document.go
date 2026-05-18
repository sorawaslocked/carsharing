package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"carsharing/user-service/internal/model"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	baseuserpb "github.com/sorawaslocked/car-rental-protos/gen/base/user"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
)

func FromCreateDocumentRequest(req *usersvc.CreateDocumentRequest) (objectKey string, imageType model.ImageType, err error) {
	objectKey = req.GetObjectKey()
	if objectKey == "" {
		return "", "", model.ValidationErrors{"object_key": model.ErrRequiredField}
	}

	imageType, err = model.ImageTypeFromString(req.GetImageType())
	if err != nil {
		return "", "", model.ValidationErrors{"image_type": model.ErrInvalidImageType}
	}
	return objectKey, imageType, nil
}

func FromGetUploadDocumentDataRequest(req *usersvc.GetUploadDocumentDataRequest) (model.ImageType, error) {
	imageType, err := model.ImageTypeFromString(req.GetImageType())
	if err != nil {
		return "", model.ValidationErrors{"image_type": model.ErrInvalidImageType}
	}
	return imageType, nil
}

func FromCheckDocumentRequest(req *usersvc.CheckDocumentRequest) (docID string, status model.DocumentStatus, docError *string, err error) {
	docID = req.GetDocId()

	status, err = model.DocumentStatusFromString(req.GetStatus())
	if err != nil {
		return "", "", nil, model.ValidationErrors{"status": model.ErrInvalidDocumentStatus}
	}

	if e := req.GetError(); e != "" {
		docError = &e
	}
	return
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
	if doc.Image != nil && doc.Image.URL != "" {
		d.ImageUrl = doc.Image.URL
	}

	return d
}

func ImageUploadDataToProto(data model.ImageUploadData) *basepb.ImageUploadData {
	return &basepb.ImageUploadData{
		PresignedPutUrl: data.PresignedPutURL,
		ObjectKey:       data.ObjectKey,
	}
}
