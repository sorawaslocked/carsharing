package service

import (
	"context"
	"errors"
	"time"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/pkg/security"
	"carsharing/user-service/internal/validation"
)

func (s *UserService) SendActivationCode(ctx context.Context) error {
	md := utils.MetadataFromCtx(ctx)
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "SendActivationCode"), md)

	if md.UserID == nil {
		return model.ErrUnauthenticated
	}
	userID := *md.UserID

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
		logger.Error("brevo: sending activation code", pkglog.Err(err))
		return err
	}

	return nil
}

func (s *UserService) CheckActivationCode(ctx context.Context, code string) error {
	md := utils.MetadataFromCtx(ctx)
	logger := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CheckActivationCode"), md)

	if md.UserID == nil {
		return model.ErrUnauthenticated
	}
	userID := *md.UserID

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

	if err := validation.ValidateActivationCode(s.validate, code); err != nil {
		return err
	}

	codeHash, err := s.activationCodeStorage.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return validation.Errors{"code": validation.ErrInvalidActivationCode}
		}
		logger.Error("redis: getting activation code", pkglog.Err(err))
		return err
	}

	if err := security.CheckStringHash(code, codeHash); err != nil {
		return validation.Errors{"code": validation.ErrInvalidActivationCode}
	}

	isEmailVerified := true
	if err := s.userRepo.Update(ctx, userID, model.UserUpdate{
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
