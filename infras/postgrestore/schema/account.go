package schema

import (
	"go-clean-template/entity"
	"time"
)

type AccountSchema struct {
	ID            string    `gorm:"column:id;primaryKey"`
	UserID        string    `gorm:"column:user_id;not null"`
	BankName      string    `gorm:"column:bank_name;not null"`
	AccountNumber string    `gorm:"column:account_number;not null"`
	IsLinked      bool      `gorm:"column:is_linked;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (*AccountSchema) TableName() string {
	return "accounts"
}

func (a *AccountSchema) ToAccount() *entity.Account {
	return &entity.Account{
		ID:            a.ID,
		UserID:        a.UserID,
		BankName:      a.BankName,
		AccountNumber: a.AccountNumber,
		IsLinked:      a.IsLinked,
	}
}
