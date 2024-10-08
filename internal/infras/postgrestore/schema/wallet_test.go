package schema

import (
	"reflect"
	"testing"
	"time"

	"go-clean-template/internal/entity"
)

func TestWalletSchema_ToWallet(t *testing.T) {
	type fields struct {
		ID         string
		UserID     string
		WalletName string
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *entity.Wallet
	}{
		{
			name: "Test ToWallet",
			fields: fields{
				ID:         "1",
				UserID:     "1",
				WalletName: "My wallet",
			},
			want: &entity.Wallet{
				ID:         "1",
				UserID:     "1",
				WalletName: "My wallet",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WalletSchema{
				ID:         tt.fields.ID,
				UserID:     tt.fields.UserID,
				WalletName: tt.fields.WalletName,
				CreatedAt:  tt.fields.CreatedAt,
				UpdatedAt:  tt.fields.UpdatedAt,
			}
			if got := w.ToWallet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToWallet() = %v, want %v", got, tt.want)
			}
		})
	}
}
