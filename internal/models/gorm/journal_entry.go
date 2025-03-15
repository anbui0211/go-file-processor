package gormmodels

import "time"

type JournalEntry struct {
	ID               uint      `gorm:"primaryKey"`
	JournalVoucherID uint      `gorm:"not null"`
	AccountID        uint      `gorm:"not null"`
	DebitAmount      float64   `gorm:"type:decimal(18,2);default:0"`
	CreditAmount     float64   `gorm:"type:decimal(18,2);default:0"`
	Description      string    `gorm:"type:text"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}
