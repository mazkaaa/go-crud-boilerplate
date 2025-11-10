package models

import (
	"time"

	"github.com/go-playground/validator"
)

type (
	User struct {
		ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
		Name      string    `json:"name" gorm:"not null; index:idx_name"`
		Email     string    `json:"email" gorm:"unique; not null; index:idx_email"`
		Password  string    `json:"-"`
		CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
		UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
		RoleID    string    `json:"role_id"`
		Role      Role      `json:"-" gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL"`
	}

	UserValidator struct {
		validator *validator.Validate
	}
)
