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
			ID:         walletID,
			UserID:     userId,
			WalletName: "My Wallet",
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
		got, err := repo.GetWalletByID(ctx, "w_0002")

		//Assert
		assert.NoError(t, err)
		assert.Nil(t, got)
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

		account := &schema.LinkedAccountSchema{
			ID:          accountID,
			UserID:      userId,
			AccountName: "momo",
		}
		assert.NoError(t, repo.db.Table(LinkedAccountTable).Create(account).Error)

		wallet := &schema.WalletSchema{
			ID:         walletID,
			UserID:     userId,
			WalletName: "My wallet",
		}
		assert.NoError(t, repo.db.Table(WalletTable).Create(wallet).Error)

		want := entity.NewTransaction(transID, walletID, accountID, 1000,
			"USD", entity.TransactionIn, "", entity.TransactionStatusNew)

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

		want := &entity.LinkedAccount{
			ID:          accountID,
			UserID:      userId,
			AccountName: "momo",
		}
		assert.NoError(t, repo.db.Table(LinkedAccountTable).Create(want).Error)
		//Act
		got, err := repo.GetLinkedAccountByID(ctx, accountID)

		//Assert
		assert.NoError(t, err)
		assertAccount(t, want, got)
	})

	t.Run("record not found", func(t *testing.T) {
		//Act
		got, err := repo.GetLinkedAccountByID(ctx, "acc_0002")

		//Assert
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestTransactionRepo_GetBalanceByWalletID(t *testing.T) {
	db := testutil.CreateConnection(t, "test1", "test1", "123456")
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	repo := NewTransactionRepo(db)
	ctx := context.Background()

	t.Run("success: get balance by wallet id", func(t *testing.T) {
		//Arrange
		walletID := "1"
		userId := uuid.New().String()

		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
		VALUES (?, 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		assert.NoError(t, repo.db.Exec(query, userId).Error)

		wallet := &schema.WalletSchema{
			ID:         walletID,
			UserID:     userId,
			WalletName: "My wallet",
		}

		assert.NoError(t, repo.db.Table(WalletTable).Create(wallet).Error)

		linkedAccount := &schema.LinkedAccountSchema{
			ID:          "acc_0001",
			UserID:      userId,
			AccountName: "momo",
		}
		assert.NoError(t, repo.db.Table(LinkedAccountTable).Create(linkedAccount).Error)

		transIn := entity.NewTransaction(uuid.New().String(), walletID, "acc_0001", 1000, "VND", entity.TransactionIn, "", entity.TransactionStatusSuccessful)
		transOut := entity.NewTransaction(uuid.New().String(), walletID, "acc_0001", 500, "VND", entity.TransactionOut, "", entity.TransactionStatusSuccessful)

		transInSchema := schema.ToTransactionSchema(transIn)
		transOutSchema := schema.ToTransactionSchema(transOut)

		assert.NoError(t, repo.db.Table(TransactionsTable).Create(transInSchema).Error)
		assert.NoError(t, repo.db.Table(TransactionsTable).Create(transOutSchema).Error)

		//Act
		got, err := repo.GetBalanceByWalletID(ctx, walletID)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, 500.0, got)
	})
}

func TestTransactionRepo_GetTransactionByID(t *testing.T) {
	db := testutil.CreateConnection(t, "test1", "test1", "123456")
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	repo := NewTransactionRepo(db)
	ctx := context.Background()

	t.Run("success: get transaction by id", func(t *testing.T) {
		//Arrange
		transID := "t_001"
		walletID := "1"
		userId := uuid.New().String()

		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
		VALUES (?, 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		assert.NoError(t, repo.db.Exec(query, userId).Error)

		wallet := &schema.WalletSchema{
			ID:         walletID,
			UserID:     userId,
			WalletName: "My wallet",
		}

		linkedAccount := &schema.LinkedAccountSchema{
			ID:          "acc_0001",
			UserID:      userId,
			AccountName: "momo",
		}
		assert.NoError(t, repo.db.Table(LinkedAccountTable).Create(linkedAccount).Error)

		assert.NoError(t, repo.db.Table(WalletTable).Create(wallet).Error)

		transIn := entity.NewTransaction(transID, walletID, "acc_0001", 1000, "VND", entity.TransactionIn, "", entity.TransactionStatusNew)
		transSchema := schema.ToTransactionSchema(transIn)
		assert.NoError(t, repo.db.Table(TransactionsTable).Create(transSchema).Error)

		//Act
		got, err := repo.GetTransactionByID(ctx, transID)

		//Assert
		assert.NoError(t, err)
		assertTransaction(t, transIn, got)
	})
}

func TestTransactionRepo_UpdateTransactionStatus(t *testing.T) {
	db := testutil.CreateConnection(t, "test1", "test1", "123456")
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	repo := NewTransactionRepo(db)
	ctx := context.Background()

	t.Run("success: update transaction status", func(t *testing.T) {
		//Arrange
		transID := "t_001"
		walletID := "1"
		userId := uuid.New().String()

		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
    		VALUES (?, 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		assert.NoError(t, repo.db.Exec(query, userId).Error)

		wallet := &schema.WalletSchema{
			ID:         walletID,
			UserID:     userId,
			WalletName: "My wallet",
		}

		assert.NoError(t, repo.db.Table(WalletTable).Create(wallet).Error)

		linkedAccount := &schema.LinkedAccountSchema{
			ID:          "acc_0001",
			UserID:      userId,
			AccountName: "momo",
		}
		assert.NoError(t, repo.db.Table(LinkedAccountTable).Create(linkedAccount).Error)

		transSchema := schema.TransactionSchema{
			ID:              transID,
			WalletID:        walletID,
			AccountID:       "acc_0001",
			Amount:          1000,
			Currency:        "VND",
			TransactionKind: string(entity.TransactionIn),
			Note:            "",
			Status:          string(entity.TransactionStatusNew),
		}

		assert.NoError(t, repo.db.Table(TransactionsTable).Create(&transSchema).Error)

		//Act
		err := repo.UpdateTransactionStatus(ctx, transID, entity.TransactionStatusSuccessful)

		//Assert
		assert.NoError(t, err)
		var got *entity.Transaction
		assert.NoError(t, repo.db.Raw("SELECT * from transactions").Scan(&got).Error)
		assert.Equal(t, entity.TransactionStatusSuccessful, got.Status)
	})
}

func assertWallet(t testing.TB, want *entity.Wallet, got *entity.Wallet) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.UserID, got.UserID)
	assert.Equal(t, want.WalletName, got.WalletName)
}

func assertTransaction(t testing.TB, want *entity.Transaction, got *entity.Transaction) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.WalletID, got.WalletID)
	assert.Equal(t, want.AccountID, got.AccountID)
	assert.Equal(t, want.Amount, got.Amount)
	assert.Equal(t, want.Currency, got.Currency)
	assert.Equal(t, want.TransactionKind, got.TransactionKind)
	assert.Equal(t, want.Note, got.Note)
	assert.Equal(t, want.Status, got.Status)
}

func assertAccount(t testing.TB, want *entity.LinkedAccount, got *entity.LinkedAccount) {
	t.Helper()

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.UserID, got.UserID)
	assert.Equal(t, want.AccountName, got.AccountName)
}
