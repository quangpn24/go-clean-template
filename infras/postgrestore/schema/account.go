package schema

import (
	"time"

	"go-clean-template/entity"
)

type LinkedAccountSchema struct {
	ID          string    `gorm:"column:id;primaryKey"`
	UserID      string    `gorm:"column:user_id;not null"`
	AccountName string    `gorm:"column:account_name;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (*LinkedAccountSchema) TableName() string {
	return "linked_accounts"
}

func (a *LinkedAccountSchema) ToLinkedAccount() *entity.LinkedAccount {
	return &entity.LinkedAccount{
		ID:          a.ID,
		UserID:      a.UserID,
		AccountName: a.AccountName,
	}
}
