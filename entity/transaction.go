package entity

import "fmt"

type TransactionKind string

const (
	TransactionOut TransactionKind = "OUT"
	TransactionIn  TransactionKind = "IN"
)

type TransactionStatus string

const (
	TransactionStatusNew        TransactionStatus = "NEW"
	TransactionStatusSuccessful TransactionStatus = "SUCCESSFUL"
	TransactionStatusFailed     TransactionStatus = "FAILED"
)

type Transaction struct {
	ID              string
	WalletID        string
	AccountID       string
	Amount          float64
	Currency        string
	TransactionKind TransactionKind
	Note            string
	Status          TransactionStatus
}

func NewTransaction(id string, walletID string, accountID string, amount float64, currency string, transKind TransactionKind, note string, status TransactionStatus) *Transaction {
	return &Transaction{
		ID:              id,
		WalletID:        walletID,
		AccountID:       accountID,
		Amount:          amount,
		Currency:        currency,
		TransactionKind: transKind,
		Note:            note,
		Status:          status,
	}
}

func (t *Transaction) ToSuccessful() error {
	if t.Status != TransactionStatusNew {
		return fmt.Errorf("cant update transaction status from %s to %s", t.Status, TransactionStatusSuccessful)
	}
	t.Status = TransactionStatusSuccessful
	return nil
}

func (t *Transaction) ToFailed() error {
	if t.Status != TransactionStatusNew {
		return fmt.Errorf("cant update transaction status from %s to %s", t.Status, TransactionStatusFailed)
	}
	t.Status = TransactionStatusFailed
	return nil
}
