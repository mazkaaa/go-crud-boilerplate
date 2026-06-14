package models

import (
	"time"
)

type Role struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Users     []User    `json:"users"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
