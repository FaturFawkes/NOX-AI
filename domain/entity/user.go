package entity

import (
	"time"

	"gorm.io/gorm"
)

type Plan string

const (
	Free    Plan = "free"
	Basic   Plan = "basic"
	Premium Plan = "premium"
)

type User struct {
	gorm.Model
	Name      string
	Number    string
	ExpiredAt time.Time
	Plan      Plan `gorm:"type:enum('free', 'basic', 'premium')"`
	Limit     int
}