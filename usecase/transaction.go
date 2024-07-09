package usecase

import (
	"context"
	"fmt"
	"go-clean-template/entity"
	"go-clean-template/pkg/apperror"
	"go-clean-template/usecase/interfaces"

	"github.com/google/uuid"
)

type TransactionUseCase struct {
	repo          interfaces.ITransactionRepository
	bankSvc       interfaces.IBankService
	notifiers     []interfaces.INotifier
	dbTransaction interfaces.IDBTransaction
}

func NewTransactionUseCase(repo interfaces.ITransactionRepository, bankSvc interfaces.IBankService, dbTransaction interfaces.IDBTransaction) *TransactionUseCase {
	return &TransactionUseCase{
		repo:          repo,
		bankSvc:       bankSvc,
		notifiers:     []interfaces.INotifier{},
		dbTransaction: dbTransaction,
	}
}

func (uc *TransactionUseCase) SetNotifiers(notifiers ...interfaces.INotifier) {
	uc.notifiers = append(uc.notifiers, notifiers...)
}

func (uc *TransactionUseCase) Deposit(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error {
	var (
		transID = uuid.New().String()
		err     error
		trans   *entity.Transaction
	)

	// check account linking status
	account, err := uc.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return apperror.ErrGet(err, "failed to get account by id")
	}
	if account == nil {
		return apperror.ErrInvalidParams(fmt.Errorf("account not found"))
	}

	if !account.IsLinked {
		return apperror.ErrInvalidParams(fmt.Errorf("account not linked"))
	}

	// create new transaction
	trans = entity.NewTransaction(transID, "", walletID, accountID, amount, currency, entity.CategoryDeposit, note)

	// get wallet
	wallet, err := uc.repo.GetWalletByID(ctx, walletID)
	if err != nil {
		return apperror.ErrGet(err, "failed to get wallet by id")
	}

	if wallet == nil {
		return apperror.ErrInvalidParams(fmt.Errorf("wallet not found"))
	}

	// deposit to wallet
	if err = wallet.Deposit(amount); err != nil {
		return apperror.ErrInvalidParams(err)
	}

	// start transaction
	tx, err := uc.dbTransaction.Begin(ctx)
	if err != nil {
		return apperror.ErrOtherInternalServerError(err, "error when starting transaction")
	}

	repo := uc.repo.WithDBTransaction(tx)

	// update wallet
	if err := repo.UpdateWalletBalance(ctx, walletID, wallet.Balance); err != nil {
		tx.Rollback(ctx)
		return apperror.ErrUpdate(err, "failed to update balance")
	}

	// save transaction
	if err := repo.SaveTransaction(ctx, trans); err != nil {
		tx.Rollback(ctx)
		return apperror.ErrCreate(err, "failed to create deposit transaction")
	}

	// call bank service
	if err := uc.bankSvc.Deposit(account.AccountNumber, account.BankName, amount, currency, note); err != nil {
		tx.Rollback(ctx)
		return apperror.ErrThirdParty(err, "error when calling api deposit bank service")
	}

	// Commit and finish transaction
	if err := tx.Commit(ctx); err != nil {
		return apperror.ErrOtherInternalServerError(err, "error when committing transaction")
	}

	// notification, don't care result
	for _, notifier := range uc.notifiers {
		notifier.SendNotification(ctx, "Deposit success")
	}

	return nil
}
