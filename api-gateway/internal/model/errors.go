package model

import (
	"bytes"
	"fmt"
	"strings"
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
