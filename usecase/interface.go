package usecase

import (
	"context"
	"go-clean-template/entity"
)

type ITransactionUseCase interface {
	Deposit(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error
	//Withdraw(ctx context.Context) apperror
}

type IBankService interface {
	Deposit(accountNumber string, bankName string, amount float64, currency string, note string) error
}

type ITransactionRepository interface {
	WithDBTransaction(tx IDBTransaction) ITransactionRepository
	GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID string, balance float64) error
	SaveTransaction(ctx context.Context, trans *entity.Transaction) error
	GetAccountByID(ctx context.Context, accountID string) (*entity.Account, error)
}

type INotifier interface {
	SendNotification(ctx context.Context, message string)
}

type IDBTransaction interface {
	Begin(ctx context.Context) (IDBTransaction, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}
