package entity

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id_user" db:"id_user"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Password    string    `json:"password" db:"password"`
	ImageUrl    string    `json:"image_url" db:"image_url"`
}

type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
