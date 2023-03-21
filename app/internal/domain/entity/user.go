package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash []byte
	Created      time.Time
}

type Image struct {
	Name        string
	ContentType string
	URL         string
	Size        int64
	Created     time.Time
}
