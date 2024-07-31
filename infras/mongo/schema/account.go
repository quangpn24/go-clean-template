package schema

import (
	"time"

	"go-clean-template/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkedAccountSchema struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id,omitempty"`
	AccountName string             `bson:"account_name,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

func (a *LinkedAccountSchema) ToLinkedAccount() *entity.LinkedAccount {
	return &entity.LinkedAccount{
		ID:          a.ID.String(),
		UserID:      a.UserID,
		AccountName: a.AccountName,
	}
}
