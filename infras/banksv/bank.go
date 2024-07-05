package banksv

type BankService struct {
}

func NewBankService() *BankService {
	return &BankService{}
}

func (b *BankService) Deposit(accountNumber string, bankName string, amount float64, currency string, note string) error {
	//call bank api to deposit
	return nil
}
