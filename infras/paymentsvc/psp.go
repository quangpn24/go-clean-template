package paymentsvc

type PaymentServiceProvider struct {
}

func NewPaymentServiceProvider() *PaymentServiceProvider {
	return &PaymentServiceProvider{}
}

func (b *PaymentServiceProvider) Deposit(accountNumber string, bankName string, amount float64, currency string, note string) error {
	//call psp api to deposit
	return nil
}
