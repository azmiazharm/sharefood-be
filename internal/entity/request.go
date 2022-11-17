package entity

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	ID        uuid.UUID `json:"id_request" db:"id_request"`
	IDUser    uuid.UUID `json:"id_user" db:"id_user"`
	IDFood    uuid.UUID `json:"id_food" db:"id_food"`
	Status    int       `json:"status" db:"status"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type RequestAction struct {
	ID     uuid.UUID `json:"id_request" db:"id_request"`
	Action string    `json:"action"`
}

type RequestWithFood struct {
	ID         uuid.UUID `json:"id_request" db:"requests.id_request"`
	IDUser     uuid.UUID `json:"id_user" db:"requests.id_user"`
	IDFood     uuid.UUID `json:"id_food" db:"requests.id_food"`
	Status     int       `json:"status" db:"requests.status"`
	Quantity   int       `json:"quantity" db:"requests.quantity"`
	IDUserFood uuid.UUID `json:"id_user_food" db:"foods.id_user"`
	Stock      int       `json:"stock" db:"foods.quantity"`
}

// type RequestInput struct {
// 	IDFood    uuid.UUID `json:"id_food" db:"id_food"`
// 	Quantity  int       `json:"quantity" db:"quantity"`
// 	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
// }
