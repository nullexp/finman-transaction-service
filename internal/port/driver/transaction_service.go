package driver

import (
	"context"

	"github.com/nullexp/finman-transaction-service/internal/port/model"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, request model.CreateTransactionRequest) (*model.CreateTransactionResponse, error)
	GetTransactionById(ctx context.Context, request model.GetTransactionByIdRequest) (*model.GetTransactionByIdResponse, error)
	GetOwnTransactionById(ctx context.Context, request model.GetOwnTransactionByIdRequest) (*model.GetOwnTransactionByIdResponse, error)
	GetAllTransactions(ctx context.Context) (*model.GetAllTransactionsResponse, error)
	UpdateTransaction(ctx context.Context, request model.UpdateTransactionRequest) error
	DeleteTransaction(ctx context.Context, request model.DeleteTransactionRequest) error
	GetTransactionsByUserId(ctx context.Context, request model.GetTransactionsByUserIdRequest) (*model.GetTransactionsByUserIdResponse, error)
	GetTransactionsWithPagination(ctx context.Context, request model.GetTransactionsWithPaginationRequest) (*model.GetTransactionsWithPaginationResponse, error)
}
