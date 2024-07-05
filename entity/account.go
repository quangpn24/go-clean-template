package entity

type Account struct {
	ID            string
	UserID        string
	BankName      string
	AccountNumber string
	IsLinked      bool
}

func NewAccount(id string, userId string, accountNumber string) *Account {
	return &Account{
		ID:            id,
		UserID:        userId,
		AccountNumber: accountNumber,
	}
}
