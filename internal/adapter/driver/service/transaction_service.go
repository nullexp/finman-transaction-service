package service

import (
	"context"
	"log"
	"time"

	"github.com/nullexp/finman-transaction-service/internal/domain"
	domainModel "github.com/nullexp/finman-transaction-service/internal/domain/model"
	"github.com/nullexp/finman-transaction-service/internal/port/driven"
	"github.com/nullexp/finman-transaction-service/internal/port/driven/db"
	"github.com/nullexp/finman-transaction-service/internal/port/driven/db/repository"
	"github.com/nullexp/finman-transaction-service/internal/port/model"
)

func ToDomainTransaction(m model.Transaction) domainModel.Transaction {
	return domainModel.Transaction{
		Id:          m.Id,
		UserId:      m.UserId,
		Type:        m.Type,
		Amount:      m.Amount,
		Date:        m.Date,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func ToModelTransaction(dm domainModel.Transaction) model.Transaction {
	return model.Transaction{
		Id:          dm.Id,
		UserId:      dm.UserId,
		Type:        dm.Type,
		Amount:      dm.Amount,
		Date:        dm.Date,
		Description: dm.Description,
		CreatedAt:   dm.CreatedAt,
		UpdatedAt:   dm.UpdatedAt,
	}
}

func ToDomainTransactions(ms []model.Transaction) []domainModel.Transaction {
	var domains []domainModel.Transaction
	for _, m := range ms {
		domains = append(domains, ToDomainTransaction(m))
	}
	return domains
}

func ToModelTransactions(dms []domainModel.Transaction) []model.Transaction {
	var models []model.Transaction
	for _, dm := range dms {
		models = append(models, ToModelTransaction(dm))
	}
	return models
}

type transactionService struct {
	transactionRepositoryFactory repository.TransactionRepositoryFactory
	dbTransactionFactory         db.DbTransactionFactory
	notificationService          driven.NotificationService
}

func NewTransactionService(trf repository.TransactionRepositoryFactory, dtf db.DbTransactionFactory, ns driven.NotificationService) *transactionService {
	return &transactionService{transactionRepositoryFactory: trf, dbTransactionFactory: dtf, notificationService: ns}
}

func (ts *transactionService) CreateTransaction(ctx context.Context, request model.CreateTransactionRequest) (*model.CreateTransactionResponse, error) {
	if err := request.Validate(ctx); err != nil {
		return nil, err
	}

	tx := ts.dbTransactionFactory.NewTransaction()
	handler, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repository := ts.transactionRepositoryFactory.New(handler)

	// Check balance for withdraw transactions
	if request.Type == "withdrawal" {
		balance, err := repository.GetBalanceByUserId(ctx, request.UserId)
		if err != nil {
			return nil, err
		}
		if balance < request.Amount {
			return nil, domain.ErrInsufficientBalance
		}
	}

	transaction := domainModel.Transaction{
		UserId:      request.UserId,
		Type:        request.Type,
		Amount:      request.Amount,
		Date:        time.Now(),
		Description: request.Description,
	}

	id, err := repository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Send notification asynchronously
	go func(ctx context.Context) {
		message := "Transaction created with ID: " + id
		err := ts.notificationService.SendTransactionNotification(ctx, request.UserId, message)
		if err != nil {
			log.Printf("Failed to send notification: %v\n", err)
		}
	}(ctx)

	return &model.CreateTransactionResponse{Id: id}, nil
}

func (ts *transactionService) GetTransactionById(ctx context.Context, request model.GetTransactionByIdRequest) (*model.GetTransactionByIdResponse, error) {
	if err := request.Validate(ctx); err != nil {
		return nil, err
	}

	tx := ts.dbTransactionFactory.NewTransaction()
	handler, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repository := ts.transactionRepositoryFactory.New(handler)

	transaction, err := repository.GetTransactionById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, domain.ErrTransactionNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &model.GetTransactionByIdResponse{Transaction: ToModelTransaction(*transaction)}, nil
}

func (ts *transactionService) GetOwnTransactionById(ctx context.Context, request model.GetOwnTransactionByIdRequest) (*model.GetOwnTransactionByIdResponse, error) {
	if err := request.Validate(ctx); err != nil {
		return nil, err
	}

	tx := ts.dbTransactionFactory.NewTransaction()
	handler, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repository := ts.transactionRepositoryFactory.New(handler)

	transaction, err := repository.GetTransactionById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, domain.ErrTransactionNotFound
	}

	if transaction.UserId != request.UserId {
		return nil, domain.ErrTransactionNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &model.GetOwnTransactionByIdResponse{Transaction: ToModelTransaction(*transaction)}, nil
}

func (ts *transactionService) GetAllTransactions(ctx context.Context) (*model.GetAllTransactionsResponse, error) {
	tx := ts.dbTransactionFactory.NewTransaction()
	handler, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repository := ts.transactionRepositoryFactory.New(handler)

	transactions, err := repository.GetAllTransactions(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &model.GetAllTransactionsResponse{Transactions: ToModelTransactions(transactions)}, nil
}

func (ts *transactionService) UpdateTransaction(ctx context.Context, request model.UpdateTransactionRequest) error {
	if err := request.Validate(ctx); err != nil {
		return err
	}

	tx := ts.dbTransactionFactory.NewTransaction()
	handler, err := tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repository := ts.transactionRepositoryFactory.New(handler)

	transaction := domainModel.Transaction{
		Id:          request.Id,
		UserId:      request.UserId,
		Type:        request.Type,
		Amount:      request.Amount,
		Description: request.Description,
	}

	err = repository.UpdateTransaction(ctx, transaction)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (ts *transactionService) DeleteTransaction(ctx context.Context, request model.DeleteTransactionRequest) error {
	tx := ts.dbTransactionFactory.NewTransaction()
	handler, err := tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repository := ts.transactionRepositoryFactory.New(handler)

	err = repository.DeleteTransaction(ctx, request.Id)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (ts *transactionService) GetTransactionsByUserId(ctx context.Context, request model.GetTransactionsByUserIdRequest) (*model.GetTransactionsByUserIdResponse, error) {
	if err := request.Validate(ctx); err != nil {
		return nil, err
	}

	tx := ts.dbTransactionFactory.NewTransaction()
	handler, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repository := ts.transactionRepositoryFactory.New(handler)

	transactions, err := repository.GetTransactionsByUserId(ctx, request.UserId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &model.GetTransactionsByUserIdResponse{Transactions: ToModelTransactions(transactions)}, nil
}

func (ts *transactionService) GetTransactionsWithPagination(ctx context.Context, request model.GetTransactionsWithPaginationRequest) (*model.GetTransactionsWithPaginationResponse, error) {
	if err := request.Validate(ctx); err != nil {
		return nil, err
	}

	tx := ts.dbTransactionFactory.NewTransaction()
	handler, err := tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.RollbackUnlessCommitted(ctx)

	repository := ts.transactionRepositoryFactory.New(handler)

	transactions, err := repository.GetTransactionsWithPagination(ctx, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &model.GetTransactionsWithPaginationResponse{Transactions: ToModelTransactions(transactions)}, nil
}
