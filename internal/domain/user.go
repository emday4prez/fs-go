package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string
	Username       string
	HashedPassword string
}

func NewUser(id, username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:             id,
		Username:       username,
		HashedPassword: string(hashedPassword),
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))

	return err == nil
}
