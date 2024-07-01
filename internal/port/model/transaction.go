package model

import (
	"context"
	"time"

	validator "github.com/go-playground/validator/v10"
)

type Transaction struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	Type        string    `json:"type"`
	Amount      int64     `json:"amount"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateTransactionRequest struct {
	UserId      string `json:"userId" validate:"required,uuid"`
	Type        string `json:"type" validate:"required,oneof=deposit withdrawal"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Description string `json:"description"`
}

func (dto CreateTransactionRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type CreateTransactionResponse struct {
	Id string `json:"id"`
}

type GetTransactionByIdRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

func (dto GetTransactionByIdRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type GetTransactionByIdResponse struct {
	Transaction Transaction `json:"transaction"`
}

type GetOwnTransactionByIdRequest struct {
	Id     string `json:"id" validate:"required,uuid"`
	UserId string `json:"userId" validate:"required,uuid"`
}

func (dto GetOwnTransactionByIdRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type GetOwnTransactionByIdResponse struct {
	Transaction Transaction `json:"transaction"`
}

type GetAllTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type UpdateTransactionRequest struct {
	Id          string `json:"id" validate:"required,uuid"`
	UserId      string `json:"userId" validate:"required,uuid"`
	Type        string `json:"type" validate:"required,oneof=deposit withdrawal"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Description string `json:"description"`
}

func (dto UpdateTransactionRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type DeleteTransactionRequest struct {
	Id string `json:"id" validate:"required,uuid"`
}

type GetTransactionsByUserIdRequest struct {
	UserId string `json:"userId" validate:"required,uuid"`
}

func (dto GetTransactionsByUserIdRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type GetTransactionsByUserIdResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type GetTransactionsWithPaginationRequest struct {
	Offset int `json:"offset" validate:"gte=0"`
	Limit  int `json:"limit" validate:"gt=0"`
}

func (dto GetTransactionsWithPaginationRequest) Validate(ctx context.Context) error {
	validate := validator.New()
	return validate.StructCtx(ctx, dto)
}

type GetTransactionsWithPaginationResponse struct {
	Transactions []Transaction `json:"transactions"`
}
