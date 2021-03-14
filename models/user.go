package models

type User struct {
	ID       uint   `gorm:"primaryKey,autoIncrement" json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Secret   string `json:"-"`
	Image    string `json:"image,omitempty"`
}
