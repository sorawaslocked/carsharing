package validation

import (
	"fmt"
	"log/slog"
	"regexp"
	"time"

	sharedmodel "carsharing/shared/model"
	"carsharing/user-service/internal/model"

	"github.com/go-playground/validator/v10"
)

var (
	hasUpperRX   = regexp.MustCompile(`[A-Z]`)
	hasLowerRX   = regexp.MustCompile(`[a-z]`)
	hasNumberRX  = regexp.MustCompile(`[0-9]`)
	hasSpecialRX = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>?~]`)
)

type ErrRegisterValidator struct {
	Tag string
}

func (e ErrRegisterValidator) Error() string {
	return fmt.Sprintf("failed to register validator %q", e.Tag)
}

func RegisterCustomValidators(v *validator.Validate, log *slog.Logger) error {
	validators := []struct {
		tag string
		fn  validator.Func
	}{
		{"min_age", minAge},
		{"complex_password", complexPassword},
		{"document_image_type", imageTypeValidator},
		{"document_status", documentStatusValidator},
		{"role", roleValidator},
	}

	for _, vd := range validators {
		if err := v.RegisterValidation(vd.tag, vd.fn); err != nil {
			log.Error("registering validator", slog.String("tag", vd.tag), slog.Any("error", err))
			return ErrRegisterValidator{Tag: vd.tag}
		}
	}

	return nil
}

func minAge(fl validator.FieldLevel) bool {
	minAgeParam := fl.Param()
	minAgeVal := 18

	if minAgeParam != "" {
		fmt.Sscanf(minAgeParam, "%d", &minAgeVal)
	}

	birthDate := fl.Field().Interface().(time.Time)
	now := time.Now()

	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age >= minAgeVal
}

func complexPassword(fl validator.FieldLevel) bool {
	p := fl.Field().String()
	return hasUpperRX.MatchString(p) &&
		hasLowerRX.MatchString(p) &&
		hasNumberRX.MatchString(p) &&
		hasSpecialRX.MatchString(p)
}

func imageTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.DocumentImageTypeFromString(fl.Field().String())
	return ok
}

func documentStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.DocumentStatusFromString(fl.Field().String())
	return ok
}

func roleValidator(fl validator.FieldLevel) bool {
	_, ok := sharedmodel.RoleFromString(fl.Field().String())
	return ok
}
