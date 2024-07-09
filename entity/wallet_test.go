package entity

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewWallet(t *testing.T) {
	type args struct {
		id       string
		userID   string
		balance  float64
		currency string
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
				id:       "0001",
				userID:   "1",
				balance:  10000,
				currency: "VND",
			},
			want: &Wallet{
				ID:       "0001",
				UserID:   "1",
				Balance:  10000,
				Currency: "VND",
			},
			wantErr: nil,
		},
		{
			name: "create wallet failed",
			args: args{
				id:       "",
				userID:   "1",
				balance:  10000,
				currency: "VND",
			},
			want:    nil,
			wantErr: fmt.Errorf("id must not be empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWallet(tt.args.id, tt.args.userID, tt.args.balance, tt.args.currency)

			assert.Equal(t, tt.wantErr, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWallet_Deposit(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		Balance float64
		wantErr error
	}{
		{
			name:    "deposit amount greater than 0",
			amount:  1000,
			Balance: 10000,
			wantErr: nil,
		},
		{
			name:    "deposit amount less than 0",
			amount:  -1000,
			Balance: 10000,
			wantErr: fmt.Errorf("amount must be greater than 0"),
		},
		{
			name:    "deposit amount equal to 0",
			amount:  0,
			Balance: 10000,
			wantErr: fmt.Errorf("amount must be greater than 0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{
				Balance: tt.Balance,
			}

			gotErr := w.Deposit(tt.amount)

			assert.Equal(t, tt.wantErr, gotErr)

			// Check if the balance is updated correctly
			if gotErr == nil {
				assert.Equal(t, tt.Balance+tt.amount, w.Balance)
			}
		})
	}
}

func TestWallet_Withdraw(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		Balance float64
		wantErr error
	}{
		{
			name:    "withdraw success",
			amount:  1000,
			Balance: 10000,
			wantErr: nil,
		},
		{
			name:    "withdraw amount less than 0",
			amount:  -1000,
			Balance: 10000,
			wantErr: fmt.Errorf("amount must be greater than 0"),
		},
		{
			name:    "withdraw amount equal to 0",
			amount:  0,
			Balance: 10000,
			wantErr: fmt.Errorf("amount must be greater than 0"),
		},
		{
			name:    "withdraw amount less than balance",
			amount:  100000,
			Balance: 10000,
			wantErr: fmt.Errorf("insufficient balance"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wallet{
				Balance: tt.Balance,
			}

			gotErr := w.Withdraw(tt.amount)

			assert.Equal(t, tt.wantErr, gotErr)

			// Check if the balance is updated correctly
			if gotErr == nil {
				assert.Equal(t, tt.Balance-tt.amount, w.Balance)
			}
		})
	}
}
