package handler

import (
	"carsharing/api-gateway/internal/adapter/http/dto"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CreateDocument"), utils.MetadataFromCtx(c))

	imageType, objectKey, err := dto.FromCreateDocumentRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	id, err := h.svc.CreateDocument(c, objectKey, imageType)
	if err != nil {
		log.Warn("creating document", pkglog.Err(err))

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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetProfileImageUploadData"), utils.MetadataFromCtx(c))

	data, err := h.svc.GetProfileImageUploadData(c)
	if err != nil {
		log.Warn("getting profile image upload data", pkglog.Err(err))

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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetUploadDocumentData"), utils.MetadataFromCtx(c))

	imageType, err := dto.FromGetUploadDocumentDataRequest(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	data, err := h.svc.GetUploadDocumentData(c, imageType)
	if err != nil {
		log.Warn("getting document upload data", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.Ok(c, gin.H{"uploadData": dto.ToImageUploadDataResponse(data)})
}

// ListDocuments godoc
// @Summary      List documents for a user
// @Description  Returns documents for the specified user, optionally filtered by status and image type, with pagination and sorting support.
// @Tags         documents
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      string   true   "User ID"
// @Param        status     query     string   false  "Filter by status"     Enums(pending, processed, approved, rejected)
// @Param        imageType  query     string   false  "Filter by image type" Enums(id_front, id_back, driving_license_front, driving_license_back)
// @Param        sort       query     string   false  "Sort order"           Enums(+createdAt, -createdAt)
// @Param        limit      query     integer  false  "Pagination limit"
// @Param        offset     query     integer  false  "Pagination offset"
// @Success      200        {object}  dto.DocumentsResponse
// @Failure      400        {object}  dto.ErrorResponse
// @Failure      401        {object}  dto.ErrorResponse
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /users/{id}/documents [get]
func (h *UserHandler) ListDocuments(c *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "ListDocuments"), utils.MetadataFromCtx(c))

	filter, err := dto.DocumentFilterFromCtx(c)
	if err != nil {
		dto.FromError(c, err)

		return
	}

	documents, err := h.svc.ListDocuments(c, filter)
	if err != nil {
		log.Warn("listing documents", pkglog.Err(err))

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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CheckDocument"), utils.MetadataFromCtx(c))

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

	if err = h.svc.CheckDocument(c, docID, status, documentError); err != nil {
		log.Warn("checking document", pkglog.Err(err))

		dto.FromError(c, err)

		return
	}

	dto.NoContent(c)
}
