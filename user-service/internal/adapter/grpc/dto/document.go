package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	baseuserpb "github.com/sorawaslocked/car-rental-protos/gen/base/user"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
)

func FromCreateDocumentRequest(req *usersvc.CreateDocumentRequest) (validation.DocumentCreate, error) {
	objectKey := req.GetObjectKey()
	if objectKey == "" {
		return validation.DocumentCreate{}, validation.Errors{"object_key": validation.ErrRequiredField}
	}

	imageType := req.GetImageType()
	if _, err := model.ImageTypeFromString(imageType); err != nil {
		return validation.DocumentCreate{}, validation.Errors{"image_type": model.ErrInvalidImageType}
	}

	return validation.DocumentCreate{
		ObjectKey: objectKey,
		ImageType: imageType,
	}, nil
}

func FromGetUploadDocumentDataRequest(req *usersvc.GetUploadDocumentDataRequest) (model.ImageType, error) {
	imageType, err := model.ImageTypeFromString(req.GetImageType())
	if err != nil {
		return "", validation.Errors{"image_type": model.ErrInvalidImageType}
	}
	return imageType, nil
}

func FromCheckDocumentRequest(req *usersvc.CheckDocumentRequest) (docID string, data validation.DocumentUpdate, err error) {
	docID = req.GetDocId()

	statusStr := req.GetStatus()
	if _, err = model.DocumentStatusFromString(statusStr); err != nil {
		return "", validation.DocumentUpdate{}, validation.Errors{"status": model.ErrInvalidDocumentStatus}
	}

	data = validation.DocumentUpdate{Status: statusStr}
	if e := req.GetError(); e != "" {
		data.Error = &e
	}
	return docID, data, nil
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
