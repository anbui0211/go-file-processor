package repository

import (
	"context"
	"gofile/global"
	gormmodels "gofile/internal/models/gorm"
)

type AccountRepository struct {
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (r *AccountRepository) FindAll(ctx context.Context) ([]gormmodels.Account, error) {
	var accounts []gormmodels.Account

	result := global.Mdb.WithContext(ctx).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}
