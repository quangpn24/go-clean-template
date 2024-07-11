package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-clean-template/handler/httpserver/model"
	"go-clean-template/pkg/testutil"
	"go-clean-template/usecase/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupDeposit(t testing.TB, req interface{}) (echo.Context, *httptest.ResponseRecorder) {
	body, err := json.Marshal(req)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/v1/transactions/deposit", bytes.NewReader(body))
	r.Header.Set("Content-type", echo.MIMEApplicationJSON)
	r.Header.Set("User-agent", "testing")
	w := httptest.NewRecorder()

	return echo.New().NewContext(r, w), w
}

func TestServer_Deposit(t *testing.T) {
	db := testutil.CreateConnection(t, "test1", "user", "123456")
	testutil.MigrateTestDatabase(t, db, "../../migrations")

	transUCMock := mocks.NewITransactionUseCase(t)
	s := Server{
		TransactionUseCase: transUCMock,
		Logger:             zap.S(),
	}

	t.Run("200: success", func(t *testing.T) {
		// Arrange
		req := model.DepositRequest{
			WalletID:  "wallet1",
			AccountID: "account1",
			Amount:    1000,
			Currency:  "USD",
			Note:      "deposit",
		}
		c, resp := setupDeposit(t, req)
		transUCMock.EXPECT().Deposit(c.Request().Context(), req.WalletID, req.AccountID, req.Amount, req.Currency,
			req.Note).Return(nil).Once()

		// Act
		err := s.Deposit(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		expectedData := "OK"
		actual := extractSuccessData[string](t, resp.Body)
		assert.Equal(t, expectedData, actual)
	})

	t.Run("400: failed to bind", func(t *testing.T) {
		// Arrange
		req := map[string]interface{}{
			"wallet_id":  "w-001",
			"account_id": "a_001",
			"amount":     "1000abc",
			"currency":   "USD",
			"note":       "deposit",
		}
		c, resp := setupDeposit(t, req)

		// Act
		err := s.Deposit(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		expectedData := "invalid params"
		actual := extractErrorData(t, resp.Body)
		assert.Equal(t, expectedData, actual.Message)
	})

	t.Run("400: failed to validate", func(t *testing.T) {
		// Arrange
		req := model.DepositRequest{
			WalletID:  "w_001",
			AccountID: "a_001",
			Amount:    -10,
			Currency:  "USD",
			Note:      "deposit",
		}
		c, resp := setupDeposit(t, req)

		// Act
		err := s.Deposit(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		expectedData := "invalid params"
		actual := extractErrorData(t, resp.Body)
		assert.Equal(t, expectedData, actual.Message)
	})

	t.Run("500: failed to deposit", func(t *testing.T) {
		// Arrange
		req := model.DepositRequest{
			WalletID:  "w_001",
			AccountID: "a_001",
			Amount:    1000,
			Currency:  "USD",
			Note:      "deposit",
		}
		c, resp := setupDeposit(t, req)
		transUCMock.EXPECT().Deposit(c.Request().Context(), req.WalletID, req.AccountID, req.Amount, req.Currency,
			req.Note).Return(fmt.Errorf("unexpected error")).Once()

		// Act
		err := s.Deposit(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		expectedData := "unexpected error"
		actual := extractErrorData(t, resp.Body)
		assert.Equal(t, expectedData, actual.Message)
	})
}
