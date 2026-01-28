package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	"net/http"
)

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
			"message": "invalid query params",
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
