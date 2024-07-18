package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go-clean-template/entity"
	"go-clean-template/pkg/apperror"
	"go-clean-template/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTransactionUseCase(t *testing.T) {
	type args struct {
		repo       ITransactionRepository
		paymentSvc IPaymentServiceProvider
	}
	tests := []struct {
		name string
		args args
		want *TransactionUseCase
	}{
		{
			name: "create new transaction use case",
			args: args{
				repo:       mocks.NewITransactionRepository(t),
				paymentSvc: mocks.NewIPaymentServiceProvider(t),
			},
			want: &TransactionUseCase{
				repo:       mocks.NewITransactionRepository(t),
				paymentSvc: mocks.NewIPaymentServiceProvider(t),
				notifiers:  []INotifier{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransactionUseCase(tt.args.repo, tt.args.paymentSvc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionUseCase_SetNotifiers(t *testing.T) {
	tests := []struct {
		name      string
		notifiers []INotifier
		object    *TransactionUseCase
		want      *TransactionUseCase
	}{
		{
			name:      "set notifiers",
			notifiers: []INotifier{mocks.NewINotifier(t)},
			object: &TransactionUseCase{
				notifiers: []INotifier{},
			},
			want: &TransactionUseCase{
				notifiers: []INotifier{mocks.NewINotifier(t)},
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
	notifier := mocks.NewINotifier(t)
	uc := TransactionUseCase{
		repo:       transRepo,
		paymentSvc: paymentSvc,
		notifiers:  []INotifier{notifier},
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

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      userID,
			AccountName: "momo",
		}

		walletMock := &entity.Wallet{
			ID:         walletID,
			UserID:     userID,
			WalletName: "quangpn's wallet",
		}

		newTransMock := &entity.Transaction{
			ID:              "t_00001",
			WalletID:        walletID,
			AccountID:       accountID,
			Amount:          amount,
			Currency:        currency,
			TransactionKind: entity.TransactionIn,
			Note:            note,
			Status:          entity.TransactionStatusNew,
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().SaveTransaction(ctx, IsMatchByTransaction(newTransMock)).Return(nil).Once()

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

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(nil, errDB).Once()

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

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(nil, nil).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("account not found"))
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

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      "u_00001",
			AccountName: "momo",
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
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

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      "u_00001",
			AccountName: "momo",
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(nil, nil).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("wallet not found"))
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

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      "u_00001",
			AccountName: "momo",
		}

		walletMock := &entity.Wallet{
			ID:         walletID,
			UserID:     "u_00001",
			WalletName: "quangpn's wallet",
		}

		newTrans := &entity.Transaction{
			ID:              "t_00001",
			WalletID:        walletID,
			AccountID:       accountID,
			Amount:          amount,
			Currency:        currency,
			TransactionKind: entity.TransactionIn,
			Note:            note,
			Status:          entity.TransactionStatusNew,
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().SaveTransaction(ctx, IsMatchByTransaction(newTrans)).Return(errSaveTrans).Once()

		//Act
		err := uc.Deposit(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrCreate(errSaveTrans, "failed to create deposit transaction")
		assert.Equal(t, expectedErr, err)
	})
}

func TestTransactionUseCase_Withdraw(t *testing.T) {
	transRepo := mocks.NewITransactionRepository(t)
	paymentSvc := mocks.NewIPaymentServiceProvider(t)
	notifier := mocks.NewINotifier(t)
	uc := TransactionUseCase{
		repo:       transRepo,
		paymentSvc: paymentSvc,
		notifiers:  []INotifier{notifier},
	}
	t.Run("success", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		userID := "u_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Withdraw 1,000,000 VND"
		balance := 10000000.0

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      userID,
			AccountName: "momo",
		}

		walletMock := &entity.Wallet{
			ID:         walletID,
			UserID:     userID,
			WalletName: "quangpn's wallet",
		}

		newTransMock := &entity.Transaction{
			ID:              "t_00001",
			WalletID:        walletID,
			AccountID:       accountID,
			Amount:          amount,
			Currency:        currency,
			TransactionKind: entity.TransactionOut,
			Note:            note,
			Status:          entity.TransactionStatusNew,
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().GetBalanceByWalletID(ctx, walletID).Return(balance, nil).Once()
		transRepo.EXPECT().SaveTransaction(ctx, IsMatchByTransaction(newTransMock)).Return(nil).Once()

		//Act
		err := uc.Withdraw(ctx, walletID, accountID, amount, currency, note)

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
		note := "Withdraw 1,000,000 VND"
		errDB := fmt.Errorf("unexpected error")

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(nil, errDB).Once()

		//Act
		err := uc.Withdraw(ctx, walletID, accountID, amount, currency, note)

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
		note := "Withdraw 1,000,000 VND"

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(nil, nil).Once()

		//Act
		err := uc.Withdraw(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("account not found"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to get wallet by id", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Withdraw 1,000,000 VND"
		errDB := fmt.Errorf("unexpected error")

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      "u_00001",
			AccountName: "momo",
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(nil, errDB).Once()

		//Act
		err := uc.Withdraw(ctx, walletID, accountID, amount, currency, note)

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
		note := "Withdraw 1,000,000 VND"

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      "u_00001",
			AccountName: "momo",
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(nil, nil).Once()

		//Act
		err := uc.Withdraw(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("wallet not found"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to get balance by wallet id", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Withdraw 1,000,000 VND"
		errDB := fmt.Errorf("unexpected error")

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      "u_00001",
			AccountName: "momo",
		}

		walletMock := &entity.Wallet{
			ID:         walletID,
			UserID:     "u_00001",
			WalletName: "quangpn's wallet",
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().GetBalanceByWalletID(ctx, walletID).Return(0, errDB).Once()

		//Act
		err := uc.Withdraw(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrGet(errDB, "failed to get balance by wallet id")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Withdraw 1,000,000 VND"
		balance := 1000.0

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      "u_00001",
			AccountName: "momo",
		}

		walletMock := &entity.Wallet{
			ID:         walletID,
			UserID:     "u_00001",
			WalletName: "quangpn's wallet",
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().GetBalanceByWalletID(ctx, walletID).Return(balance, nil).Once()

		//Act
		err := uc.Withdraw(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("insufficient balance"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to create withdraw transaction", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		walletID := "w_00001"
		accountID := "a_00001"
		amount := 1000000.0
		currency := "VND"
		note := "Withdraw 1,000,000 VND"
		balance := 10000000.0
		errSaveTrans := fmt.Errorf("unexpected error")

		accountMock := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      "u_00001",
			AccountName: "momo",
		}

		walletMock := &entity.Wallet{
			ID:         walletID,
			UserID:     "u_00001",
			WalletName: "quangpn's wallet",
		}

		newTrans := &entity.Transaction{
			ID:              "t_00001",
			WalletID:        walletID,
			AccountID:       accountID,
			Amount:          amount,
			Currency:        currency,
			TransactionKind: entity.TransactionOut,
			Note:            note,
			Status:          entity.TransactionStatusNew,
		}

		transRepo.EXPECT().GetLinkedAccountByID(ctx, accountID).Return(accountMock, nil).Once()
		transRepo.EXPECT().GetWalletByID(ctx, walletID).Return(walletMock, nil).Once()
		transRepo.EXPECT().GetBalanceByWalletID(ctx, walletID).Return(balance, nil).Once()
		transRepo.EXPECT().SaveTransaction(ctx, IsMatchByTransaction(newTrans)).Return(errSaveTrans).Once()

		//Act
		err := uc.Withdraw(ctx, walletID, accountID, amount, currency, note)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrCreate(errSaveTrans, "failed to create withdraw transaction")
		assert.Equal(t, expectedErr, err)
	})
}

func TestTransactionUseCase_PayTransaction(t *testing.T) {
	transRepo := mocks.NewITransactionRepository(t)
	paymentSvc := mocks.NewIPaymentServiceProvider(t)
	notifier := mocks.NewINotifier(t)
	uc := TransactionUseCase{
		repo:       transRepo,
		paymentSvc: paymentSvc,
		notifiers:  []INotifier{notifier},
	}

	t.Run("success: withdraw", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		trans := &entity.Transaction{
			ID:              "t_00001",
			WalletID:        "w_00001",
			AccountID:       "a_00001",
			Amount:          1000000.0,
			Currency:        "VND",
			TransactionKind: entity.TransactionOut,
			Note:            "Withdraw 1,000,000 VND",
			Status:          entity.TransactionStatusNew,
		}

		transRepo.EXPECT().GetTransactionByID(ctx, trans.ID).Return(trans, nil).Once()

		paymentSvc.EXPECT().Withdraw(ctx, trans.Amount, trans.Currency, trans.Note).Return(nil).Once()

		transRepo.EXPECT().UpdateTransactionStatus(ctx, trans.ID, entity.TransactionStatusSuccessful).Return(nil).Once()

		//Act
		err := uc.PayTransaction(ctx, trans.ID)

		//Assert
		assert.NoError(t, err)
	})

	t.Run("success: deposit", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		trans := &entity.Transaction{
			ID:              "t_00001",
			WalletID:        "w_00001",
			AccountID:       "a_00001",
			Amount:          1000000.0,
			Currency:        "VND",
			TransactionKind: entity.TransactionIn,
			Note:            "Deposit 1,000,000 VND",
			Status:          entity.TransactionStatusNew,
		}

		transRepo.EXPECT().GetTransactionByID(ctx, trans.ID).Return(trans, nil).Once()

		paymentSvc.EXPECT().Deposit(ctx, trans.Amount, trans.Currency, trans.Note).Return(nil).Once()

		transRepo.EXPECT().UpdateTransactionStatus(ctx, trans.ID, entity.TransactionStatusSuccessful).Return(nil).Once()

		//Act
		err := uc.PayTransaction(ctx, trans.ID)

		//Assert
		assert.NoError(t, err)
	})

	t.Run("failed to get transaction by id", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		transID := "t_00001"
		errDB := fmt.Errorf("unexpected error")

		transRepo.EXPECT().GetTransactionByID(ctx, transID).Return(nil, errDB).Once()

		//Act
		err := uc.PayTransaction(ctx, transID)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrGet(errDB, "failed to get transaction by id")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("no transactions found in ready-to-pay status", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		transID := "t_00001"

		transRepo.EXPECT().GetTransactionByID(ctx, transID).Return(nil, nil).Once()

		//Act
		err := uc.PayTransaction(ctx, transID)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("no transactions found in ready-to-pay status"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("transaction status is not new", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		transID := "t_00001"

		trans := &entity.Transaction{
			ID:              transID,
			WalletID:        "w_00001",
			AccountID:       "a_00001",
			Amount:          1000000.0,
			Currency:        "VND",
			TransactionKind: entity.TransactionOut,
			Note:            "Withdraw 1,000,000 VND",
			Status:          entity.TransactionStatusSuccessful,
		}

		transRepo.EXPECT().GetTransactionByID(ctx, transID).Return(trans, nil).Once()

		//Act
		err := uc.PayTransaction(ctx, transID)

		//Assert
		assert.Error(t, err)
		expectedErr := apperror.ErrInvalidParams(fmt.Errorf("transaction status is not new"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failed to withdraw", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		trans := &entity.Transaction{
			ID:              "t_00001",
			WalletID:        "w_00001",
			AccountID:       "a_00001",
			Amount:          1000000.0,
			Currency:        "VND",
			TransactionKind: entity.TransactionOut,
			Note:            "Withdraw 1,000,000 VND",
			Status:          entity.TransactionStatusNew,
		}
		errorWithdraw := fmt.Errorf("unexpected error")

		transRepo.EXPECT().GetTransactionByID(ctx, trans.ID).Return(trans, nil).Once()
		paymentSvc.EXPECT().Withdraw(ctx, trans.Amount, trans.Currency, trans.Note).Return(errorWithdraw).Once()

		transRepo.EXPECT().UpdateTransactionStatus(ctx, trans.ID, entity.TransactionStatusFailed).Return(nil).Once()

		//Act
		err := uc.PayTransaction(ctx, trans.ID)

		//Assert
		assert.NoError(t, err)
	})

	t.Run("failed to deposit", func(t *testing.T) {
		//Arrange
		ctx := context.Background()
		trans := &entity.Transaction{
			ID:              "t_00001",
			WalletID:        "w_00001",
			AccountID:       "a_00001",
			Amount:          1000000.0,
			Currency:        "VND",
			TransactionKind: entity.TransactionIn,
			Note:            "Deposit 1,000,000 VND",
			Status:          entity.TransactionStatusNew,
		}
		errorWithdraw := fmt.Errorf("unexpected error")

		transRepo.EXPECT().GetTransactionByID(ctx, trans.ID).Return(trans, nil).Once()
		paymentSvc.EXPECT().Deposit(ctx, trans.Amount, trans.Currency, trans.Note).Return(errorWithdraw).Once()

		transRepo.EXPECT().UpdateTransactionStatus(ctx, trans.ID, entity.TransactionStatusFailed).Return(nil).Once()

		//Act
		err := uc.PayTransaction(ctx, trans.ID)

		//Assert
		assert.NoError(t, err)
	})
}

func IsMatchByTransaction(a *entity.Transaction) interface{} {
	return mock.MatchedBy(func(b *entity.Transaction) bool {
		return a.WalletID == b.WalletID &&
			a.AccountID == b.AccountID &&
			a.Amount == b.Amount &&
			a.Currency == b.Currency &&
			a.TransactionKind == b.TransactionKind &&
			a.Note == b.Note && a.Status == b.Status
	})
}
