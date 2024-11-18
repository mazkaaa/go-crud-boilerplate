package models

import (
	"time"
)

type Role struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"unique; not null"`
	Users     []User    `json:"users" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
