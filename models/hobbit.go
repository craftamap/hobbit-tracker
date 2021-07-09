package models

import "time"

type Hobbit struct {
	ID          uint            `gorm:"primaryKey,autoIncrement" json:"id"`
	UserID      uint            `json:"-"`
	User        User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name        string          `json:"name,omitempty"`
	Image       string          `json:"image,omitempty"`
	Description string          `json:"description,omitempty"`
	Records     []NumericRecord `gorm:"foreignKey:HobbitID" json:"records,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
