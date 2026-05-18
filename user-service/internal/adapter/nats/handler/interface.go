package handler

import (
	"context"

	"carsharing/user-service/internal/model"
)

type DocumentService interface {
	HandleDocumentAnalyzed(ctx context.Context, event model.DocumentAnalyzedEvent) error
}
