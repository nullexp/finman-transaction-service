package domain

import "errors"

var (
	ErrTransactionNotFound = errors.New("ErrTransactionNotFound: Transaction not found")
	ErrInsufficientBalance = errors.New("ErrTransactionNotFound: Insufficient balance")
)
