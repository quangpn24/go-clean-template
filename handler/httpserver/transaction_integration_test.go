package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-clean-template/handler/httpserver/model"
	"go-clean-template/infras/notification"
	"go-clean-template/infras/paymentsvc"
	"go-clean-template/infras/postgrestore"
	"go-clean-template/infras/postgrestore/schema"
	"go-clean-template/pkg/config"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/testutil"
	"go-clean-template/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDepositAPI(t testing.TB, req model.DepositRequest) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	body, err := json.Marshal(req)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/v1/transactions/deposit", bytes.NewReader(body))
	r.Header.Set("Content-type", echo.MIMEApplicationJSON)
	r.Header.Set("User-agent", "testing")
	w := httptest.NewRecorder()

	return r, w
}

func newTransactionServerForTest(t testing.TB, db *gorm.DB) Server {
	t.Helper()

	transRepo := postgrestore.NewTransactionRepo(db)
	paymentSvc := paymentsvc.NewPaymentServiceProvider()
	dbTransaction := postgrestore.NewDBTransaction(db)
	transUseCase := usecase.NewTransactionUseCase(transRepo, paymentSvc, dbTransaction)
	transUseCase.SetNotifiers(notification.NewEmailNotifier(), notification.NewAppNotifier())

	router := echo.New()

	applog, err := logger.NewAppLogger()
	assert.NoError(t, err)

	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	s := Server{
		Router:             router,
		TransactionUseCase: transUseCase,
		Logger:             applog,
		Config:             cfg,
	}
	s.RegisterTransactionRoutesV1(router.Group("/api/v1/transactions"))
	return s
}

func initDataForDeposit(t testing.TB, db *gorm.DB) (*schema.WalletSchema, *schema.AccountSchema) {
	t.Helper()

	userId := uuid.New().String()
	w := &schema.WalletSchema{
		ID:       uuid.New().String(),
		UserID:   userId,
		Balance:  100000,
		Currency: "VND",
	}

	//create user
	query := `INSERT INTO users (id,full_name, email, phone_number,current_address)
        VALUES (?, 'Phan Ngoc Quang', 'quangpn@tm.teqn.asia', '0123456789', 'HCM')`
	err := db.Exec(query, userId).Error
	assert.NoError(t, err)

	//create wallet
	err = db.Table(postgrestore.WalletTable).Create(&w).Error
	assert.NoError(t, err)

	//account
	a := &schema.AccountSchema{
		ID:            uuid.New().String(),
		UserID:        userId,
		BankName:      "Vietcombank",
		AccountNumber: "123456789",
		IsLinked:      true,
	}
	assert.NoError(t, db.Table(postgrestore.AccountTable).Create(&a).Error)

	return w, a
}

func TestDepositAPI(t *testing.T) {
	dbName, dbUser, dbPass := "server", "server", "123456"
	db := testutil.CreateConnection(t, dbName, dbUser, dbPass)
	testutil.MigrateTestDatabase(t, db, "../../migrations")
	s := newTransactionServerForTest(t, db)

	t.Run("deposit successfully", func(t *testing.T) {
		// Arrange
		wallet, account := initDataForDeposit(t, db)

		req := model.DepositRequest{
			WalletID:  wallet.ID,
			AccountID: account.ID,
			Amount:    100000,
			Currency:  "USD",
			Note:      "Deposit 100000 VND",
		}

		request, resp := setupTestDepositAPI(t, req)

		// Act
		s.ServeHTTP(resp, request)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
		var walletAfterDeposit schema.WalletSchema
		assert.NoError(t, db.Table(postgrestore.WalletTable).Where("id = ?", wallet.ID).Take(&walletAfterDeposit).Error)
		assert.Equal(t, req.Amount+wallet.Balance, walletAfterDeposit.Balance)
	})
}
