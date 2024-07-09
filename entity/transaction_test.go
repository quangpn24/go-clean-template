package entity

import (
	"reflect"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	type args struct {
		id        string
		sender    string
		receiver  string
		accountID string
		amount    float64
		currency  string
		category  TransactionCategory
		note      string
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
				sender:    "",
				receiver:  "wallet002",
				accountID: "a0001",
				amount:    10000,
				currency:  "VND",
				category:  CategoryDeposit,
				note:      "test",
			},
			want: &Transaction{
				ID:               "trans001",
				SenderWalletID:   "",
				ReceiverWalletID: "wallet002",
				AccountID:        "a0001",
				Amount:           10000,
				Currency:         "VND",
				Category:         CategoryDeposit,
				Note:             "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransaction(tt.args.id, tt.args.sender, tt.args.receiver, tt.args.accountID, tt.args.amount, tt.args.currency, tt.args.category, tt.args.note); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
