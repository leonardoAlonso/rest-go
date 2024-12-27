package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type AccountStorage interface {
	CreateAccount(account *Account) error
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error
	GetAccountByID(id int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStorage, error) {
	conStr := "user=leonardo.alonso dbname=bank password=aea4f8261e sslmode=disable"
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStorage) CreateAccountTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS accounts (
      id SERIAL PRIMARY KEY,
      first_name varchar(50),
      last_name varchar(50),
      number bigint,
      balance double precision,
      created_at timestamp default current_timestamp,
      modified_at timestamp default current_timestamp
    ) 
  `

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateAccount(account *Account) error {
	query := `
      INSERT INTO accounts (first_name, last_name, number, balance)
      VALUES ($1, $2, $3, $4)
  `
	_, err := s.db.Query(
		query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStorage) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM accounts WHERE id = $1", id)
	return err
}

func (s *PostgresStorage) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(
		"SELECT id, first_name, last_name, number, balance, created_at FROM accounts where id = $1",
		id,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	fmt.Println("account not found")
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at FROM accounts`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, error := scanIntoAccount(rows)
		if error != nil {
			return nil, error
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)
	return account, err
}
