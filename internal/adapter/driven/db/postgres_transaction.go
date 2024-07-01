package db

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/nullexp/finman-transaction-service/internal/port/driven/db"
)

type PostgresDbTransaction struct {
	db        *sql.DB
	tx        *sql.Tx
	committed bool
}

func NewPostgresDbTransaction(db *sql.DB) *PostgresDbTransaction {
	return &PostgresDbTransaction{db: db}
}

func (p *PostgresDbTransaction) Begin(ctx context.Context) (db.DbHandler, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	p.tx = tx
	return p, nil
}

func (p *PostgresDbTransaction) Commit(ctx context.Context) error {
	if p.tx == nil {
		return errors.New("no transaction in progress")
	}
	err := p.tx.Commit()
	if err == nil {
		p.committed = true
	}
	return err
}

func (p *PostgresDbTransaction) Rollback(ctx context.Context) error {
	if p.tx == nil {
		return errors.New("no transaction in progress")
	}
	return p.tx.Rollback()
}

func (p *PostgresDbTransaction) RollbackUnlessCommitted(ctx context.Context) {
	if !p.committed {
		err := p.Rollback(ctx)
		if err != nil {
			log.Print(err)
		}
	}
}

func (p *PostgresDbTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.tx.QueryContext(ctx, query, args...)
}

func (p *PostgresDbTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.tx.ExecContext(ctx, query, args...)
}

func (p *PostgresDbTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.tx.QueryRowContext(ctx, query, args...)
}

type PostgresDbTransactionFactory struct {
	db *sql.DB
}

func NewPostgresDbTransactionFactory(db *sql.DB) *PostgresDbTransactionFactory {
	return &PostgresDbTransactionFactory{db: db}
}

func (f *PostgresDbTransactionFactory) NewTransaction() db.DbTransaction {
	return &PostgresDbTransaction{db: f.db}
}
