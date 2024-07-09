package entity

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewUser(t *testing.T) {
	type args struct {
		id             string
		fullName       string
		email          string
		phoneNumber    string
		currentAddress string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr error
	}{
		{
			name: "create new user success",
			args: args{
				id:             "u0001",
				fullName:       "John Doe",
				email:          "john.doe@gmail.com",
				phoneNumber:    "08123456789",
				currentAddress: "HCM",
			},
			want: &User{
				ID:             "u0001",
				FullName:       "John Doe",
				Email:          "john.doe@gmail.com",
				PhoneNumber:    "08123456789",
				CurrentAddress: "HCM",
			},
			wantErr: nil,
		},
		{
			name: "create new user failed",
			args: args{
				id:             "",
				fullName:       "John Doe",
				email:          "john.doe@gmail.com",
				phoneNumber:    "08123456789",
				currentAddress: "HCM",
			},
			want:    nil,
			wantErr: fmt.Errorf("id must not be empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.id, tt.args.fullName, tt.args.email, tt.args.phoneNumber, tt.args.currentAddress)

			assert.Equal(t, tt.wantErr, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
