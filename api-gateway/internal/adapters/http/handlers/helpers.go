package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func malformedJson(ctx *gin.Context) {
	errors := make(map[string]string)
	errors["json"] = "malformed json body"

	badRequest(ctx, errors)
}

func badRequest(ctx *gin.Context, errors map[string]string) {
	ctx.JSON(http.StatusBadRequest, gin.H{"errors": errors})
}

func ok(ctx *gin.Context, body map[string]any) {
	ctx.JSON(http.StatusOK, body)
}
