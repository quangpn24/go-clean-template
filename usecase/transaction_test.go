package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go-clean-template/entity"
	"go-clean-template/pkg/apperror"
	"go-clean-template/usecase/interfaces"
	"go-clean-template/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTransactionUseCase(t *testing.T) {
	type args struct {
		repo          interfaces.ITransactionRepository
		paymentSvc    interfaces.IPaymentServiceProvider
		dbTransaction interfaces.IDBTransaction
	}
	tests := []struct {
		name string
		args args
		want *TransactionUseCase
	}{
		{
			name: "create new transaction use case",
			args: args{
				repo:          mocks.NewITransactionRepository(t),
				paymentSvc:    mocks.NewIPaymentServiceProvider(t),
				dbTransaction: mocks.NewIDBTransaction(t),
			},
			want: &TransactionUseCase{
				repo:          mocks.NewITransactionRepository(t),
				paymentSvc:    mocks.NewIPaymentServiceProvider(t),
				notifiers:     []interfaces.INotifier{},
				dbTransaction: mocks.NewIDBTransaction(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionUseCase(tt.args.repo, tt.args.paymentSvc, tt.args.dbTransaction); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionUseCase_SetNotifiers(t *testing.T) {
	tests := []struct {
		name      string
		notifiers []interfaces.INotifier
		object    *TransactionUseCase
		want      *TransactionUseCase
	}{
		{
			name:      "set notifiers",
			notifiers: []interfaces.INotifier{mocks.NewINotifier(t)},
			object: &TransactionUseCase{
				notifiers: []interfaces.INotifier{},
			},
			want: &TransactionUseCase{
				notifiers: []interfaces.INotifier{mocks.NewINotifier(t)},
			},
		},
	}
	for _, tt := range tests {
		tt.object.SetNotifiers(tt.notifiers...)

		assert.Equal(t, tt.want, tt.object)
	}
}

func TestTransactionUseCase_Deposit(t *testing.T) {
	transRepo := mocks.NewITransactionRepository(t)
	paymentSvc := mocks.NewIPaymentServiceProvider(t)
	dbTrans := mocks.NewIDBTransaction(t)
	notifier := mocks.NewINotifier(t)
	uc := TransactionUseCase{
		repo:          transRepo,
		paymentSvc:    paymentSvc,
		notifiers:     []interfaces.INotifier{notifier},
		dbTransaction: dbTrans,
	}
	t.Run("success", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		userID := "u_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        userID,
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		walletMock := &entity.Wallet{
			ID:       walletID,
			UserID:   userID,
			Balance:  10000,
			Currency: currency,
		}

		newTransMock := &entity.Transaction{
			ID:               "t_00001",
			ReceiverWalletID: walletID,
			AccountID:        accountID,
			Amount:           amount,
			Currency:         currency,
			Category:         entity.CategoryDeposit,
			Note:             note,
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().UpdateWalletBalance(ctx, walletID, walletMock.Balance+amount).Return(nil).Once()
		transRepo.EXPECT().SaveTransaction(ctx, IsMatchByTransaction(newTransMock)).Return(nil).Once()
		transRepo.EXPECT().WithDBTransaction(dbTrans).Return(transRepo).Once()

		dbTrans.EXPECT().Begin(ctx).Return(dbTrans, nil).Once()

		dbTrans.EXPECT().Commit(ctx).Return(nil).Once()

		paymentSvc.EXPECT().Deposit(accountMock.AccountNumber, accountMock.BankName, amount, currency, note).Return(nil).Once()

		notifier.EXPECT().SendNotification(ctx, "Deposit success").Return().Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.NoError(t, err)

	})

	t.Run("failed to get account by id", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"
		errDB := fmt.Errorf("unexpected error")

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(nil, errDB).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrGet(errDB, "failed to get account by id")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("account not found", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(nil, nil).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("account not found"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("account not linked", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      false,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("account not linked"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to get wallet by id", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"
		errDB := fmt.Errorf("unexpected error")

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(nil, errDB).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrGet(errDB, "failed to get wallet by id")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("wallet not found", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(nil, nil).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("wallet not found"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to deposit to wallet", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := -1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		walletMock := &entity.Wallet{
			ID:       walletID,
			UserID:   "u_00001",
			Balance:  10000,
			Currency: currency,
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("amount must be greater than 0"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to start transaction", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"
		errBeginTrans := fmt.Errorf("unexpected error")

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		walletMock := &entity.Wallet{
			ID:       walletID,
			UserID:   "u_00001",
			Balance:  10000,
			Currency: currency,
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()

		dbTrans.EXPECT().Begin(ctx).Return(nil, errBeginTrans).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrOtherInternalServerError(errBeginTrans, "error when starting transaction")
		assert.Equal(t, expectedErr, err)

	})

	t.Run("failed to update wallet balance", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"
		errUpdateWallet := fmt.Errorf("unexpected error")

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		walletMock := &entity.Wallet{
			ID:       walletID,
			UserID:   "u_00001",
			Balance:  10000,
			Currency: currency,
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().WithDBTransaction(dbTrans).Return(transRepo).Once()
		transRepo.EXPECT().UpdateWalletBalance(ctx, walletID, walletMock.Balance+amount).Return(errUpdateWallet).Once()

		dbTrans.EXPECT().Begin(ctx).Return(dbTrans, nil).Once()
		dbTrans.EXPECT().Rollback(ctx).Return().Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrUpdate(errUpdateWallet, "failed to update balance")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to create deposit transaction", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"
		errSaveTrans := fmt.Errorf("unexpected error")

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		walletMock := &entity.Wallet{
			ID:       walletID,
			UserID:   "u_00001",
			Balance:  10000,
			Currency: currency,
		}

		newTrans := &entity.Transaction{
			ID:               "t_00001",
			ReceiverWalletID: walletID,
			AccountID:        accountID,
			Amount:           amount,
			Currency:         currency,
			Category:         entity.CategoryDeposit,
			Note:             note,
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().WithDBTransaction(dbTrans).Return(transRepo).Once()
		transRepo.EXPECT().UpdateWalletBalance(ctx, walletID, walletMock.Balance+amount).Return(nil).Once()
		transRepo.EXPECT().SaveTransaction(ctx, IsMatchByTransaction(newTrans)).Return(errSaveTrans).Once()

		dbTrans.EXPECT().Begin(ctx).Return(dbTrans, nil).Once()
		dbTrans.EXPECT().Rollback(ctx).Return().Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrCreate(errSaveTrans, "failed to create deposit transaction")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to call payment service", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"
		errPaymentSvc := fmt.Errorf("unexpected error")

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		walletMock := &entity.Wallet{
			ID:       walletID,
			UserID:   "u_00001",
			Balance:  10000,
			Currency: currency,
		}

		newTrans := &entity.Transaction{
			ID:               "t_00001",
			ReceiverWalletID: walletID,
			AccountID:        accountID,
			Amount:           amount,
			Currency:         currency,
			Category:         entity.CategoryDeposit,
			Note:             note,
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().WithDBTransaction(dbTrans).Return(transRepo).Once()
		transRepo.EXPECT().UpdateWalletBalance(ctx, walletID, walletMock.Balance+amount).Return(nil).Once()
		transRepo.EXPECT().SaveTransaction(ctx, IsMatchByTransaction(newTrans)).Return(nil).Once()

		dbTrans.EXPECT().Begin(ctx).Return(dbTrans, nil).Once()
		dbTrans.EXPECT().Rollback(ctx).Return().Once()

		paymentSvc.EXPECT().Deposit(accountMock.AccountNumber, accountMock.BankName, amount, currency, note).Return(errPaymentSvc).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrThirdParty(errPaymentSvc, "error when calling api deposit payment service")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to commit transaction", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Deposit 1,000,000 VND"
		errCommitTrans := fmt.Errorf("unexpected error")

		accountMock := &entity.Account{
			ID:            accountID,
			UserID:        "u_00001",
			IsLinked:      true,
			AccountNumber: "1234567890",
			BankName:      "Bank A",
		}

		walletMock := &entity.Wallet{
			ID:       walletID,
			UserID:   "u_00001",
			Balance:  10000,
			Currency: currency,
		}

		newTrans := &entity.Transaction{
			ID:               "t_00001",
			ReceiverWalletID: walletID,
			AccountID:        accountID,
			Amount:           amount,
			Currency:         currency,
			Category:         entity.CategoryDeposit,
			Note:             note,
		}

		transRepo.EXPECT().GetAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().WithDBTransaction(dbTrans).Return(transRepo).Once()
		transRepo.EXPECT().UpdateWalletBalance(ctx, walletID, walletMock.Balance+amount).Return(nil).Once()
		transRepo.EXPECT().SaveTransaction(ctx, IsMatchByTransaction(newTrans)).Return(nil).Once()

		dbTrans.EXPECT().Begin(ctx).Return(dbTrans, nil).Once()
		dbTrans.EXPECT().Commit(ctx).Return(errCommitTrans).Once()

		paymentSvc.EXPECT().Deposit(accountMock.AccountNumber, accountMock.BankName, amount, currency, note).Return(nil).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrOtherInternalServerError(errCommitTrans, "error when committing transaction")
		assert.Equal(t, expectedErr, err)
	})
}

func IsMatchByTransaction(a *entity.Transaction) interface{} {
	return mock.MatchedBy(func(b *entity.Transaction) bool {
		return a.SenderWalletID == b.SenderWalletID &&
			a.ReceiverWalletID == b.ReceiverWalletID &&
			a.AccountID == b.AccountID &&
			a.Amount == b.Amount &&
			a.Currency == b.Currency &&
			a.Category == b.Category &&
			a.Note == b.Note
	})
}
