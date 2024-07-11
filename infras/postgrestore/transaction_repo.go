package postgrestore

import (
	"context"

	"go-clean-template/entity"
	"go-clean-template/infras/postgrestore/schema"
	"go-clean-template/usecase/interfaces"

	"gorm.io/gorm"
)

const (
	WalletTable       = "wallets"
	TransactionsTable = "transactions"
	AccountTable      = "accounts"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) WithDBTransaction(tx interfaces.IDBTransaction) interfaces.ITransactionRepository {
	trans, ok := tx.(*DBTransaction)
	if !ok {
		return r
	}
	return &TransactionRepo{
		db: trans.db,
	}
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

func (r *TransactionRepo) UpdateWalletBalance(ctx context.Context, walletID string, balance float64) error {
	return r.db.WithContext(ctx).Table(WalletTable).Where("id = ?", walletID).
		Update("balance", balance).Error
}

func (r *TransactionRepo) SaveTransaction(ctx context.Context, trans *entity.Transaction) error {
	transSchema := schema.ToTransactionSchema(trans)
	return r.db.WithContext(ctx).Table(TransactionsTable).Create(transSchema).Error
}

func (r *TransactionRepo) GetAccountByID(ctx context.Context, accountID string) (*entity.Account, error) {
	var (
		account       *entity.Account
		accountSchema schema.AccountSchema
	)
	if err := r.db.WithContext(ctx).Table(AccountTable).Where("id = ?", accountID).Take(&accountSchema).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	account = accountSchema.ToAccount()
	return account, nil
}
