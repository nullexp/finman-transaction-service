package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	adapterDriven "github.com/nullexp/finman-transaction-service/internal/adapter/driven"
	db "github.com/nullexp/finman-transaction-service/internal/adapter/driven/db"
	repository "github.com/nullexp/finman-transaction-service/internal/adapter/driven/db/repository"
	"github.com/nullexp/finman-transaction-service/internal/adapter/driver/service"
	domainModel "github.com/nullexp/finman-transaction-service/internal/domain/model"
	"github.com/nullexp/finman-transaction-service/internal/port/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	repo := repository.NewInMemoryTransactionRepository()
	repoFactory := repository.NewInMemoryTransactionRepositoryFactory(repo)
	txFactory := &db.PostgresTransactionMockFactory{}

	service := service.NewTransactionService(repoFactory, txFactory, adapterDriven.NewMockNotificationService())

	request := model.CreateTransactionRequest{
		UserId:      uuid.New().String(),
		Type:        "deposit",
		Amount:      100,
		Description: "Test transaction",
	}

	ctx := context.Background()
	response, err := service.CreateTransaction(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Id)
}

func TestGetTransactionById(t *testing.T) {
	repo := repository.NewInMemoryTransactionRepository()
	repoFactory := repository.NewInMemoryTransactionRepositoryFactory(repo)
	txFactory := &db.PostgresTransactionMockFactory{}
	service := service.NewTransactionService(repoFactory, txFactory, adapterDriven.NewMockNotificationService())

	transaction := domainModel.Transaction{
		UserId:      uuid.New().String(),
		Type:        "deposit",
		Amount:      100,
		Date:        time.Now(),
		Description: "Test transaction",
	}

	// Create transaction
	ctx := context.Background()
	id, err := repo.CreateTransaction(ctx, transaction)
	assert.NoError(t, err)

	// Get transaction by ID
	request := model.GetTransactionByIdRequest{Id: id}
	response, err := service.GetTransactionById(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, transaction.UserId, response.Transaction.UserId)
	assert.Equal(t, transaction.Type, response.Transaction.Type)
	assert.Equal(t, transaction.Amount, response.Transaction.Amount)
	assert.Equal(t, transaction.Description, response.Transaction.Description)
}

func TestGetOwnTransactionById(t *testing.T) {
	repo := repository.NewInMemoryTransactionRepository()
	repoFactory := repository.NewInMemoryTransactionRepositoryFactory(repo)
	txFactory := &db.PostgresTransactionMockFactory{}
	service := service.NewTransactionService(repoFactory, txFactory, adapterDriven.NewMockNotificationService())

	userId := uuid.New().String()
	transaction := domainModel.Transaction{
		UserId:      userId,
		Type:        "deposit",
		Amount:      100,
		Date:        time.Now(),
		Description: "Test transaction",
	}

	// Create transaction
	ctx := context.Background()
	id, err := repo.CreateTransaction(ctx, transaction)
	assert.NoError(t, err)

	// Get own transaction by ID
	request := model.GetOwnTransactionByIdRequest{Id: id, UserId: userId}
	response, err := service.GetOwnTransactionById(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, transaction.UserId, response.Transaction.UserId)
	assert.Equal(t, transaction.Type, response.Transaction.Type)
	assert.Equal(t, transaction.Amount, response.Transaction.Amount)
	assert.Equal(t, transaction.Description, response.Transaction.Description)
}

func TestGetAllTransactions(t *testing.T) {
	repo := repository.NewInMemoryTransactionRepository()
	repoFactory := repository.NewInMemoryTransactionRepositoryFactory(repo)
	txFactory := &db.PostgresTransactionMockFactory{}
	service := service.NewTransactionService(repoFactory, txFactory, adapterDriven.NewMockNotificationService())

	transaction1 := domainModel.Transaction{
		UserId:      uuid.New().String(),
		Type:        "deposit",
		Amount:      100,
		Date:        time.Now(),
		Description: "Test transaction 1",
	}

	transaction2 := domainModel.Transaction{
		UserId:      uuid.New().String(),
		Type:        "withdrawal",
		Amount:      50,
		Date:        time.Now(),
		Description: "Test transaction 2",
	}

	// Create transactions
	ctx := context.Background()
	_, err := repo.CreateTransaction(ctx, transaction1)
	assert.NoError(t, err)
	_, err = repo.CreateTransaction(ctx, transaction2)
	assert.NoError(t, err)

	// Get all transactions
	response, err := service.GetAllTransactions(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Transactions))
}

func TestUpdateTransaction(t *testing.T) {
	repo := repository.NewInMemoryTransactionRepository()
	repoFactory := repository.NewInMemoryTransactionRepositoryFactory(repo)
	txFactory := &db.PostgresTransactionMockFactory{}
	service := service.NewTransactionService(repoFactory, txFactory, adapterDriven.NewMockNotificationService())

	transaction := domainModel.Transaction{
		UserId:      uuid.New().String(),
		Type:        "deposit",
		Amount:      100,
		Date:        time.Now(),
		Description: "Test transaction",
	}

	// Create transaction
	ctx := context.Background()
	id, err := repo.CreateTransaction(ctx, transaction)
	assert.NoError(t, err)

	// Update transaction
	request := model.UpdateTransactionRequest{
		Id:          id,
		UserId:      uuid.New().String(),
		Type:        "deposit",
		Amount:      200,
		Description: "Updated transaction",
	}

	err = service.UpdateTransaction(ctx, request)
	assert.NoError(t, err)

	// Get transaction by ID to check update
	updatedTransaction, err := repo.GetTransactionById(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, updatedTransaction)
	assert.Equal(t, request.Amount, updatedTransaction.Amount)
	assert.Equal(t, request.Description, updatedTransaction.Description)
}

func TestDeleteTransaction(t *testing.T) {
	repo := repository.NewInMemoryTransactionRepository()
	repoFactory := repository.NewInMemoryTransactionRepositoryFactory(repo)
	txFactory := &db.PostgresTransactionMockFactory{}
	service := service.NewTransactionService(repoFactory, txFactory, adapterDriven.NewMockNotificationService())

	transaction := domainModel.Transaction{
		UserId:      uuid.New().String(),
		Type:        "deposit",
		Amount:      100,
		Date:        time.Now(),
		Description: "Test transaction",
	}

	// Create transaction
	ctx := context.Background()
	id, err := repo.CreateTransaction(ctx, transaction)
	assert.NoError(t, err)

	// Delete transaction
	request := model.DeleteTransactionRequest{Id: id}
	err = service.DeleteTransaction(ctx, request)
	assert.NoError(t, err)

	// Try to get deleted transaction
	deletedTransaction, err := repo.GetTransactionById(ctx, id)
	assert.Error(t, err)
	assert.Nil(t, deletedTransaction)
}

func TestGetTransactionsByUserId(t *testing.T) {
	repo := repository.NewInMemoryTransactionRepository()
	repoFactory := repository.NewInMemoryTransactionRepositoryFactory(repo)
	txFactory := &db.PostgresTransactionMockFactory{}
	service := service.NewTransactionService(repoFactory, txFactory, adapterDriven.NewMockNotificationService())

	userId := uuid.New().String()
	transaction1 := domainModel.Transaction{
		UserId:      userId,
		Type:        "deposit",
		Amount:      100,
		Date:        time.Now(),
		Description: "Test transaction 1",
	}

	transaction2 := domainModel.Transaction{
		UserId:      userId,
		Type:        "withdrawal",
		Amount:      50,
		Date:        time.Now(),
		Description: "Test transaction 2",
	}

	// Create transactions
	ctx := context.Background()
	_, err := repo.CreateTransaction(ctx, transaction1)
	assert.NoError(t, err)
	_, err = repo.CreateTransaction(ctx, transaction2)
	assert.NoError(t, err)

	// Get transactions by user ID
	request := model.GetTransactionsByUserIdRequest{UserId: userId}
	response, err := service.GetTransactionsByUserId(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Transactions))
}

func TestGetTransactionsWithPagination(t *testing.T) {
	repo := repository.NewInMemoryTransactionRepository()
	repoFactory := repository.NewInMemoryTransactionRepositoryFactory(repo)
	txFactory := &db.PostgresTransactionMockFactory{}
	service := service.NewTransactionService(repoFactory, txFactory, adapterDriven.NewMockNotificationService())

	for i := 0; i < 10; i++ {
		transaction := domainModel.Transaction{
			UserId:      uuid.New().String(),
			Type:        "deposit",
			Amount:      100,
			Date:        time.Now(),
			Description: "Test transaction",
		}
		ctx := context.Background()
		_, err := repo.CreateTransaction(ctx, transaction)
		assert.NoError(t, err)
	}

	// Get transactions with pagination
	ctx := context.Background()
	request := model.GetTransactionsWithPaginationRequest{Offset: 0, Limit: 5}
	response, err := service.GetTransactionsWithPagination(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 5, len(response.Transactions))
}
