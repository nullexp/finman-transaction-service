package grpc

import (
	"context"
	"time"

	transactionv1 "github.com/nullexp/finman-transaction-service/internal/adapter/driver/grpc/proto/transaction/v1"
	"github.com/nullexp/finman-transaction-service/internal/port/driver"
	"github.com/nullexp/finman-transaction-service/internal/port/model"
)

type TransactionService struct {
	transactionv1.UnimplementedTransactionServiceServer
	service driver.TransactionService
}

func NewTransactionService(us driver.TransactionService) *TransactionService {
	return &TransactionService{service: us}
}

func CastTransactionToProto(tx *model.Transaction) *transactionv1.Transaction {
	return &transactionv1.Transaction{
		Id:          tx.Id,
		UserId:      tx.UserId,
		Type:        tx.Type,
		Amount:      tx.Amount,
		Date:        tx.Date.Format(time.RFC3339), // Assuming protobuf uses RFC3339 string format for dates
		Description: tx.Description,
		CreatedAt:   tx.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   tx.UpdatedAt.Format(time.RFC3339),
	}
}

func CastTransactionsToProtoArray(transactions []model.Transaction) []*transactionv1.Transaction {
	if transactions == nil {
		return nil
	}

	var protoTransactions []*transactionv1.Transaction
	for _, tx := range transactions {
		protoTx := CastTransactionToProto(&tx)
		protoTransactions = append(protoTransactions, protoTx)
	}

	return protoTransactions
}

func (ts TransactionService) CreateTransaction(ctx context.Context, request *transactionv1.CreateTransactionRequest) (*transactionv1.CreateTransactionResponse, error) {
	rs, err := ts.service.CreateTransaction(ctx, model.CreateTransactionRequest{
		UserId:      request.UserId,
		Type:        request.Type,
		Amount:      request.Amount,
		Description: request.Description,
	})
	if err != nil {
		return nil, err
	}

	return &transactionv1.CreateTransactionResponse{Id: rs.Id}, nil
}

func (ts TransactionService) GetTransactionById(ctx context.Context, request *transactionv1.GetTransactionByIdRequest) (*transactionv1.GetTransactionByIdResponse, error) {
	rs, err := ts.service.GetTransactionById(ctx, model.GetTransactionByIdRequest{Id: request.Id})
	if err != nil {
		return nil, err
	}

	return &transactionv1.GetTransactionByIdResponse{Transaction: CastTransactionToProto(&rs.Transaction)}, nil
}

func (ts TransactionService) GetTransactionsByUserId(ctx context.Context, request *transactionv1.GetTransactionsByUserIdRequest) (*transactionv1.GetTransactionsByUserIdResponse, error) {
	rs, err := ts.service.GetTransactionsByUserId(ctx, model.GetTransactionsByUserIdRequest{UserId: request.UserId})
	if err != nil {
		return nil, err
	}

	return &transactionv1.GetTransactionsByUserIdResponse{Transactions: CastTransactionsToProtoArray(rs.Transactions)}, nil
}

func (ts TransactionService) GetOwnTransactionById(ctx context.Context, request *transactionv1.GetOwnTransactionByIdRequest) (*transactionv1.GetOwnTransactionByIdResponse, error) {
	rs, err := ts.service.GetOwnTransactionById(ctx, model.GetOwnTransactionByIdRequest{UserId: request.UserId, Id: request.Id})
	if err != nil {
		return nil, err
	}

	return &transactionv1.GetOwnTransactionByIdResponse{Transaction: CastTransactionToProto(&rs.Transaction)}, nil
}

func (ts TransactionService) GetAllTransactions(ctx context.Context, request *transactionv1.GetAllTransactionsRequest) (*transactionv1.GetAllTransactionsResponse, error) {
	rs, err := ts.service.GetAllTransactions(ctx)
	if err != nil {
		return nil, err
	}

	return &transactionv1.GetAllTransactionsResponse{Transactions: CastTransactionsToProtoArray(rs.Transactions)}, nil
}

func (ts TransactionService) UpdateTransaction(ctx context.Context, request *transactionv1.UpdateTransactionRequest) (*transactionv1.UpdateTransactionResponse, error) {
	err := ts.service.UpdateTransaction(ctx, model.UpdateTransactionRequest{
		Id:          request.Id,
		UserId:      request.UserId,
		Type:        request.Type,
		Amount:      request.Amount,
		Description: request.Description,
	})
	if err != nil {
		return nil, err
	}

	return &transactionv1.UpdateTransactionResponse{}, nil
}

func (ts TransactionService) DeleteTransaction(ctx context.Context, request *transactionv1.DeleteTransactionRequest) (*transactionv1.DeleteTransactionResponse, error) {
	err := ts.service.DeleteTransaction(ctx, model.DeleteTransactionRequest{Id: request.Id})
	if err != nil {
		return nil, err
	}

	return &transactionv1.DeleteTransactionResponse{}, nil
}

func (ts TransactionService) GetTransactionsWithPagination(ctx context.Context, request *transactionv1.GetTransactionsWithPaginationRequest) (*transactionv1.GetTransactionsWithPaginationResponse, error) {
	rs, err := ts.service.GetTransactionsWithPagination(ctx, model.GetTransactionsWithPaginationRequest{Offset: int(request.Offset), Limit: int(request.Limit)})
	if err != nil {
		return nil, err
	}

	return &transactionv1.GetTransactionsWithPaginationResponse{Transactions: CastTransactionsToProtoArray(rs.Transactions)}, nil
}
