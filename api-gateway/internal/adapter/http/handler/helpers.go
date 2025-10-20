package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ok(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusOK, body)
}

func created(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusCreated, body)
}

func badRequest(ctx *gin.Context, body any) {
	ctx.JSON(http.StatusBadRequest, body)
}

func internalServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "something went wrong",
	})
}

func malformedJson(ctx *gin.Context) {
	errors := make(map[string]string)
	errors["json"] = "malformed json body"

	badRequest(ctx, errors)
}
