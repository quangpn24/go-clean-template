package schema

import (
	"time"

	"go-clean-template/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletSchema struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id,omitempty"`
	WalletName string             `bson:"wallet_name,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `bson:"updated_at,omitempty"`
}

func (w *WalletSchema) ToWallet() *entity.Wallet {
	return &entity.Wallet{
		ID:         w.ID.String(),
		UserID:     w.UserID,
		WalletName: w.WalletName,
	}
}
