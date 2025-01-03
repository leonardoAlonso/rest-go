package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type AccountStorage interface {
	CreateAccount(account *Account) error
	UpdateAccount(account *Account) error
	DeleteAccount(id int) error
	GetAccountByID(id int) (*Account, error)
	GetAccountByNumber(number int64) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStorage, error) {
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	conStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		db_user,
		db_name,
		db_password,
		db_host,
		db_port,
	)

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
	if err := s.CreateAccountTable(); err != nil {
		return err
	}
	return s.addPasswordColumn()
}

func (s *PostgresStorage) addPasswordColumn() error {
	query := `ALTER TABLE accounts ADD COLUMN IF NOT EXISTS encripted_password varchar(255)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateAccountTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS accounts (
      id SERIAL PRIMARY KEY,
      first_name varchar(50),
      last_name varchar(50),
      encripted_password varchar(255),
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
      INSERT INTO accounts (first_name, last_name, number, balance, encripted_password)
      VALUES ($1, $2, $3, $4, $5)
  `
	_, err := s.db.Query(
		query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.EncriptedPassword,
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

func (s *PostgresStorage) GetAccountByNumber(number int64) (*Account, error) {
	rows, err := s.db.Query(
		"SELECT id, first_name, last_name, encripted_password, number, balance, created_at FROM accounts where number = $1",
		number,
	)
	if err != nil {
		ErrorLogger.Println("Error getting account by number: ", err)
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	ErrorLogger.Println("account not found")
	return nil, fmt.Errorf("account %d not found", number)
}

func (s *PostgresStorage) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(
		"SELECT id, first_name, last_name, encripted_password, number, balance, created_at FROM accounts where id = $1",
		id,
	)
	if err != nil {
		ErrorLogger.Println("Error getting account by id: ", err)
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	fmt.Println("account not found")
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	query := `SELECT id, first_name, last_name, encripted_password, number, balance, created_at FROM accounts`

	rows, err := s.db.Query(query)
	if err != nil {
		ErrorLogger.Println("Error getting accounts: ", err)
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
		&account.EncriptedPassword,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)
	if err != nil {
		ErrorLogger.Println("Error scanning account: ", err)
	}
	return account, err
}
