package schema

import (
	"go-clean-template/entity"
	"time"
)

type TransactionSchema struct {
	ID               string
	SenderWalletID   *string
	ReceiverWalletID *string
	AccountID        *string
	Amount           float64
	Currency         string
	Category         string
	TransactionKind  string
	Note             string
	CreatedAt        time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (*TransactionSchema) TableName() string {
	return "transactions"
}

func ToTransactionSchema(trans *entity.Transaction) *TransactionSchema {
	var (
		senderWalletID   *string
		receiverWalletID *string
		accountId        *string
	)

	if trans.SenderWalletID != "" {
		senderWalletID = &trans.SenderWalletID
	}

	if trans.ReceiverWalletID != "" {
		receiverWalletID = &trans.ReceiverWalletID
	}

	if trans.AccountID != "" {
		accountId = &trans.AccountID
	}

	return &TransactionSchema{
		ID:               trans.ID,
		SenderWalletID:   senderWalletID,
		ReceiverWalletID: receiverWalletID,
		AccountID:        accountId,
		Amount:           trans.Amount,
		Currency:         trans.Currency,
		Category:         string(trans.Category),
		Note:             trans.Note,
	}
}

func (trans *TransactionSchema) ToTransaction() *entity.Transaction {
	var (
		senderWalletID   string
		receiverWalletID string
		accountId        string
	)

	if trans.SenderWalletID != nil {
		senderWalletID = *trans.SenderWalletID
	}

	if trans.ReceiverWalletID != nil {
		receiverWalletID = *trans.ReceiverWalletID
	}

	if trans.AccountID != nil {
		accountId = *trans.AccountID
	}

	return &entity.Transaction{
		ID:               trans.ID,
		SenderWalletID:   senderWalletID,
		ReceiverWalletID: receiverWalletID,
		AccountID:        accountId,
		Amount:           trans.Amount,
		Currency:         trans.Currency,
		Category:         entity.TransactionCategory(trans.Category),
		Note:             trans.Note,
	}
}
