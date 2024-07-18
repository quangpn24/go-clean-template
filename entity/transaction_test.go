package entity

import (
	"reflect"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	type args struct {
		id        string
		walletID  string
		accountID string
		amount    float64
		currency  string
		kind      TransactionKind
		note      string
		status    TransactionStatus
	}
	tests := []struct {
		name string
		args args
		want *Transaction
	}{
		{
			name: "create new transaction",
			args: args{
				id:        "trans001",
				walletID:  "wallet002",
				accountID: "a0001",
				amount:    10000,
				currency:  "VND",
				kind:      TransactionIn,
				note:      "test",
				status:    TransactionStatusNew,
			},
			want: &Transaction{
				ID:              "trans001",
				WalletID:        "wallet002",
				AccountID:       "a0001",
				Amount:          10000,
				Currency:        "VND",
				TransactionKind: TransactionIn,
				Note:            "test",
				Status:          TransactionStatusNew,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransaction(tt.args.id, tt.args.walletID, tt.args.accountID, tt.args.amount, tt.args.currency, tt.args.kind, tt.args.note, tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
