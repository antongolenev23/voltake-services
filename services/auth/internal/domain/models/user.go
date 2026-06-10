package models

import "github.com/google/uuid"

type User struct {
	Email string
	PassHash []byte
	ID uuid.UUID
	IsAdmin bool
}