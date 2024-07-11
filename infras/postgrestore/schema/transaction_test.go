package schema

import (
	"reflect"
	"testing"
	"time"

	"go-clean-template/entity"
)

func TestTransactionSchema_ToTransaction(t *testing.T) {
	type fields struct {
		ID               string
		SenderWalletID   *string
		ReceiverWalletID *string
		AccountID        *string
		Amount           float64
		Currency         string
		Category         string
		Note             string
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}
	receiverWalletID := "w_001"
	senderWalletID := "w_002"
	accountID := "a_001"
	tests := []struct {
		name   string
		fields fields
		want   *entity.Transaction
	}{
		{
			name: "Sender nil",
			fields: fields{
				ID:               "1",
				SenderWalletID:   nil,
				ReceiverWalletID: &receiverWalletID,
				AccountID:        &accountID,
				Amount:           100,
				Currency:         "USD",
				Category:         "DEPOSIT",
				Note:             "deposit 100",
			},
			want: &entity.Transaction{
				ID:               "1",
				SenderWalletID:   "",
				ReceiverWalletID: receiverWalletID,
				AccountID:        accountID,
				Amount:           100,
				Currency:         "USD",
				Category:         entity.CategoryDeposit,
				Note:             "deposit 100",
			},
		},
		{
			name: "receiver nil",
			fields: fields{
				ID:               "1",
				SenderWalletID:   &senderWalletID,
				ReceiverWalletID: nil,
				AccountID:        &accountID,
				Amount:           100,
				Currency:         "USD",
				Category:         "WITHDRAW",
				Note:             "withdraw 100",
			},
			want: &entity.Transaction{
				ID:               "1",
				SenderWalletID:   senderWalletID,
				ReceiverWalletID: "",
				AccountID:        accountID,
				Amount:           100,
				Currency:         "USD",
				Category:         entity.CategoryWithdraw,
				Note:             "withdraw 100",
			},
		},
		{
			name: "accountID nil",
			fields: fields{
				ID:               "1",
				SenderWalletID:   &senderWalletID,
				ReceiverWalletID: &receiverWalletID,
				AccountID:        nil,
				Amount:           100,
				Currency:         "USD",
				Category:         "TRANSFER",
				Note:             "transfer 100",
			},
			want: &entity.Transaction{
				ID:               "1",
				SenderWalletID:   senderWalletID,
				ReceiverWalletID: receiverWalletID,
				AccountID:        "",
				Amount:           100,
				Currency:         "USD",
				Category:         entity.CategoryTransfer,
				Note:             "transfer 100",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trans := &TransactionSchema{
				ID:               tt.fields.ID,
				SenderWalletID:   tt.fields.SenderWalletID,
				ReceiverWalletID: tt.fields.ReceiverWalletID,
				AccountID:        tt.fields.AccountID,
				Amount:           tt.fields.Amount,
				Currency:         tt.fields.Currency,
				Category:         tt.fields.Category,
				Note:             tt.fields.Note,
				CreatedAt:        tt.fields.CreatedAt,
				UpdatedAt:        tt.fields.UpdatedAt,
			}
			if got := trans.ToTransaction(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToTransactionSchema(t *testing.T) {
	type args struct {
		trans *entity.Transaction
	}
	receiverWalletID := "w_001"
	senderWalletID := "w_002"
	accountID := "a_001"

	tests := []struct {
		name string
		args args
		want *TransactionSchema
	}{
		{
			name: "Sender is empty",
			args: args{
				&entity.Transaction{
					ID:               "1",
					SenderWalletID:   "",
					ReceiverWalletID: receiverWalletID,
					AccountID:        accountID,
					Amount:           100,
					Currency:         "USD",
					Category:         entity.CategoryDeposit,
					Note:             "deposit 100",
				},
			},
			want: &TransactionSchema{
				ID:               "1",
				SenderWalletID:   nil,
				ReceiverWalletID: &receiverWalletID,
				AccountID:        &accountID,
				Amount:           100,
				Currency:         "USD",
				Category:         "DEPOSIT",
				Note:             "deposit 100",
			},
		},
		{
			name: "Receiver is empty",
			args: args{
				&entity.Transaction{
					ID:               "1",
					SenderWalletID:   senderWalletID,
					ReceiverWalletID: "",
					AccountID:        accountID,
					Amount:           100,
					Currency:         "USD",
					Category:         entity.CategoryWithdraw,
					Note:             "withdraw 100",
				},
			},
			want: &TransactionSchema{
				ID:               "1",
				SenderWalletID:   &senderWalletID,
				ReceiverWalletID: nil,
				AccountID:        &accountID,
				Amount:           100,
				Currency:         "USD",
				Category:         "WITHDRAW",
				Note:             "withdraw 100",
			},
		},
		{
			name: "account is empty",
			args: args{
				&entity.Transaction{
					ID:               "1",
					SenderWalletID:   senderWalletID,
					ReceiverWalletID: receiverWalletID,
					AccountID:        "",
					Amount:           100,
					Currency:         "USD",
					Category:         entity.CategoryTransfer,
					Note:             "transfer 100",
				},
			},
			want: &TransactionSchema{
				ID:               "1",
				SenderWalletID:   &senderWalletID,
				ReceiverWalletID: &receiverWalletID,
				AccountID:        nil,
				Amount:           100,
				Currency:         "USD",
				Category:         "TRANSFER",
				Note:             "transfer 100",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTransactionSchema(tt.args.trans); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTransactionSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
