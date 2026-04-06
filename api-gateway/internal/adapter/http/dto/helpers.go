package dto

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ImageUploadData struct {
	PresignedUrl string `json:"presignedUrl"`
	ObjectKey    string `json:"objectKey"`
}

func ToImageUploadDataResponse(m model.ImageUploadData) ImageUploadData {
	return ImageUploadData{
		PresignedUrl: m.PresignedUrl,
		ObjectKey:    m.ObjectKey,
	}
}

func Ok(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusOK, body)
}

func Created(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusCreated, body)
}

func NoContent(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, nil)
}

func badRequest(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusBadRequest, body)
}

func MalformedJson(ctx *gin.Context) {
	md := make(map[string]any)
	md["type"] = "json"
	body := errorBody("malformed json format", md)

	badRequest(ctx, body)
}

func InvalidQueryParams(ctx *gin.Context) {
	body := gin.H{
		"error": gin.H{
			"message": model.ErrInvalidQueryParam.Error(),
		},
	}

	badRequest(ctx, body)
}

func EmptyIDParam(ctx *gin.Context) {
	body := gin.H{
		"error": gin.H{
			"message": model.ErrEmptyIDParam.Error(),
		},
	}

	badRequest(ctx, body)
}

func validationError(ctx *gin.Context, ve model.ValidationErrors) {
	md := make(map[string]any)
	md["type"] = "validation"
	md["validation"] = ve
	body := errorBody("validation error", md)

	badRequest(ctx, body)
}

func unauthorized(ctx *gin.Context) {
	body := errorBody("not authorized", nil)

	ctx.JSON(http.StatusUnauthorized, body)
}

func forbidden(ctx *gin.Context) {
	body := errorBody("insufficient permissions", nil)

	ctx.JSON(http.StatusForbidden, body)
}

func notFound(ctx *gin.Context) {
	body := errorBody("resource not found", nil)

	ctx.JSON(http.StatusNotFound, body)
}

func conflict(ctx *gin.Context) {
	body := errorBody("resource already exists", nil)

	ctx.JSON(http.StatusConflict, body)
}

func internalServerError(ctx *gin.Context) {
	body := errorBody("something went wrong", nil)

	ctx.JSON(http.StatusInternalServerError, body)
}

func pagination(ctx *gin.Context) (*model.Pagination, error) {
	var p model.Pagination
	paginationEmpty := true

	if v := ctx.Query("limit"); v != "" {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return nil, model.ErrInvalidQueryParam
		}

		p.Limit = int64(vInt)
		paginationEmpty = false
	}
	if v := ctx.Query("offset"); v != "" {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return nil, model.ErrInvalidQueryParam
		}

		p.Offset = int64(vInt)
		paginationEmpty = false
	}

	if paginationEmpty {
		return nil, nil
	}

	return &p, nil
}

func IDParam(ctx *gin.Context) (string, error) {
	id := ctx.Param("id")

	if id == "" {
		return "", model.ErrEmptyIDParam
	}

	return id, nil
}
