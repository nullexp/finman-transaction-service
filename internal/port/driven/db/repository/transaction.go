package repository

import (
	"context"

	"github.com/nullexp/finman-transaction-service/internal/domain/model"
	"github.com/nullexp/finman-transaction-service/internal/port/driven/db"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction model.Transaction) (string, error)
	GetTransactionById(ctx context.Context, id string) (*model.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]model.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction model.Transaction) error
	DeleteTransaction(ctx context.Context, id string) error
	GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error)
	GetTransactionsWithPagination(ctx context.Context, offset, limit int) ([]model.Transaction, error)
	GetBalanceByUserId(ctx context.Context, userId string) (int64, error)
}

type TransactionRepositoryFactory interface {
	New(handler db.DbHandler) TransactionRepository
}
