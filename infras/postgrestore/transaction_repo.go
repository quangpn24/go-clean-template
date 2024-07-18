package postgrestore

import (
	"context"

	"go-clean-template/entity"
	"go-clean-template/infras/postgrestore/schema"

	"gorm.io/gorm"
)

const (
	WalletTable        = "wallets"
	TransactionsTable  = "transactions"
	LinkedAccountTable = "linked_accounts"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	var (
		wallet       *entity.Wallet
		walletSchema schema.WalletSchema
	)
	if err := r.db.WithContext(ctx).Table(WalletTable).Where("id = ?", walletID).Take(&walletSchema).Error;
		err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	wallet = walletSchema.ToWallet()
	return wallet, nil
}

func (r *TransactionRepo) SaveTransaction(ctx context.Context, trans *entity.Transaction) error {
	transSchema := schema.ToTransactionSchema(trans)
	return r.db.WithContext(ctx).Table(TransactionsTable).Create(transSchema).Error
}

func (r *TransactionRepo) GetLinkedAccountByID(ctx context.Context, accountID string) (*entity.LinkedAccount, error) {
	var (
		account       *entity.LinkedAccount
		accountSchema schema.LinkedAccountSchema
	)
	if err := r.db.WithContext(ctx).Table(LinkedAccountTable).Where("id = ?", accountID).Take(&accountSchema).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	account = accountSchema.ToLinkedAccount()
	return account, nil
}

func (r *TransactionRepo) GetBalanceByWalletID(ctx context.Context, walletID string) (float64, error) {
	var balance float64
	selectQuery := `SUM(CASE WHEN transaction_kind = ? THEN amount ELSE -amount END)`
	if err := r.db.WithContext(ctx).Table(TransactionsTable).
		Select(selectQuery, entity.TransactionIn).
		Where("wallet_id = ? and status = ?", walletID, entity.TransactionStatusSuccessful).
		Row().Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *TransactionRepo) GetTransactionByID(ctx context.Context, transID string) (*entity.Transaction, error) {
	var transSchema schema.TransactionSchema
	if err := r.db.WithContext(ctx).Table(TransactionsTable).Where("id = ?", transID).Take(&transSchema).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return transSchema.ToTransaction(), nil
}

func (r *TransactionRepo) UpdateTransactionStatus(ctx context.Context, transID string, status entity.TransactionStatus) error {
	return r.db.WithContext(ctx).Table(TransactionsTable).Where("id = ?", transID).
		Update("status", string(status)).Error
}
