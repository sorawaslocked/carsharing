package validate

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

var (
	hasUpperRX   = regexp.MustCompile(`[A-Z]`)
	hasLowerRX   = regexp.MustCompile(`[a-z]`)
	hasNumberRX  = regexp.MustCompile(`[0-9]`)
	hasSpecialRX = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]`)
)

func MinAge(minAge int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		birthDate := fl.Field().Interface().(time.Time)
		now := time.Now()

		age := now.Year() - birthDate.Year()

		if now.Month() < birthDate.Month() ||
			(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
			age--
		}

		return age >= minAge
	}
}

func ComplexPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	hasUpper := hasUpperRX.MatchString(password)
	hasLower := hasLowerRX.MatchString(password)
	hasNumber := hasNumberRX.MatchString(password)
	hasSpecial := hasSpecialRX.MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}
