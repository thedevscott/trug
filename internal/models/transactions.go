package models

import (
	"database/sql"
	"errors"
	"time"
)

type TransactionModelInterface interface {
	Insert(userID int, title string, isIncome bool, amount int64, category string, description string, transactionDate time.Time) (int, error)
	Update(userID int, id int, title string, isIncome bool, amount int64, category string, description string, transactionDate time.Time) (int, error)
	Get(userID int, id int) (Transaction, error)
	Delete(userID int, id int) (int, error)
	Latest(userID int) ([]Transaction, error)
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

func (m *TransactionModel) Update(userID int, id int, title string, isIncome bool, amount int64, category string, description string, transactionDate time.Time) (int, error) {
	stmt := `UPDATE transactions 
	SET title = ?, isIncome = ?, amountInCents = ?, category = ?, description = ?, transactionDate = ?, updated = UTC_TIMESTAMP()
	WHERE id = ? and userid = ?`

	result, err := m.DB.Exec(stmt, title, isIncome, amount, category, description, transactionDate, id, userID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (m *TransactionModel) Delete(userID int, id int) (int, error) {
	stmt := `DELETE FROM transactions WHERE id = ? and userid = ?`

	result, err := m.DB.Exec(stmt, id, userID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
func (m *TransactionModel) Insert(userID int, title string, isIncome bool, amount int64, category string, description string, transactionDate time.Time) (int, error) {
	stmt := `INSERT INTO transactions (title, isIncome, amountInCents, category, description, transactionDate, created, updated, userid)
    VALUES(?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, title, isIncome, amount, category, description, transactionDate, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *TransactionModel) Get(userID int, id int) (Transaction, error) {
	stmt := `SELECT id, title, isIncome, amountInCents, category, description, transactionDate, created, updated FROM transactions
    WHERE id = ? and userid = ?`

	row := m.DB.QueryRow(stmt, id, userID)

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

func (m *TransactionModel) Latest(userID int) ([]Transaction, error) {
	stmt := `SELECT id, title, isIncome, amountInCents, category, description, transactionDate, created, updated FROM transactions WHERE userid = ? ORDER BY id`

	rows, err := m.DB.Query(stmt, userID)
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
