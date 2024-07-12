package interfaces

import (
	"context"

	"go-clean-template/entity"
)

type ITransactionUseCase interface {
	Deposit(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error
	//Withdraw(ctx context.Context) apperror
}

type IPaymentServiceProvider interface {
	Deposit(accountNumber string, bankName string, amount float64, currency string, note string) error
}

type ITransactionRepository interface {
	// WithDBTransaction used to set up a transaction for operations.
	WithDBTransaction(tx IDBTransaction) ITransactionRepository

	// GetWalletByID get a wallet by id. If wallet not found, return nil - nil
	GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error)

	//UpdateWalletBalance update balance by wallet id
	UpdateWalletBalance(ctx context.Context, walletID string, balance float64) error

	//SaveTransaction insert a transaction
	SaveTransaction(ctx context.Context, trans *entity.Transaction) error

	// GetAccountByID get accoutn by id. If account not found, return nil - nil
	GetAccountByID(ctx context.Context, accountID string) (*entity.Account, error)
}

type INotifier interface {
	SendNotification(ctx context.Context, message string)
}

type IDBTransaction interface {
	// Begin begins a transaction
	Begin(ctx context.Context) (IDBTransaction, error)

	// Commit commits the changes in a transaction
	Commit(ctx context.Context) error

	// Rollback rollbacks the changes in a transaction
	Rollback(ctx context.Context)
}
