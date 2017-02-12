package utils

import (
	"github.com/satori/go.uuid"
)

// CreateCorrelationID generates a new uuid
func CreateCorrelationID() uuid.UUID {
	return uuid.NewV4()
}
