package entity

import (
	"github.com/google/uuid"
)

type Metadata struct {
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"`
	PerPage       int        `json:"per_page"`
	Page          int        `json:"page"`
	Total         int        `json:"total"`
	OrderBy       string     `json:"order_by,omitempty"`
	OrderType     string     `json:"order_type,omitempty"`
}

type Meta struct {
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"`
}
