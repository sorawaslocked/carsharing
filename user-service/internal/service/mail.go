package service

import (
	"context"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/logger"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/security"
	"time"
)

func (s *UserService) CheckActivationCode(ctx context.Context, code string) error {
	id, err := userIDFromCtx(ctx)
	if err != nil {
		return err
	}

	filter := model.UserFilter{
		ID: &id,
	}

	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		return err
	}

	if user.IsActive {
		return model.ErrActivatedUser
	}

	codeValidation := &activationCodeValidation{
		Code: code,
	}

	err = validateInput(s.validate, codeValidation)
	if err != nil {
		return err
	}

	codeHash, err := s.activationCodeStorage.Get(ctx, id)
	if err != nil {
		return err
	}

	err = security.CheckStringHash(code, codeHash)
	if err != nil {
		return model.ValidationErrors{
			"activationCode": model.ErrInvalidActivationCode,
		}
	}

	isActive := true
	update := model.UserUpdate{
		UpdatedAt: time.Now(),
		IsActive:  &isActive,
	}

	err = s.userRepo.Update(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) SendActivationCode(ctx context.Context) error {
	id, err := userIDFromCtx(ctx)
	if err != nil {
		return err
	}

	filter := model.UserFilter{
		ID: &id,
	}

	user, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		return err
	}

	code, err := s.activationCodeStorage.Save(ctx, user.ID)
	if err != nil {
		return err
	}

	err = s.mailer.SendActivationCode(ctx, user.Email, code)
	if err != nil {
		s.log.Error("Mailer", logger.Err(err))

		return err
	}

	return nil
}
