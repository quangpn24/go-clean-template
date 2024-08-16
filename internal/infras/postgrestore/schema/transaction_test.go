package schema

import (
	"reflect"
	"testing"
	"time"

	"go-clean-template/internal/entity"
)

func TestTransactionSchema_ToTransaction(t *testing.T) {
	type fields struct {
		ID              string
		WalletID        string
		AccountID       string
		Amount          float64
		Currency        string
		TransactionKind string
		Status          string
		Note            string
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}
	walletID := "w_001"
	accountID := "a_001"
	tests := []struct {
		name   string
		fields fields
		want   *entity.Transaction
	}{
		{
			name: "test to transaction",
			fields: fields{
				ID:              "1",
				WalletID:        walletID,
				AccountID:       accountID,
				Amount:          100,
				Currency:        "USD",
				TransactionKind: "IN",
				Status:          "NEW",
				Note:            "deposit 100",
			},
			want: &entity.Transaction{
				ID:              "1",
				WalletID:        walletID,
				AccountID:       accountID,
				Amount:          100,
				Currency:        "USD",
				TransactionKind: entity.TransactionIn,
				Status:          entity.TransactionStatusNew,
				Note:            "deposit 100",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trans := &TransactionSchema{
				ID:              tt.fields.ID,
				WalletID:        tt.fields.WalletID,
				AccountID:       tt.fields.AccountID,
				Amount:          tt.fields.Amount,
				Currency:        tt.fields.Currency,
				TransactionKind: tt.fields.TransactionKind,
				Status:          tt.fields.Status,
				Note:            tt.fields.Note,
				CreatedAt:       tt.fields.CreatedAt,
				UpdatedAt:       tt.fields.UpdatedAt,
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
	walletID := "w_001"
	accountID := "a_001"

	tests := []struct {
		name string
		args args
		want *TransactionSchema
	}{
		{
			name: "To TransactionSchema",
			args: args{
				&entity.Transaction{
					ID:              "1",
					WalletID:        walletID,
					AccountID:       accountID,
					Amount:          100,
					Currency:        "USD",
					TransactionKind: entity.TransactionOut,
					Status:          entity.TransactionStatusSuccessful,
					Note:            "Withdraw 100",
				},
			},
			want: &TransactionSchema{
				ID:              "1",
				WalletID:        walletID,
				AccountID:       accountID,
				Amount:          100,
				Currency:        "USD",
				TransactionKind: "OUT",
				Status:          "SUCCESSFUL",
				Note:            "Withdraw 100",
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
