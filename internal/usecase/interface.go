package usecase

import (
	"context"

	"go-clean-template/internal/entity"
)

type ITransactionUseCase interface {
	Deposit(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error
	Withdraw(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error
	PayTransaction(ctx context.Context, transID string) error
}

type IPaymentServiceProvider interface {
	Deposit(ctx context.Context, amount float64, currency string, note string) error
	Withdraw(ctx context.Context, amount float64, currency string, note string) error
}

type ITransactionRepository interface {
	// GetWalletByID get a wallet by id. If wallet not found, return nil - nil
	GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error)

	//SaveTransaction insert a transaction
	SaveTransaction(ctx context.Context, trans *entity.Transaction) error

	// GetLinkedAccountByID get account by id. If account not found, return nil - nil
	GetLinkedAccountByID(ctx context.Context, accountID string) (*entity.LinkedAccount, error)

	// GetBalanceByWalletID get balance by wallet id
	GetBalanceByWalletID(ctx context.Context, walletID string) (float64, error)

	// GetTransactionByID get transaction by id. If Transaction not found, return nil - nil
	GetTransactionByID(ctx context.Context, transID string) (*entity.Transaction, error)

	// UpdateTransactionStatus update transaction status
	UpdateTransactionStatus(ctx context.Context, transID string, status entity.TransactionStatus) error
}

type INotifier interface {
	SendNotification(ctx context.Context, message string)
}
