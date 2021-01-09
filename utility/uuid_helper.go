package utility

import (
	"github.com/google/uuid"
	"strings"
)

func GetUUIDWithoutDash() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
