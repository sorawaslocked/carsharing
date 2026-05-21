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
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "SendActivationCode"), md)

	if md.UserID == nil {
		return model.ErrUnauthenticated
	}

	user, err := s.userRepo.FindByID(ctx, *md.UserID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: finding user", pkglog.Err(err))
		}
		return err
	}

	code, err := s.activationCodeStorage.Save(ctx, *md.UserID)
	if err != nil {
		log.Error("cache: saving activation code", pkglog.Err(err))

		return err
	}

	if err := s.mailer.SendActivationCode(ctx, user.Email, code); err != nil {
		log.Error("mailer: sending activation code", pkglog.Err(err))

		return err
	}

	return nil
}

func (s *UserService) CheckActivationCode(ctx context.Context, code string) error {
	md := utils.MetadataFromCtx(ctx)
	log := pkglog.WithMetadata(pkglog.WithMethod(s.log, "CheckActivationCode"), md)

	if md.UserID == nil {
		return model.ErrUnauthenticated
	}

	user, err := s.userRepo.FindByID(ctx, *md.UserID)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: finding user", pkglog.Err(err))
		}
		return err
	}

	if user.IsEmailVerified {
		return model.ErrEmailVerified
	}

	if err := validation.ValidateActivationCode(s.validate, code); err != nil {
		return err
	}

	codeHash, err := s.activationCodeStorage.Get(ctx, *md.UserID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return validation.Errors{"code": validation.ErrInvalidActivationCode}
		}
		log.Error("cache: getting activation code", pkglog.Err(err))

		return err
	}

	if err := security.CheckStringHash(code, codeHash); err != nil {
		return validation.Errors{"code": validation.ErrInvalidActivationCode}
	}

	isEmailVerified := true
	if err := s.userRepo.Update(ctx, *md.UserID, model.UserUpdate{
		IsEmailVerified: &isEmailVerified,
		UpdatedAt:       time.Now(),
	}); err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			log.Error("repo: updating user email verified", pkglog.Err(err))
		}
		return err
	}

	if err := s.publisher.PublishUserUpdated(ctx, *md.UserID, false); err != nil {
		log.Error("event: publishing user updated", pkglog.Err(err))
	}

	return nil
}
