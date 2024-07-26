package models

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusProcessing Status = "PROCESSING"
	StatusCompleted  Status = "COMPLETED"
)

type Message struct {
	ID        *uuid.UUID
	Message   *string
	Status    *Status
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type CreateRequest struct {
	Message *string `json:"message,omitempty"`
}

type FilterRequest struct {
	Fields Message
	Limit  uint64
	Offset uint64
}

type FilterResponse struct {
	Messages []Message
	Total    int64
}
