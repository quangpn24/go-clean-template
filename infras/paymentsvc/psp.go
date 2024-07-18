package paymentsvc

import (
	"context"
	"fmt"
)

type PaymentServiceProvider struct {
}

func NewPaymentServiceProvider() *PaymentServiceProvider {
	return &PaymentServiceProvider{}
}

func (b *PaymentServiceProvider) Deposit(ctx context.Context, amount float64, currency string, note string) error {
	//call psp api to deposit
	fmt.Printf("Deposit %.2f %s successfully\n", amount, currency)
	return nil
}

func (b *PaymentServiceProvider) Withdraw(ctx context.Context, amount float64, currency string, note string) error {
	//call psp api to withdraw
	fmt.Printf("Withdraw %.2f %s successfully\n", amount, currency)
	return nil
}
