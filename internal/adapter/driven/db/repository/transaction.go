package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/nullexp/finman-transaction-service/internal/domain/model"
	"github.com/nullexp/finman-transaction-service/internal/port/driven/db"
	"github.com/nullexp/finman-transaction-service/internal/port/driven/db/repository"
)

type TransactionRepositoryFactory struct{}

func NewTransactionRepositoryFactory() *TransactionRepositoryFactory {
	return &TransactionRepositoryFactory{}
}

func (f *TransactionRepositoryFactory) New(handler db.DbHandler) repository.TransactionRepository {
	return NewTransactionRepository(handler)
}

type TransactionRepository struct {
	handler db.DbHandler
}

func NewTransactionRepository(handler db.DbHandler) *TransactionRepository {
	return &TransactionRepository{handler: handler}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction model.Transaction) (string, error) {
	query := `INSERT INTO transactions (user_id, type, amount, date, description) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id`
	var id string
	err := r.handler.QueryRowContext(ctx, query, transaction.UserId, transaction.Type, transaction.Amount, transaction.Date, transaction.Description).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *TransactionRepository) GetTransactionById(ctx context.Context, id string) (*model.Transaction, error) {
	query := `SELECT id, user_id, type, amount, date, description, created_at, updated_at 
	          FROM transactions 
	          WHERE id = $1`
	row := r.handler.QueryRowContext(ctx, query, id)

	var transaction model.Transaction
	err := row.Scan(&transaction.Id, &transaction.UserId, &transaction.Type, &transaction.Amount, &transaction.Date, &transaction.Description, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) GetAllTransactions(ctx context.Context) ([]model.Transaction, error) {
	query := `SELECT id, user_id, type, amount, date, description, created_at, updated_at 
	          FROM transactions`
	rows, err := r.handler.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.Id, &transaction.UserId, &transaction.Type, &transaction.Amount, &transaction.Date, &transaction.Description, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, rows.Err()
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction model.Transaction) error {
	query := `UPDATE transactions 
	          SET user_id = $1, type = $2, amount = $3, description = $4, updated_at = $5 
	          WHERE id = $6`
	_, err := r.handler.ExecContext(ctx, query, transaction.UserId, transaction.Type, transaction.Amount, transaction.Description, time.Now(), transaction.Id)
	return err
}

func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id string) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := r.handler.ExecContext(ctx, query, id)
	return err
}

func (r *TransactionRepository) GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error) {
	query := `SELECT id, user_id, type, amount, date, description, created_at, updated_at 
	          FROM transactions 
	          WHERE user_id = $1`
	rows, err := r.handler.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.Id, &transaction.UserId, &transaction.Type, &transaction.Amount, &transaction.Date, &transaction.Description, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, rows.Err()
}

func (r *TransactionRepository) GetTransactionsWithPagination(ctx context.Context, offset, limit int) ([]model.Transaction, error) {
	query := `SELECT id, user_id, type, amount, date, description, created_at, updated_at 
	          FROM transactions 
	          ORDER BY date DESC 
	          LIMIT $1 OFFSET $2`
	rows, err := r.handler.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.Id, &transaction.UserId, &transaction.Type, &transaction.Amount, &transaction.Date, &transaction.Description, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, rows.Err()
}

func (r *TransactionRepository) GetBalanceByUserId(ctx context.Context, userId string) (int64, error) {
	var balance int64

	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'deposit' THEN amount ELSE 0 END) - SUM(CASE WHEN type = 'withdrawal' THEN amount ELSE 0 END), 0) 
		FROM 
			transactions 
		WHERE 
			user_id = $1`

	err := r.handler.QueryRowContext(ctx, query, userId).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return balance, nil
}
