package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

type DocumentService interface {
	HandleDocumentAnalyzed(ctx context.Context, event model.DocumentAnalyzedEvent) error
}
