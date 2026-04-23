package models

import (
	"database/sql"
	"errors"
	"time"
)

type TransactionModelInterface interface {
	Insert(title string, isIncome bool, amount int64, category string, description string, transactionDate time.Time) (int, error)
	Get(id int) (Transaction, error)
	Latest() ([]Transaction, error)
}

type Transaction struct {
	ID              int
	Title           string
	IsIncome        bool
	Amount          int64
	Category        string
	Description     string
	TransactionDate time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TransactionModel struct {
	DB *sql.DB
}

func (m *TransactionModel) Insert(title string, isIncome bool, amount int64, category string, description string, transactionDate time.Time) (int, error) {
	stmt := `INSERT INTO transactions (title, isIncome, amountInCents, category, description, transactionDate, created, updated)
    VALUES(?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, isIncome, amount, category, description, transactionDate)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *TransactionModel) Get(id int) (Transaction, error) {
	stmt := `SELECT id, title, isIncome, amountInCents, category, description, transactionDate, created, updated FROM transactions
    WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	var t Transaction

	err := row.Scan(&t.ID, &t.Title, &t.IsIncome, &t.Amount, &t.Category, &t.Description, &t.TransactionDate, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Transaction{}, ErrNoRecord
		} else {
			return Transaction{}, err
		}
	}
	return t, nil
}

func (m *TransactionModel) Latest() ([]Transaction, error) {
	stmt := `SELECT id, title, isIncome, amountInCents, category, description, transactionDate, created, updated FROM transactions ORDER BY id DESC LIMIT 20`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var t Transaction
		err = rows.Scan(&t.ID, &t.Title, &t.IsIncome, &t.Amount, &t.Category, &t.Description, &t.TransactionDate, &t.CreatedAt, &t.UpdatedAt)

		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
