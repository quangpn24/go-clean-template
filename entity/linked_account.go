package entity

type LinkedAccount struct {
	ID          string
	UserID      string
	AccountName string
}

func NewLinkedAccount(id string, userId string, accountName string) *LinkedAccount {
	return &LinkedAccount{
		ID:          id,
		UserID:      userId,
		AccountName: accountName,
	}
}
