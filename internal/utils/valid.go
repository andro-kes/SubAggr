package utils

import (
	"time"

	"github.com/google/uuid"
)

func IsValidUUID(id uuid.UUID) bool {
    _, err := uuid.Parse(id.String())
    return err == nil
}

func IsValidDate(date string) bool {
    _, err := time.Parse("01-2006", date)
    return err == nil
}