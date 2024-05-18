package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostGresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=golang password=root sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    	id SERIAL PRIMARY KEY,
    	first_name VARCHAR(50) NOT NULL,
    	last_name VARCHAR(50) NOT NULL,
    	number SERIAL,
    	balance SERIAL,
    	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
	INSERT INTO account(first_name, last_name, number, balance)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	err := s.db.QueryRow(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance).Scan(&acc.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", acc.ID)
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		acc := new(Account)
		err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Balance, &acc.Number, &acc.CreatedAt)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}
