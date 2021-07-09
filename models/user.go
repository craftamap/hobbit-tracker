package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint    `gorm:"primaryKey,autoIncrement" json:"id,omitempty"`
	Username string  `json:"username,omitempty"`
	Secret   string  `json:"-"`
	Image    string  `json:"image,omitempty"`
	Follows  []*User `json:"-" gorm:"many2many:user_follows;"`
}

type AppPassword struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id,omitempty"`
	Description string         `json:"description,omitempty"`
	Secret      string         `json:"secret,omitempty"`
	UserID      uint           `gorm:"index" json:"user_id,omitempty"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LastUsedAt  time.Time      `json:"last_used_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (ap *AppPassword) BeforeCreate(scope *gorm.DB) error {
	ap.ID = uuid.New()
	return nil
}
