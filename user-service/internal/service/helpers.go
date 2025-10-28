package service

import (
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"strings"
)

func uncapitalize(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func toRoleStrings(roles []model.Role) []string {
	var result []string
	for _, role := range roles {
		result = append(result, role.String())
	}

	return result
}
