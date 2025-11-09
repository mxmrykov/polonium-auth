package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NewSession() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
