package entity

import "fmt"

type User struct {
	ID             string
	FullName       string
	Email          string
	PhoneNumber    string
	CurrentAddress string
}

func NewUser(id string, fullName string, email string, phoneNumber string, currentAddress string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("id must not be empty")
	}
	return &User{
		ID:             id,
		FullName:       fullName,
		Email:          email,
		PhoneNumber:    phoneNumber,
		CurrentAddress: currentAddress,
	}, nil
}
