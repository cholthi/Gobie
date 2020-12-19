package model

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

type tnxType string

var TransactionNotFoundError error = errors.New("Query returned No trsansactions")

const (
	DEPOSIT tnxType = "deposit"
	TOPUP   tnxType = "topup"
)

type TnxStatus int

const (
	TNX_FAILED  TnxStatus = 4
	TNX_PENDING TnxStatus = 2
	TNX_SUCCESS TnxStatus = 0
)

type Transaction struct {
	ID        uuid.UUID `json:"transaction_id" db:"id"`
	Ref       string    `json:"reference" db:"ref"`
	Currency  string    `json:"currency" db:"currency"`
	Type      tnxType   `json:"transaction_type" db:"type"` //type can be topup or deposit
	Recipient string    `json:"recipient" db:"recipient"`
	Amount    float64   `json:"amount" db:"amount"`
	Status    TnxStatus `json:"status" db:"status"`
	Account   Account   `belongs_to:"account"`
	AccountID uuid.UUID `json:"-" db:"account_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewTransaction(recipient string, amount float64, accountid uuid.UUID, ref string) (*Transaction, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	obj := &Transaction{
		ID:        id,
		Amount:    amount,
		Recipient: recipient,
		AccountID: accountid,
		Ref:       ref,
	}
	return obj, nil
}

func (t *Transaction) Update(includedColumns ...string) error {
	return DB.UpdateColumns(t, includedColumns...)
}

func (t *Transaction) Save(excludedColumns ...string) error {
	return DB.Create(t, excludedColumns...)
}
func (t *Transaction) SetStatus(status TnxStatus) error {
	t.Status = status
	return DB.UpdateColumns(t, "status")
}

func findTransaction(query string, params ...interface{}) ([]Transaction, error) {
	tnxslice := []Transaction{}
	if err := DB.Q().Where(query, params...).All(tnxslice); err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows {
			return nil, TransactionNotFoundError
		}
		return nil, err
	}

	return tnxslice, nil
}

func FindTransactionByID(id uuid.UUID) (*Transaction, error) {
	tnx := &Transaction{}
	err := DB.Q().Find(tnx, id)
	if errors.Unwrap(err) == sql.ErrNoRows {
		return nil, TransactionNotFoundError
	}

	return tnx, err
}
