package types

import "github.com/google/uuid"

type UUID = uuid.UUID

func NewUUID() UUID {
	return uuid.New()
}

func ParseUUID(id string) (UUID, error) {
	return uuid.Parse(id)
}
