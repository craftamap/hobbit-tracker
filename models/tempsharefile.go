package models

import "time"

type TemporaryShareFile struct {
	CreatedAt time.Time
	MimeType  string
	Bytes     []byte
	UserID      uint            `json:"-"`
	User        User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ID        uint `gorm:"primaryKey,autoIncrement" json:"id"`
}
