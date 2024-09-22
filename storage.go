package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type AccountStorage interface {
	CreateAccount(account *Account) error
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error
	GetAccountByID(id int) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStorage, error) {
	conStr := "user=postgre dbname=banl password=tedst sslmode=disable"
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
	return nil
}

func (s *PostgresStorage) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStorage) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStorage) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
