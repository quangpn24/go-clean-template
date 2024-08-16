package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-clean-template/internal/entity"
	"go-clean-template/internal/handler/httpserver/model"
	"go-clean-template/internal/infras/notification"
	"go-clean-template/internal/infras/paymentsvc"
	"go-clean-template/internal/infras/postgrestore"
	"go-clean-template/internal/infras/postgrestore/schema"
	"go-clean-template/internal/usecase"
	"go-clean-template/pkg/config"
	"go-clean-template/pkg/logger"
	"go-clean-template/pkg/testutil"

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
	transUseCase := usecase.NewTransactionUseCase(transRepo, paymentSvc)
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

func initDataForDeposit(t testing.TB, db *gorm.DB) (*schema.WalletSchema, *schema.LinkedAccountSchema) {
	t.Helper()

	userId := uuid.New().String()
	w := &schema.WalletSchema{
		ID:         uuid.New().String(),
		UserID:     userId,
		WalletName: "My Wallet",
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
	a := &schema.LinkedAccountSchema{
		ID:          uuid.New().String(),
		UserID:      userId,
		AccountName: "My Account",
	}
	assert.NoError(t, db.Table(postgrestore.LinkedAccountTable).Create(&a).Error)

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
		var trans *schema.TransactionSchema
		assert.NoError(t, db.Table(postgrestore.TransactionsTable).
			Where("wallet_id = ? AND account_id = ?", wallet.ID, account.ID).
			Take(&trans).Error)
		assert.NotNil(t, trans)
		assert.Equal(t, string(entity.TransactionStatusNew), trans.Status)
		assert.Equal(t, string(entity.TransactionIn), trans.TransactionKind)
	})
}
