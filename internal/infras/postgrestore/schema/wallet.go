package schema

import (
	"time"

	"go-clean-template/internal/entity"
)

type WalletSchema struct {
	ID         string    `gorm:"column:id;primaryKey"`
	UserID     string    `gorm:"column:user_id;not null"`
	WalletName string    `gorm:"column:wallet_name;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (*WalletSchema) TableName() string {
	return "wallets"
}

func (w *WalletSchema) ToWallet() *entity.Wallet {
	return &entity.Wallet{
		ID:         w.ID,
		UserID:     w.UserID,
		WalletName: w.WalletName,
	}
}
