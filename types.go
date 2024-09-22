package main

import "math/rand"

type Account struct {
	// this struct uses json tags to define the json keys
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Number    int64   `json:"number"`
	Balance   float64 `json:"balance"`
}

func NewAccount(firstName string, lastName string) *Account {
	// This is a simple way to generate an account
	// we use pointers to use the same instance of the account

	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000000)),
	}
}
