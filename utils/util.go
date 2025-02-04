package utils

import "github.com/google/uuid"

func UUIDToString(id uuid.UUID) string {
	return id.String()
}