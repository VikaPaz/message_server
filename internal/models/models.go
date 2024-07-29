package models

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusCompleted  Status = "COMPLETED"
	StatusProcessing Status = "PROCESSING"
	StatusFailed     Status = "FAILED"
)

type Message struct {
	ID        *uuid.UUID
	Message   *string
	Status    *Status
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type CreateRequest struct {
	Message string `json:"message"`
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

type WriteRequest struct {
	ID      uuid.UUID
	Message *string
}

type UpdateRequest struct {
	ID     uuid.UUID
	Status *Status
}
