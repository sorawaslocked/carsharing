package service

import (
	"context"

	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
)

func (s *UserService) StreamDocumentAnalyzed(ctx context.Context, userID *string, passed *bool, send func(model.DocumentAnalyzedEvent) error) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "StreamDocumentAnalyzed"), utils.MetadataFromCtx(ctx))

	if err := s.presenter.StreamDocumentAnalyzed(ctx, userID, passed, send); err != nil {
		log.Warn("streaming document analyzed", pkglog.Err(err))
		return err
	}

	return nil
}
