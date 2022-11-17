package entity

import (
	"time"

	"github.com/google/uuid"
)

type Food struct {
	ID          uuid.UUID `json:"id_food" db:"id_food"`
	IDUser      uuid.UUID `json:"id_user" db:"id_user"`
	Name        string    `json:"name,omitempty" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Category    string    `json:"category,omitempty" db:"category"`
	Quantity    int       `json:"quantity,omitempty" db:"quantity"`
	ImageUrl    string    `json:"image_url,omitempty" db:"image_url"`
	IsActive    string    `json:"is_active,omitempty" db:"is_active"`
	ExpiredAt   time.Time `json:"expired_at,omitempty" db:"expired_at"`
	Latitude    string    `json:"latitude,omitempty" db:"latitude"`
	Longitude   string    `json:"longitude,omitempty" db:"longitude"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
	// Location    string    `json:"location" db:"location"`
	// Status      int64     `json:"status" db:"status"`
}
