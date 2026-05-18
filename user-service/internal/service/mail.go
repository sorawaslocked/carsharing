package service

import (
	"context"
	"errors"
	"time"

	"carsharing/user-service/internal/model"
	pkglog "carsharing/user-service/internal/pkg/log"
	"carsharing/user-service/internal/pkg/security"
	"carsharing/user-service/internal/pkg/utils"
)

func (s *UserService) SendActivationCode(ctx context.Context) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "SendActivationCode"), utils.MetadataFromCtx(ctx))

	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			logger.Error("repo: finding user", pkglog.Err(err))
		}
		return err
	}

	code, err := s.activationCodeStorage.Save(ctx, userID)
	if err != nil {
		logger.Error("redis: saving activation code", pkglog.Err(err))
		return err
	}

	if err := s.mailer.SendActivationCode(ctx, user.Email, code); err != nil {
		logger.Error("mailer: sending activation code", pkglog.Err(err))
		return err
	}

	return nil
}

func (s *UserService) CheckActivationCode(ctx context.Context, code string) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CheckActivationCode"), utils.MetadataFromCtx(ctx))

	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			logger.Error("repo: finding user", pkglog.Err(err))
		}
		return err
	}

	if user.IsEmailVerified {
		return model.ErrAlreadyExists
	}

	if err := validateInput(s.validate, &activationCodeValidation{Code: code}); err != nil {
		return err
	}

	codeHash, err := s.activationCodeStorage.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.ValidationErrors{"code": model.ErrInvalidActivationCode}
		}
		logger.Error("redis: getting activation code", pkglog.Err(err))
		return err
	}

	if err := security.CheckStringHash(code, codeHash); err != nil {
		return model.ValidationErrors{"code": model.ErrInvalidActivationCode}
	}

	isEmailVerified := true
	if err := s.userRepo.Update(ctx, userID, model.UserRepoUpdate{
		IsEmailVerified: &isEmailVerified,
		UpdatedAt:       time.Now(),
	}); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			logger.Error("repo: marking email as verified", pkglog.Err(err))
		}
		return err
	}

	if err := s.publisher.PublishUserUpdated(ctx, userID, false); err != nil {
		logger.Error("nats: publishing user.updated", pkglog.Err(err))
	}

	return nil
}
