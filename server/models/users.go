package models

import "time"

type UserInterface interface {
	GetID() int
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	IsAdmin      bool      `json:"is_admin"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Pin          *string   `json:"pin"`
}

func (u *User) GetID() int {
	return u.ID
}
