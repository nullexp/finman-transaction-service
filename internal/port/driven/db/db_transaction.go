package db

import (
	"context"
)

type DbTransaction interface {
	Begin(ctx context.Context) (DbHandler, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	RollbackUnlessCommitted(ctx context.Context)
}

type DbTransactionFactory interface {
	NewTransaction() DbTransaction
}
