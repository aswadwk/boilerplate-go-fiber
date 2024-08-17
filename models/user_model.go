package models

import (
	"time"
)

type User struct {
	ID         string     `json:"id" gorm:"primary_key"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Password   string     `json:"-"`
	Role       string     `json:"role"`
	AvatarPath string     `json:"avatar_path"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"update_at"`
	// DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
