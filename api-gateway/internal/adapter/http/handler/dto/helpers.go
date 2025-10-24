package dto

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ok(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusOK, body)
}

func Created(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusCreated, body)
}

func badRequest(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusBadRequest, body)
}

func MalformedJson(ctx *gin.Context) {
	body := errorBody("malformed json format")
	body["type"] = "json"

	badRequest(ctx, body)
}

func unauthorized(ctx *gin.Context) {
	body := errorBody("not authorized")

	ctx.JSON(http.StatusUnauthorized, body)
}

func forbidden(ctx *gin.Context) {
	body := errorBody("insufficient permissions")

	ctx.JSON(http.StatusForbidden, body)
}

func notFound(ctx *gin.Context) {
	body := errorBody("resource not found")

	ctx.JSON(http.StatusNotFound, body)
}

func conflict(ctx *gin.Context) {
	body := errorBody("resource already exists")

	ctx.JSON(http.StatusConflict, body)
}

func internalServerError(ctx *gin.Context) {
	body := errorBody("something went wrong")

	ctx.JSON(http.StatusInternalServerError, body)
}
