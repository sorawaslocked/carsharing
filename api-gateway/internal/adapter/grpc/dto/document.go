package dto

import (
	"carsharing/api-gateway/internal/model"
	baseuser "carsharing/protos/gen/base/user"
)

func DocumentFromProto(d *baseuser.Document) model.Document {
	return model.Document{
		ID:        d.GetId(),
		UserID:    d.GetUserId(),
		ImageType: d.GetImageType(),
		Status:    d.GetStatus(),
		Reason:    d.Error,
		Image:     ImageFromProto(d.GetImage()),
		CreatedAt: d.GetCreatedAt().AsTime(),
		UpdatedAt: d.GetUpdatedAt().AsTime(),
	}
}
