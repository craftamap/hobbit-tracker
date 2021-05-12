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
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	Secret    string
	UserID    uint `json:"-"`
	User      User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ap *AppPassword) BeforeCreate(scope *gorm.DB) error {
	ap.ID = uuid.New()
	return nil
}
