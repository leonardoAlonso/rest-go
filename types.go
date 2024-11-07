package main

import (
	"math/rand"
	"time"
)

type TransferRequest struct {
	// this struct is used to define the request body for the transfer endpoint
	ToAccountID int     `json:"to_account_id"`
	Amount      float64 `json:"amount"`
}

type CreateAccountRequest struct {
	// this struct is used to define the request body for the create account endpoint
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	// this struct uses json tags to define the json keys
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    int64     `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName string, lastName string) *Account {
	// This is a simple way to generate an account
	// we use pointers to use the same instance of the account

	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000000)),
		CreatedAt: time.Now().UTC(),
	}
}
