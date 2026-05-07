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
	GetUserStats(userID int, start, end string) (TransactionStats, error)
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

type TransactionStats struct {
	DateRange     string
	MoneyIn       int64
	MoneyOut      int64
	MoneyLeftOver int64
}

type TransactionModel struct {
	DB *sql.DB
}

func (m *TransactionModel) Update(userID int, id int, title string, isIncome bool, amount int64, category string, description string, transactionDate time.Time) (int, error) {
	stmt := `
        UPDATE transactions 
        SET 
            title = ?, 
            isIncome = ?, 
            amountInCents = IF(? = TRUE, ABS(?), -ABS(?)), 
            category = ?, 
            description = ?, 
            transactionDate = ?, 
            updated = UTC_TIMESTAMP()
        WHERE id = ? AND userid = ?`

	result, err := m.DB.Exec(stmt, title, isIncome, isIncome, amount, amount, category, description, transactionDate, id, userID)
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

	stmt := `
	INSERT INTO transactions (
    title, 
    isIncome, 
    amountInCents, 
    category, 
    description, 
    transactionDate, 
    created, 
    updated, 
    userid
	)
	VALUES(
		?, 
		?, 
		IF(? = TRUE, ABS(?), -ABS(?)), -- Ensures positive if income, negative if expense
		?, 
		?, 
		?, 
		UTC_TIMESTAMP(), 
		UTC_TIMESTAMP(), 
		?
	)`

	result, err := m.DB.Exec(stmt, title, isIncome, isIncome, amount, category, description, transactionDate, userID)
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

func (m *TransactionModel) GetUserStats(userID int, start, end string) (TransactionStats, error) {
	var stats TransactionStats
	stats.DateRange = start + " to " + end

	/*
			// user selectable date range
			// We use COALESCE to ensure that if no rows are found, we get 0 instead of a NULL error.
			// The query calculates:
			// 1. The range string
			// 2. Sum of positive values (Income)
			// 3. Sum of negative values (Expenses) - we use ABS to show this as a positive "MoneyOut" total
			// 4. The total net sum (LeftOver)
			query := `
		        SELECT
		            CONCAT(?, ' to ', ?),
		            COALESCE(SUM(CASE WHEN amountInCents > 0 THEN amountInCents ELSE 0 END), 0),
		            COALESCE(SUM(CASE WHEN amountInCents < 0 THEN ABS(amountInCents) ELSE 0 END), 0),
		            COALESCE(SUM(amountInCents), 0)
		        FROM transactions
		        WHERE userid = ?
		          AND transactionDate >= ?
		          AND transactionDate <= ?`

				  err := m.DB.QueryRow(query, start, end, userID, start, end).Scan(
					&stats.DateRange,
					&stats.MoneyIn,
					&stats.MoneyOut,
					&stats.MoneyLeftOver,
				)*/

	// Data range is current month
	query := `
			SELECT
				CONCAT(DATE_FORMAT(CURDATE(), '%Y-%m-01'), ' to ', LAST_DAY(CURDATE())),
				COALESCE(SUM(CASE WHEN amountInCents > 0 THEN amountInCents ELSE 0 END), 0),
            	COALESCE(SUM(CASE WHEN amountInCents < 0 THEN ABS(amountInCents) ELSE 0 END), 0),
            	COALESCE(SUM(amountInCents), 0)
			FROM transactions
			WHERE userid = ?
			  AND transactionDate >= DATE_FORMAT(CURDATE(), '%Y-%m-01')
			  AND transactionDate <= LAST_DAY(CURDATE())`

	err := m.DB.QueryRow(query, userID).Scan(
		&stats.DateRange,
		&stats.MoneyIn,
		&stats.MoneyOut,
		&stats.MoneyLeftOver,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return TransactionStats{DateRange: stats.DateRange}, nil
		}
		return stats, err
	}

	return stats, nil
}
