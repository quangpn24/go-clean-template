package entity

import "fmt"

type Wallet struct {
	ID       string
	UserID   string
	Balance  float64
	Currency string
}

func NewWallet(id string, userID string, balance float64, currency string) (*Wallet, error) {
	if id == "" {
		return nil, fmt.Errorf("id must not be empty")
	}
	return &Wallet{
		ID:       id,
		UserID:   userID,
		Balance:  balance,
		Currency: currency,
	}, nil
}

func (w *Wallet) Deposit(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	w.Balance += amount
	return nil
}

func (w *Wallet) Withdraw(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if w.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	w.Balance -= amount
	return nil
}
