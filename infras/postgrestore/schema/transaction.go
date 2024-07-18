package schema

import (
	"time"

	"go-clean-template/entity"
)

type TransactionSchema struct {
	ID              string    `gorm:"column:id;primaryKey"`
	WalletID        string    `gorm:"column:wallet_id;not null"`
	AccountID       string    `gorm:"column:account_id;not null"`
	Amount          float64   `gorm:"column:amount;not null"`
	Currency        string    `gorm:"column:currency;not null"`
	TransactionKind string    `gorm:"column:transaction_kind;not null"`
	Status          string    `gorm:"column:status;not null"`
	Note            string    `gorm:"column:note"`
	CreatedAt       time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (*TransactionSchema) TableName() string {
	return "transactions"
}

func ToTransactionSchema(trans *entity.Transaction) *TransactionSchema {
	return &TransactionSchema{
		ID:              trans.ID,
		WalletID:        trans.WalletID,
		AccountID:       trans.AccountID,
		Amount:          trans.Amount,
		Currency:        trans.Currency,
		TransactionKind: string(trans.TransactionKind),
		Status:          string(trans.Status),
		Note:            trans.Note,
	}
}

func (trans *TransactionSchema) ToTransaction() *entity.Transaction {
	return &entity.Transaction{
		ID:              trans.ID,
		WalletID:        trans.WalletID,
		AccountID:       trans.AccountID,
		Amount:          trans.Amount,
		Currency:        trans.Currency,
		TransactionKind: entity.TransactionKind(trans.TransactionKind),
		Status:          entity.TransactionStatus(trans.Status),
		Note:            trans.Note,
	}
}
