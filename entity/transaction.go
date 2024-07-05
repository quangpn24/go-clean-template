package entity

type TransactionKind string

const (
	TransactionIn  TransactionKind = "in"
	TransactionOut TransactionKind = "out"
)

type TransactionCategory string

const (
	CategoryDeposit  TransactionCategory = "DEPOSIT"
	CategoryWithdraw TransactionCategory = "WITHDRAW"
)

type Transaction struct {
	ID               string
	SenderWalletID   string
	ReceiverWalletID string
	AccountID        string
	Amount           float64
	Currency         string
	Category         TransactionCategory
	TransactionKind  TransactionKind
	Note             string
}

func NewTransaction(id string, sender string, receiver string, accountID string, amount float64, currency string, kind TransactionKind, category TransactionCategory, note string) *Transaction {
	return &Transaction{
		ID:               id,
		SenderWalletID:   sender,
		ReceiverWalletID: receiver,
		AccountID:        accountID,
		Amount:           amount,
		Currency:         currency,
		Category:         category,
		TransactionKind:  kind,
		Note:             note,
	}
}
