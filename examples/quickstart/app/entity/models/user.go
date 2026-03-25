package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string         `gorm:"size:32;not null;uniqueIndex" json:"username"`
	Nickname     string         `gorm:"size:64;not null;default:''" json:"nickname"`
	Email        string         `gorm:"size:128;not null;uniqueIndex" json:"email"`
	Phone        string         `gorm:"size:20;not null;default:''" json:"phone"`
	PasswordHash string         `gorm:"size:60;not null" json:"-"`
	Status       int8           `gorm:"not null;default:1;index" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
