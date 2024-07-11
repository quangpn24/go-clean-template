package postgrestore

import (
	"context"
	"testing"

	"go-clean-template/entity"
	"go-clean-template/infras/postgrestore/schema"
	"go-clean-template/pkg/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepo_WithDBTransaction(t *testing.T) {
	//Arrange
	db := testutil.CreateConnection(t, "test1", "test1", "123456")
	repo := NewTransactionRepo(db)

	t.Run("success: is DBTransaction", func(t *testing.T) {
		//Act
		iRepo := repo.WithDBTransaction(NewDBTransaction(db))

		//Assert
		assert.NotNil(t, iRepo)
	})

	t.Run("fail: not DBTransaction", func(t *testing.T) {
		//Act
		iRepo := repo.WithDBTransaction(nil)

		//Assert
		assert.Equal(t, iRepo, repo)
	})
}

func TestTransactionRepo_GetWalletByID(t *testing.T) {
	db := testutil.CreateConnection(t, "test1", "test1", "123456")
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	repo := NewTransactionRepo(db)
	ctx := context.Background()

	t.Run("success: get wallet by id", func(t *testing.T) {
		//Arrange
		walletID := "1"
		userId := uuid.New().String()
		want := &entity.Wallet{
			ID:       walletID,
			UserID:   userId,
			Balance:  1000,
			Currency: "USD",
		}

		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
        VALUES (?, 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		err := repo.db.Exec(query, userId).Error
		assert.NoError(t, err)

		err = repo.db.Table(WalletTable).Create(&want).Error
		assert.NoError(t, err)

		//Act
		got, err := repo.GetWalletByID(ctx, walletID)

		//Assert
		assert.NoError(t, err)
		assertWallet(t, want, got)
	})

	t.Run("record not found", func(t *testing.T) {
		//Act
		got, err := repo.GetWalletByID(ctx, "w_0001")

		//Assert
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestTransactionRepo_UpdateWalletBalance(t *testing.T) {
	db := testutil.CreateConnection(t, "test1", "test1", "123456")
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	repo := NewTransactionRepo(db)
	ctx := context.Background()

	t.Run("success: update wallet balance", func(t *testing.T) {
		//Arrange
		walletID := "1"
		userId := uuid.New().String()
		want := &entity.Wallet{
			ID:       walletID,
			UserID:   userId,
			Balance:  1000,
			Currency: "USD",
		}

		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
        VALUES (?, 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		err := repo.db.Exec(query, userId).Error
		assert.NoError(t, err)

		err = repo.db.Table(WalletTable).Create(&want).Error
		assert.NoError(t, err)
		//Act
		err = repo.UpdateWalletBalance(ctx, walletID, 2000.0)

		//Assert
		assert.NoError(t, err)
		got, _ := repo.GetWalletByID(ctx, walletID)
		assert.Equal(t, 2000.0, got.Balance)

	})
}

func TestTransactionRepo_SaveTransaction(t *testing.T) {
	db := testutil.CreateConnection(t, "test1", "test1", "123456")
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	repo := NewTransactionRepo(db)
	ctx := context.Background()

	t.Run("success: save transaction", func(t *testing.T) {
		//Arrange
		transID := "t_001"
		accountID := "acc_0001"
		walletID := "1"
		userId := uuid.New().String()

		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
        VALUES (?, 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		assert.NoError(t, repo.db.Exec(query, userId).Error)

		account := &schema.AccountSchema{
			ID:            accountID,
			UserID:        userId,
			BankName:      "Vietcombank",
			AccountNumber: "123456789",
			IsLinked:      true,
		}
		assert.NoError(t, repo.db.Table(AccountTable).Create(account).Error)

		wallet := &schema.WalletSchema{
			ID:       walletID,
			UserID:   userId,
			Balance:  1000,
			Currency: "USD",
		}
		assert.NoError(t, repo.db.Table(WalletTable).Create(wallet).Error)

		want := entity.NewTransaction(transID, "", walletID, accountID, 1000,
			"USD", entity.CategoryDeposit, "")

		//Act
		err := repo.SaveTransaction(ctx, want)

		//Assert
		assert.NoError(t, err)
		var got *entity.Transaction
		assert.NoError(t, repo.db.Raw("SELECT * from transactions").Scan(&got).Error)
		assertTransaction(t, want, got)
	})
}

func TestTransactionRepo_GetAccountByID(t *testing.T) {
	db := testutil.CreateConnection(t, "test1", "test1", "123456")
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	repo := NewTransactionRepo(db)
	ctx := context.Background()

	t.Run("success: get account by id", func(t *testing.T) {
		//Arrange
		accountID := "acc_0001"
		userId := uuid.New().String()

		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
        VALUES (?, 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		assert.NoError(t, repo.db.Exec(query, userId).Error)

		want := &entity.Account{
			ID:            accountID,
			UserID:        userId,
			BankName:      "Vietcombank",
			AccountNumber: "123456789",
			IsLinked:      true,
		}
		assert.NoError(t, repo.db.Table(AccountTable).Create(want).Error)
		//Act
		got, err := repo.GetAccountByID(ctx, accountID)

		//Assert
		assert.NoError(t, err)
		assertAccount(t, want, got)
	})

	t.Run("record not found", func(t *testing.T) {
		//Act
		got, err := repo.GetAccountByID(ctx, "acc_0001")

		//Assert
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}

func assertWallet(t testing.TB, want *entity.Wallet, got *entity.Wallet) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.UserID, got.UserID)
	assert.Equal(t, want.Balance, got.Balance)
	assert.Equal(t, want.Currency, got.Currency)
}

func assertTransaction(t testing.TB, want *entity.Transaction, got *entity.Transaction) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.SenderWalletID, got.SenderWalletID)
	assert.Equal(t, want.ReceiverWalletID, got.ReceiverWalletID)
	assert.Equal(t, want.AccountID, got.AccountID)
	assert.Equal(t, want.Amount, got.Amount)
	assert.Equal(t, want.Currency, got.Currency)
	assert.Equal(t, want.Category, got.Category)
	assert.Equal(t, want.Note, got.Note)
}

func assertAccount(t testing.TB, want *entity.Account, got *entity.Account) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.UserID, got.UserID)
	assert.Equal(t, want.BankName, got.BankName)
	assert.Equal(t, want.AccountNumber, got.AccountNumber)
	assert.Equal(t, want.IsLinked, got.IsLinked)
}
