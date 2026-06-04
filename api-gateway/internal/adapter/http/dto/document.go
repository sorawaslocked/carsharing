package dto

import (
	"time"

	"carsharing/api-gateway/internal/model"
	"github.com/gin-gonic/gin"
)

type DocumentsResponse struct {
	Documents []Document `json:"documents"`
}

type Document struct {
	ID        string  `json:"id"`
	UserID    string  `json:"userID"`
	ImageType string  `json:"imageType" validate:"oneof=id_front id_back driving_license_front driving_license_back"`
	Status    string  `json:"status" validate:"oneof=pending processed approved rejected"`
	Reason    *string `json:"reason,omitempty"`
	Image     Image   `json:"image"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetUploadDocumentDataRequest struct {
	ImageType string `json:"imageType" binding:"required,oneof=id_front id_back driving_license_front driving_license_back"`
}

type CreateDocumentRequest struct {
	ImageType string `json:"imageType" binding:"required,oneof=id_front id_back driving_license_front driving_license_back"`
	ObjectKey string `json:"objectKey" binding:"required"`
}

type CheckDocumentRequest struct {
	Status string  `json:"status" binding:"required,oneof=pending processed approved rejected"`
	Error  *string `json:"error"`
}

func ToDocumentResponse(m model.Document) Document {
	return Document{
		ID:        m.ID,
		UserID:    m.UserID,
		ImageType: m.ImageType,
		Status:    m.Status,
		Reason:    m.Reason,
		Image:     toImage(m.Image),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromCreateDocumentRequest(c *gin.Context) (imageType, objectKey string, err error) {
	var req CreateDocumentRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		return "", "", err
	}

	return req.ImageType, req.ObjectKey, nil
}

func FromGetUploadDocumentDataRequest(c *gin.Context) (imageType string, err error) {
	var req GetUploadDocumentDataRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		return "", err
	}

	return req.ImageType, nil
}

func FromCheckDocumentRequest(c *gin.Context) (status string, documentError *string, err error) {
	var req CheckDocumentRequest

	if err = c.ShouldBindJSON(&req); err != nil {
		return "", nil, nil
	}

	return req.Status, req.Error, nil
}

func DocumentFilterFromCtx(c *gin.Context) (model.DocumentFilter, error) {
	userID, err := IDParam(c)
	if err != nil {
		return model.DocumentFilter{}, err
	}

	f := model.DocumentFilter{UserID: userID}

	if v := c.Query("status"); v != "" {
		f.Status = &v
	}
	if v := c.Query("imageType"); v != "" {
		f.ImageType = &v
	}
	if v := c.Query("sort"); v != "" {
		f.Sort = &v
	}

	p, err := pagination(c)
	if err != nil {
		return model.DocumentFilter{}, model.ErrInvalidQueryParam
	}
	f.Pagination = p

	return f, nil
}
