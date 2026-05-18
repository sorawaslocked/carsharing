package handler

import (
	"carsharing/api-gateway/internal/adapter/http/dto"
	"github.com/gin-gonic/gin"
)

// CreateDocument godoc
// @Summary      Create document record
// @Description  Creates a document record after the image has been uploaded to object storage.
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateDocumentRequest  true  "Document payload"
// @Success      201   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/documents [post]
func (h *UserHandler) CreateDocument(c *gin.Context) {
	imageType, objectKey, err := dto.FromCreateDocumentRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	id, err := h.svc.CreateDocument(c, objectKey, imageType)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.Created(c, gin.H{"id": id})
}

// GetProfileImageUploadData godoc
// @Summary      Get profile image upload URL
// @Description  Returns a presigned PUT URL and object key for uploading the current user's profile image to object storage.
// @Tags         documents
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ImageUploadResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /users/profile/image-upload [get]
func (h *UserHandler) GetProfileImageUploadData(c *gin.Context) {
	data, err := h.svc.GetProfileImageUploadData(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.Ok(c, gin.H{"uploadData": dto.ToImageUploadDataResponse(data)})
}

// GetUploadDocumentData godoc
// @Summary      Get document upload URL
// @Description  Returns a presigned PUT URL for uploading a document image to object storage.
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.GetUploadDocumentDataRequest  true  "Image type"
// @Success      200   {object}  dto.ImageUploadResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/documents/upload [post]
func (h *UserHandler) GetUploadDocumentData(c *gin.Context) {
	imageType, err := dto.FromGetUploadDocumentDataRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	data, err := h.svc.GetUploadDocumentData(c, imageType)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.Ok(c, gin.H{"uploadData": dto.ToImageUploadDataResponse(data)})
}

// GetProcessedDocumentsForUser godoc
// @Summary      Get processed documents for a user
// @Description  Returns all processed documents belonging to the specified user.
// @Tags         documents
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string  true  "User ID"
// @Success      200   {object}  dto.DocumentsResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/{id}/documents/processed [get]
func (h *UserHandler) GetProcessedDocumentsForUser(c *gin.Context) {
	userID, err := dto.IDParam(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	documents, err := h.svc.GetProcessedDocumentsForUser(c, userID)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	documentResponse := make([]dto.Document, len(documents))
	for i, doc := range documents {
		documentResponse[i] = dto.ToDocumentResponse(doc)
	}

	dto.Ok(c, gin.H{"documents": documentResponse})
}

// CheckDocument godoc
// @Summary      Review a document
// @Description  Sets the review status of a document (e.g. approved or rejected) and optionally provides a rejection reason.
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                    true  "Document ID"
// @Param        body  body      dto.CheckDocumentRequest  true  "Review payload"
// @Success      204
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /users/documents/check/{id} [post]
func (h *UserHandler) CheckDocument(c *gin.Context) {
	docID, err := dto.IDParam(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	status, documentError, err := dto.FromCheckDocumentRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	err = h.svc.CheckDocument(c, docID, status, documentError)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}
