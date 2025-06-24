package model

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    int64     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    int64     `gorm:"autoUpdateTime" json:"updated_at"`
}

type DeleteUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}
