package model

import (
	"github.com/go-playground/validator/v10"
)

type DepositRequest struct {
	WalletID  string  `json:"wallet_id"`
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	Currency  string  `json:"currency"`
	Note      string  `json:"note"`
}

func (r DepositRequest) Validate() error {
	v := validator.New()
	err := v.Struct(r)
	return err
}
