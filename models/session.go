package models

import (
	"time"
)

type Session struct {
	ID             uint      `gorm:"primaryKey,autoIncrement" json:"-"`
	Token          string    `gorm:"unique;size:36;not null" json:"-"`
	ExpirationTime time.Time `gorm:"not null" json:"-"`
	UserId         uint      `json:"-"`
}
