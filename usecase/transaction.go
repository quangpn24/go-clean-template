package usecase

import (
	"context"
	"fmt"

	"go-clean-template/entity"
	"go-clean-template/pkg/apperror"

	"github.com/google/uuid"
)

type TransactionUseCase struct {
	repo       ITransactionRepository
	paymentSvc IPaymentServiceProvider
	notifiers  []INotifier
}

func NewTransactionUseCase(repo ITransactionRepository, paymentSvc IPaymentServiceProvider) *TransactionUseCase {
	return &TransactionUseCase{
		repo:       repo,
		paymentSvc: paymentSvc,
		notifiers:  []INotifier{},
	}
}

func (uc *TransactionUseCase) SetNotifiers(notifiers ...INotifier) {
	uc.notifiers = append(uc.notifiers, notifiers...)
}

func (uc *TransactionUseCase) Deposit(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error {
	var (
		transID = uuid.New().String()
		err     error
		trans   *entity.Transaction
	)

	// check account
	account, err := uc.repo.GetLinkedAccountByID(ctx, accountID)
	if err != nil {
		return apperror.ErrGet(err, "failed to get account by id")
	}
	if account == nil {
		return apperror.ErrInvalidParams(fmt.Errorf("account not found"))
	}

	// create new transaction
	trans = entity.NewTransaction(transID, walletID, accountID, amount, currency, entity.TransactionIn, note, entity.TransactionStatusNew)

	// get wallet
	wallet, err := uc.repo.GetWalletByID(ctx, walletID)
	if err != nil {
		return apperror.ErrGet(err, "failed to get wallet by id")
	}

	if wallet == nil {
		return apperror.ErrInvalidParams(fmt.Errorf("wallet not found"))
	}

	// save transaction
	if err := uc.repo.SaveTransaction(ctx, trans); err != nil {
		return apperror.ErrCreate(err, "failed to create deposit transaction")
	}

	return nil
}

func (uc *TransactionUseCase) Withdraw(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error {
	var (
		transID = uuid.New().String()
		err     error
		trans   *entity.Transaction
	)

	// check account linking status
	account, err := uc.repo.GetLinkedAccountByID(ctx, accountID)
	if err != nil {
		return apperror.ErrGet(err, "failed to get account by id")
	}
	if account == nil {
		return apperror.ErrInvalidParams(fmt.Errorf("account not found"))
	}

	// get wallet
	wallet, err := uc.repo.GetWalletByID(ctx, walletID)
	if err != nil {
		return apperror.ErrGet(err, "failed to get wallet by id")
	}

	if wallet == nil {
		return apperror.ErrInvalidParams(fmt.Errorf("wallet not found"))
	}

	//check balance
	balance, err := uc.repo.GetBalanceByWalletID(ctx, walletID)
	if err != nil {
		return apperror.ErrGet(err, "failed to get balance by wallet id")
	}

	if balance < amount {
		return apperror.ErrInvalidParams(fmt.Errorf("insufficient balance"))
	}
	// create new transaction
	trans = entity.NewTransaction(transID, walletID, accountID, amount, currency, entity.TransactionOut, note, entity.TransactionStatusNew)

	// save transaction
	if err := uc.repo.SaveTransaction(ctx, trans); err != nil {
		return apperror.ErrCreate(err, "failed to create withdraw transaction")
	}

	return nil
}

func (uc *TransactionUseCase) PayTransaction(ctx context.Context, transID string) error {
	var (
		transStatus entity.TransactionStatus
		err         error
	)
	// get trans
	trans, err := uc.repo.GetTransactionByID(ctx, transID)

	if err != nil {
		return apperror.ErrGet(err, "failed to get transaction by id")
	}

	if trans == nil {
		return apperror.ErrInvalidParams(fmt.Errorf("no transactions found in ready-to-pay status"))
	}

	// check transaction status
	if trans.Status != entity.TransactionStatusNew {
		return apperror.ErrInvalidParams(fmt.Errorf("transaction status is not new"))
	}

	// send to Momo
	if trans.TransactionKind == entity.TransactionIn {
		err = uc.paymentSvc.Deposit(ctx, trans.Amount, trans.Currency, trans.Note)
	} else if trans.TransactionKind == entity.TransactionOut {
		err = uc.paymentSvc.Withdraw(ctx, trans.Amount, trans.Currency, trans.Note)
	}

	if err != nil {
		transStatus = entity.TransactionStatusFailed
	} else {
		transStatus = entity.TransactionStatusSuccessful
	}

	// Update transaction status
	if err := uc.repo.UpdateTransactionStatus(ctx, transID, transStatus); err != nil {
		return apperror.ErrUpdate(err, "failed to update transaction status")
	}
	return nil
}
