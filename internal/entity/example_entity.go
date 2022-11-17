// Package entity
package entity

//go:generate easytags $GOFILE db,json

type Example struct {
	ID      uint64 `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
	Email   string `db:"email" json:"email"`
	Phone   string `db:"phone" json:"phone"`
}
