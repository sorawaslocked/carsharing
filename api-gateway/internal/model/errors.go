package model

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("resource not found")
	ErrAlreadyExists       = errors.New("resource already exists")
	ErrInternalServerError = errors.New("something went wrong")
)

type ValidationErrors map[string]string

func (ve ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")

	for field, err := range ve {
		buff.WriteString(fmt.Sprintf("%s: %s", field, err))
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
