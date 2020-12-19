package model

import (
	"time"

	"github.com/gobuffalo/uuid"
)

var MaxRetry int = 3

type RetryTransaction struct {
	ID            uuid.UUID   `json:"-" db:"id"`
	TransactionID uuid.UUID   `json:"-" db:"transaction_id"`
	Transaction   Transaction `json:"transaction" belongs_to:"transaction"`
	Retry         int         `json:"retries" db:"retry"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

func (rt *RetryTransaction) TableName() string {
	return "retry_transactions"
}

func NewRetryTransaction(t *Transaction) (*RetryTransaction, error) {
	obj := &RetryTransaction{}
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	obj.ID = id
	obj.TransactionID = t.ID
	obj.Retry = 0

	return obj, nil
}

func GetTransactionsForRetry() ([]RetryTransaction, error) {
	retry := []RetryTransaction{}
	err := DB.Q().Eager("Transaction").Where("retry < ?", MaxRetry).All(&retry)
	if err != nil {
		return nil, err
	}

	return retry, nil
}

func (rt *RetryTransaction) Save(excludedColumns ...string) error {
	return DB.Create(rt, excludedColumns...)
}

func (rt *RetryTransaction) IncrementRetry() {
	rt.Retry += 1
}
