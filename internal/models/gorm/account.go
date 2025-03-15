package gormmodels

import "time"

// Account model
type Account struct {
	ID        uint      `gorm:"primaryKey"`
	Code      string    `gorm:"size:20;unique;not null"`
	Name      string    `gorm:"size:255;not null"`
	Type      string    `gorm:"type:enum('asset','liability','equity','revenue','expense');not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
