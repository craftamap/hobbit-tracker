package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey,autoIncrement" json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Secret   string `json:"-"`
	Image    string `json:"image,omitempty"`
}

type AppPassword struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Description string
	Secret      string `json:"-"`
	UserID      uint   `gorm:"index" json:"-"`
	User        *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LastUsedAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (ap *AppPassword) BeforeCreate(scope *gorm.DB) error {
	ap.ID = uuid.New()
	return nil
}
