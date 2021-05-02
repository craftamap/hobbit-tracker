package models

import "time"

type NumericRecord struct {
	ID        uint      `gorm:"primaryKey,autoIncrement" json:"id"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
	HobbitID  uint      `json:"hobbit_id,omitempty"`
	Hobbit    Hobbit    `gorm:"foreignKey:HobbitID" json:"hobbit,omitempty"`
	Value     int64     `json:"value,omitempty"`
	Comment   string    `json:"comment,omitempty"`
}
