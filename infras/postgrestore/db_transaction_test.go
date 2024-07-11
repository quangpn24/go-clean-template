package postgrestore

import (
	"context"
	"testing"

	"go-clean-template/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestDBTransaction_Begin(t *testing.T) {
	// Arrange
	dbName, dbUser, dbPass := "test1", "user1", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	trans := NewDBTransaction(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {

		//Act
		iTrans, err := trans.Begin(ctx)

		//Assert
		assert.NoError(t, err)
		assert.NotNil(t, iTrans)
	})
}

func TestDBTransaction_Commit(t *testing.T) {
	// Arrange
	dbName, dbUser, dbPass := "test1", "user1", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	trans := NewDBTransaction(db)
	ctx := context.Background()

	t.Run("commit success", func(t *testing.T) {
		//Arrange
		iTrans, err := trans.Begin(ctx)
		assert.NoError(t, err)

		tr := iTrans.(*DBTransaction)

		// create new user transaction to test commit
		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
        VALUES (gen_random_uuid(), 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		err = tr.db.Exec(query).Error
		assert.NoError(t, err)

		//Act
		err = iTrans.Commit(ctx)

		//Assert
		assert.NoError(t, err)
		var count int
		trans.db.Raw("SELECT count(*) from users").Scan(&count)
		assert.Equal(t, 1, count)
	})
}

func TestDBTransaction_Rollback(t *testing.T) {
	// Arrange
	dbName, dbUser, dbPass := "test1", "user1", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	trans := NewDBTransaction(db)
	ctx := context.Background()

	t.Run("rollback success", func(t *testing.T) {
		//Arrange
		iTrans, err := trans.Begin(ctx)
		assert.NoError(t, err)

		tr := iTrans.(*DBTransaction)

		// create new user transaction to test commit
		query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
        VALUES (gen_random_uuid(), 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
		err = tr.db.Exec(query).Error
		assert.NoError(t, err)

		//Act
		iTrans.Rollback(ctx)

		//Assert
		assert.NoError(t, err)
		var count int
		trans.db.Raw("SELECT count(*) from users").Scan(&count)
		assert.Equal(t, 0, count)
	})
}
