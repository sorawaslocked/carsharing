package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type DocumentsResponse struct {
	Documents []Document `json:"documents"`
}

type Document struct {
	ID        string  `json:"id"`
	UserID    string  `json:"userID"`
	ImageType string  `json:"imageType"`
	Status    string  `json:"status"`
	Reason    *string `json:"reason,omitempty"`
	ImageURL  string  `json:"imageURL"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetUploadDocumentDataRequest struct {
	ImageType string `json:"imageType"`
}

type CreateDocumentRequest struct {
	ImageType string `json:"imageType"`
	ObjectKey string `json:"objectKey"`
}

type CheckDocumentRequest struct {
	Status string  `json:"status"`
	Error  *string `json:"error"`
}

func ToDocumentResponse(m model.Document) Document {
	return Document{
		ID:        m.ID,
		UserID:    m.UserID,
		ImageType: m.ImageType,
		Status:    m.Status,
		Reason:    m.Reason,
		ImageURL:  m.ImageURL,
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
