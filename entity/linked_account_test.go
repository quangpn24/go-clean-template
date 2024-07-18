package entity

import (
	"reflect"
	"testing"
)

func TestNewLinkedAccount(t *testing.T) {
	type args struct {
		id          string
		userId      string
		accountName string
	}
	tests := []struct {
		name string
		args args
		want *LinkedAccount
	}{
		{
			name: "create new linked account",
			args: args{
				id:          "1",
				userId:      "u001",
				accountName: "momo",
			},
			want: &LinkedAccount{
				ID:          "1",
				UserID:      "u001",
				AccountName: "momo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLinkedAccount(tt.args.id, tt.args.userId, tt.args.accountName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
