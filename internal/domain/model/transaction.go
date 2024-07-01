package model

import "time"

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
