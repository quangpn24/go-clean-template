package entity

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewWallet(t *testing.T) {
	type args struct {
		id         string
		userID     string
		walletName string
	}
	tests := []struct {
		name    string
		args    args
		want    *Wallet
		wantErr error
	}{
		{
			name: "create wallet success",
			args: args{
				id:         "0001",
				userID:     "1",
				walletName: "quangpn's wallet",
			},
			want: &Wallet{
				ID:         "0001",
				UserID:     "1",
				WalletName: "quangpn's wallet",
			},
			wantErr: nil,
		},
		{
			name: "create wallet failed",
			args: args{
				id:         "",
				userID:     "1",
				walletName: "abc",
			},
			want:    nil,
			wantErr: fmt.Errorf("id must not be empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWallet(tt.args.id, tt.args.userID, tt.args.walletName)

			assert.Equal(t, tt.wantErr, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
