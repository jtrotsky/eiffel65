package util

import "github.com/google/uuid"

// GenerateUUID creates a new random unique ID.
func GenerateUUID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}
