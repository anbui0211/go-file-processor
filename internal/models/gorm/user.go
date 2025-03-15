package gormmodels

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:255;not null"`
    Email     string    `gorm:"size:255;unique;not null"`
    Role      string    `gorm:"type:enum('admin','accountant','viewer');not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}




