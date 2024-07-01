package db

import (
	"context"
	"database/sql"
	"log"
)

type MockDbHandler struct{}

func (m *MockDbHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	log.Printf("MockDbHandler: QueryContext called with query: %s, args: %v", query, args)
	// Simulate returning nil rows and no error
	return nil, nil
}

func (m *MockDbHandler) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	log.Printf("MockDbHandler: QueryRowContext called with query: %s, args: %v", query, args)
	// Simulate returning nil row
	return nil
}

func (m *MockDbHandler) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	log.Printf("MockDbHandler: ExecContext called with query: %s, args: %v", query, args)
	// Simulate returning nil result and no error
	return nil, nil
}
