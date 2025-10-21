package validate

import (
	"github.com/go-playground/validator/v10"
	"testing"
	"time"
)

func TestMinAge(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("min_age_18", MinAge(18))

	tests := []struct {
		name      string
		birthDate time.Time
		wantValid bool
	}{
		{
			name:      "exactly 18 years old today",
			birthDate: time.Now().AddDate(-18, 0, 0),
			wantValid: true,
		},
		{
			name:      "18 years and 1 day old",
			birthDate: time.Now().AddDate(-18, 0, -1),
			wantValid: true,
		},
		{
			name:      "17 years and 364 days old",
			birthDate: time.Now().AddDate(-18, 0, 1),
			wantValid: false,
		},
		{
			name:      "17 years old",
			birthDate: time.Now().AddDate(-17, 0, 0),
			wantValid: false,
		},
		{
			name:      "25 years old",
			birthDate: time.Now().AddDate(-25, 0, 0),
			wantValid: true,
		},
		{
			name:      "100 years old",
			birthDate: time.Now().AddDate(-100, 0, 0),
			wantValid: true,
		},
		{
			name:      "0 years old - born today",
			birthDate: time.Now(),
			wantValid: false,
		},
		{
			name:      "specific date - Feb 29, 2000 (should be 24 years old in 2024)",
			birthDate: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.birthDate, "min_age_18")
			isValid := err == nil

			if isValid != tt.wantValid {
				age := time.Now().Year() - tt.birthDate.Year()
				t.Errorf("minAge(18) for birthDate=%v (approx age %d) = %v, want %v",
					tt.birthDate.Format("2006-01-02"), age, isValid, tt.wantValid)
				if err != nil {
					t.Logf("Error: %v", err)
				}
			}
		})
	}
}

func TestComplexPassword(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("complex_password", ComplexPassword)

	tests := []struct {
		name      string
		password  string
		wantValid bool
		reason    string
	}{
		{
			name:      "valid password with all requirements",
			password:  "Password123!",
			wantValid: true,
			reason:    "has uppercase, lowercase, number, and special char",
		},
		{
			name:      "valid password with multiple special chars",
			password:  "P@ssw0rd#2024",
			wantValid: true,
			reason:    "has all requirements",
		},
		{
			name:      "missing uppercase letter",
			password:  "password123!",
			wantValid: false,
			reason:    "no uppercase letter",
		},
		{
			name:      "missing lowercase letter",
			password:  "PASSWORD123!",
			wantValid: false,
			reason:    "no lowercase letter",
		},
		{
			name:      "missing number",
			password:  "Password!",
			wantValid: false,
			reason:    "no number",
		},
		{
			name:      "missing special character",
			password:  "Password123",
			wantValid: false,
			reason:    "no special character",
		},
		{
			name:      "only uppercase letters",
			password:  "PASSWORD",
			wantValid: false,
			reason:    "missing lowercase, number, and special char",
		},
		{
			name:      "only lowercase letters",
			password:  "password",
			wantValid: false,
			reason:    "missing uppercase, number, and special char",
		},
		{
			name:      "only numbers",
			password:  "12345678",
			wantValid: false,
			reason:    "missing uppercase, lowercase, and special char",
		},
		{
			name:      "only special characters",
			password:  "!@#$%^&*",
			wantValid: false,
			reason:    "missing uppercase, lowercase, and number",
		},
		{
			name:      "empty password",
			password:  "",
			wantValid: false,
			reason:    "empty string",
		},
		{
			name:      "valid with brackets",
			password:  "Abcd123[]",
			wantValid: true,
			reason:    "has all requirements",
		},
		{
			name:      "valid with underscore",
			password:  "Pass_word1",
			wantValid: true,
			reason:    "has all requirements",
		},
		{
			name:      "valid with hyphen",
			password:  "Pass-word1",
			wantValid: true,
			reason:    "has all requirements",
		},
		{
			name:      "minimum valid password",
			password:  "aA1!",
			wantValid: true,
			reason:    "shortest possible valid password",
		},
		{
			name:      "long valid password",
			password:  "ThisIsAVeryLongPassword123!WithManyCharacters",
			wantValid: true,
			reason:    "long password with all requirements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.password, "complex_password")
			isValid := err == nil

			if isValid != tt.wantValid {
				t.Errorf("complexPassword(%q) = %v, want %v\nReason: %s",
					tt.password, isValid, tt.wantValid, tt.reason)
				if err != nil {
					t.Logf("Error: %v", err)
				}
			}
		})
	}
}
