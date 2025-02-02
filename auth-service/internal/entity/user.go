package entity

import "time"

type User struct {
	ID           int64
	Email        string
	HashPassword string
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
