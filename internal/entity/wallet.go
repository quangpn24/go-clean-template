package entity

import "fmt"

type Wallet struct {
	ID         string
	UserID     string
	WalletName string
}

func NewWallet(id string, userID string, walletName string) (*Wallet, error) {
	if id == "" {
		return nil, fmt.Errorf("id must not be empty")
	}
	return &Wallet{
		ID:         id,
		UserID:     userID,
		WalletName: walletName,
	}, nil
}
