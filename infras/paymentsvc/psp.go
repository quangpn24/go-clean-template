package paymentsvc

import "fmt"

type PaymentServiceProvider struct {
}

func NewPaymentServiceProvider() *PaymentServiceProvider {
	return &PaymentServiceProvider{}
}

func (b *PaymentServiceProvider) Deposit(accountNumber string, bankName string, amount float64, currency string, note string) error {
	//call psp api to deposit
	fmt.Println("Deposit successfully")
	return nil
}

func (b *PaymentServiceProvider) Withdraw(accountNumber string, bankName string, amount float64, currency string, note string) error {
	//call psp api to withdraw
	fmt.Println("Withdraw successfully")
	return nil
}
