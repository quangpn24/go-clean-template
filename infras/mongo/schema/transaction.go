package schema

import (
	"time"

	"go-clean-template/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionSchema struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	WalletID        string             `bson:"wallet_id,omitempty"`
	AccountID       string             `bson:"account_id,omitempty"`
	Amount          float64            `bson:"amount,omitempty"`
	Currency        string             `bson:"currency,omitempty"`
	TransactionKind string             `bson:"transaction_kind,omitempty"`
	Status          string             `bson:"status,omitempty"`
	Note            string             `bson:"note,omitempty"`
	CreatedAt       time.Time          `bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `bson:"updated_at,omitempty"`
}

func ToTransactionSchema(trans *entity.Transaction) *TransactionSchema {
	objID, _ := primitive.ObjectIDFromHex(trans.ID)
	return &TransactionSchema{
		ID:              objID,
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
		ID:              trans.ID.String(),
		WalletID:        trans.WalletID,
		AccountID:       trans.AccountID,
		Amount:          trans.Amount,
		Currency:        trans.Currency,
		TransactionKind: entity.TransactionKind(trans.TransactionKind),
		Status:          entity.TransactionStatus(trans.Status),
		Note:            trans.Note,
	}
}
