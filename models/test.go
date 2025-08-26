package models

import (
	"time"
)

type Test struct {
	ID        int       `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateTestRequest struct {
	Message string `json:"message" binding:"required"`
}
