package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nullexp/finman-transaction-service/internal/domain/model"
	"github.com/nullexp/finman-transaction-service/internal/port/driven/db"
	"github.com/nullexp/finman-transaction-service/internal/port/driven/db/repository"
)

type InMemoryTransactionRepositoryFactory struct {
	repo *InMemoryTransactionRepository
}

func NewInMemoryTransactionRepositoryFactory(repo *InMemoryTransactionRepository) *InMemoryTransactionRepositoryFactory {
	return &InMemoryTransactionRepositoryFactory{
		repo: repo,
	}
}

func (f *InMemoryTransactionRepositoryFactory) New(handler db.DbHandler) repository.TransactionRepository {
	return f.repo
}

// InMemoryTransactionRepository implements TransactionRepository using in-memory storage.
type InMemoryTransactionRepository struct {
	transactions []model.Transaction
	mu           sync.RWMutex
}

// NewInMemoryTransactionRepository creates a new instance of InMemoryTransactionRepository.
func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{
		transactions: make([]model.Transaction, 0),
	}
}

func (r *InMemoryTransactionRepository) CreateTransaction(ctx context.Context, transaction model.Transaction) (string, error) {
	transaction.Id = uuid.New().String() // Generate new UUID for transaction ID
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.transactions = append(r.transactions, transaction)
	return transaction.Id, nil
}

func (r *InMemoryTransactionRepository) GetTransactionById(ctx context.Context, id string) (*model.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, t := range r.transactions {
		if t.Id == id {
			return &t, nil
		}
	}
	return nil, errors.New("transaction not found")
}

func (r *InMemoryTransactionRepository) GetAllTransactions(ctx context.Context) ([]model.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.transactions, nil
}

func (r *InMemoryTransactionRepository) UpdateTransaction(ctx context.Context, transaction model.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, t := range r.transactions {
		if t.Id == transaction.Id {
			transaction.UpdatedAt = time.Now()
			r.transactions[i] = transaction
			return nil
		}
	}
	return errors.New("transaction not found")
}

func (r *InMemoryTransactionRepository) DeleteTransaction(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, t := range r.transactions {
		if t.Id == id {
			r.transactions = append(r.transactions[:i], r.transactions[i+1:]...)
			return nil
		}
	}
	return errors.New("transaction not found")
}

func (r *InMemoryTransactionRepository) GetTransactionsByUserId(ctx context.Context, userId string) ([]model.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userTransactions []model.Transaction
	for _, t := range r.transactions {
		if t.UserId == userId {
			userTransactions = append(userTransactions, t)
		}
	}
	return userTransactions, nil
}

func (r *InMemoryTransactionRepository) GetTransactionsWithPagination(ctx context.Context, offset, limit int) ([]model.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	start := offset
	end := offset + limit
	if start < 0 || start >= len(r.transactions) || end > len(r.transactions) {
		return nil, errors.New("invalid offset or limit")
	}
	return r.transactions[start:end], nil
}

func (r *InMemoryTransactionRepository) GetBalanceByUserId(ctx context.Context, userId string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var depositSum, withdrawalSum int64
	for _, t := range r.transactions {
		if t.UserId == userId {
			if t.Type == "deposit" {
				depositSum += t.Amount
			} else if t.Type == "withdrawal" {
				withdrawalSum += t.Amount
			}
		}
	}
	return depositSum - withdrawalSum, nil
}
