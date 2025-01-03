package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

type TransferRequest struct {
	// this struct is used to define the request body for the transfer endpoint
	ToAccountID int     `json:"to_account_id"`
	Amount      float64 `json:"amount"`
}

type CreateAccountRequest struct {
	// this struct is used to define the request body for the create account endpoint
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Account struct {
	// this struct uses json tags to define the json keys
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Number            int64     `json:"number"`
	EncriptedPassword string    `json:"-"` // this field will not be serialized
	Balance           float64   `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewAccount(firstName string, lastName string, password string) (*Account, error) {
	// This is a simple way to generate an account
	// we use pointers to use the same instance of the account
	encrpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		EncriptedPassword: string(encrpw),
		Number:            int64(rand.Intn(1000000000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func (a *Account) ComparePassword(password string) bool {
	samePassword := bcrypt.CompareHashAndPassword(
		[]byte(a.EncriptedPassword),
		[]byte(password),
	) == nil
	if !samePassword {
		WarningLogger.Println("Password does not match")
		return false
	}
	return true
}
