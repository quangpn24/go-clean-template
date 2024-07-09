package entity

import (
	"reflect"
	"testing"
)

func TestNewAccount(t *testing.T) {
	type args struct {
		id            string
		userId        string
		accountNumber string
	}
	tests := []struct {
		name string
		args args
		want *Account
	}{
		{
			name: "create new account",
			args: args{
				id:            "1",
				userId:        "u001",
				accountNumber: "123456",
			},
			want: &Account{
				ID:            "1",
				UserID:        "u001",
				AccountNumber: "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAccount(tt.args.id, tt.args.userId, tt.args.accountNumber); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
