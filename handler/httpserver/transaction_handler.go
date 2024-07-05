package httpserver

import (
	"go-clean-template/handler/httpserver/model"
	"go-clean-template/pkg/apperror"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterTransactionRoutesV1(group *echo.Group) {
	group.POST("/deposit", s.Deposit)
}

func (s *Server) Deposit(c echo.Context) error {
	var (
		req model.DepositRequest
		ctx = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		return s.handleError(c, apperror.ErrInvalidParams(err))
	}

	if err := req.Validate(); err != nil {
		return s.handleError(c, apperror.ErrInvalidParams(err))
	}

	if err := s.TransactionUseCase.Deposit(ctx, req.WalletID, req.AccountID, req.Amount, req.Currency, req.Note);
		err != nil {
		return s.handleError(c, err)
	}

	return s.handleSuccess(c, http.StatusOK, "OK")
}
