package schema

import (
	"reflect"
	"testing"
	"time"

	"go-clean-template/internal/entity"
)

func TestAccountSchema_ToLinkedAccount(t *testing.T) {
	type fields struct {
		ID          string
		UserID      string
		AccountName string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	now := time.Now()
	tests := []struct {
		name   string
		fields fields
		want   *entity.LinkedAccount
	}{
		{
			name: "Test ToLinkedAccount",
			fields: fields{
				ID:          "a_001",
				UserID:      "u_001",
				AccountName: "Momo",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			want: &entity.LinkedAccount{
				ID:          "a_001",
				UserID:      "u_001",
				AccountName: "Momo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &LinkedAccountSchema{
				ID:          tt.fields.ID,
				UserID:      tt.fields.UserID,
				AccountName: tt.fields.AccountName,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			if got := a.ToLinkedAccount(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToLinkedAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
