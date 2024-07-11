package schema

import (
	"reflect"
	"testing"
	"time"

	"go-clean-template/entity"
)

func TestAccountSchema_ToAccount(t *testing.T) {
	type fields struct {
		ID            string
		UserID        string
		BankName      string
		AccountNumber string
		IsLinked      bool
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}
	now := time.Now()
	tests := []struct {
		name   string
		fields fields
		want   *entity.Account
	}{
		{
			name: "Test ToAccount",
			fields: fields{
				ID:            "a_001",
				UserID:        "u_001",
				BankName:      "Bank A",
				AccountNumber: "1234567890",
				IsLinked:      true,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			want: &entity.Account{
				ID:            "a_001",
				UserID:        "u_001",
				BankName:      "Bank A",
				AccountNumber: "1234567890",
				IsLinked:      true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccountSchema{
				ID:            tt.fields.ID,
				UserID:        tt.fields.UserID,
				BankName:      tt.fields.BankName,
				AccountNumber: tt.fields.AccountNumber,
				IsLinked:      tt.fields.IsLinked,
				CreatedAt:     tt.fields.CreatedAt,
				UpdatedAt:     tt.fields.UpdatedAt,
			}
			if got := a.ToAccount(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
