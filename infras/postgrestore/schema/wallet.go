package schema

import (
	"go-clean-template/entity"
	"time"
)

type WalletSchema struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id;not null"`
	Balance   float64   `gorm:"column:balance;not null"`
	Currency  string    `gorm:"column:currency;not null"`
	CreatedAt time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (*WalletSchema) TableName() string {
	return "wallets"
}

func (w *WalletSchema) ToWallet() *entity.Wallet {
	return &entity.Wallet{
		ID:       w.ID,
		UserID:   w.UserID,
		Balance:  w.Balance,
		Currency: w.Currency,
	}
}
