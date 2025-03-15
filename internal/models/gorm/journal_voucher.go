package gormmodels

import "time"

// JournalVoucher model
type JournalVoucher struct {
    ID          uint      `gorm:"primaryKey"`
    VoucherNo   string    `gorm:"size:50;unique;not null"`
    Date        time.Time `gorm:"not null"`
    Description string    `gorm:"type:text"`
    Status      string    `gorm:"type:enum('pending','approved','rejected');default:'pending'"`
    CreatedBy   uint      `gorm:"not null"`
    ApprovedBy  *uint     `gorm:"default:null"`
    TotalDebit  float64   `gorm:"type:decimal(18,2);not null;default:0"`
    TotalCredit float64   `gorm:"type:decimal(18,2);not null;default:0"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
}
