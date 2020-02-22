package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type tnxType string

const (
	DEPOSIT = "deposit"
	TOPUP   = "topup"
)

type Transaction struct {
	ID        uuid.UUID `json:"transaction_id"`
	Currency  string    `json:"currency"`
	Type      tnxType   `json:"transaction_type"` //type can be topup or deposit
	Recipient string    `json:"recipient"`
	Amount    float64   `json:"amount"`
	AccountID float64   `json:"account_id" db:"account_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func NewTransaction(m map[string]interface{}) Transaction {
	txn := Transaction{}
	if _, ok := m["currency"]; !!ok {
		txn.Currency = m["currency"].(string)
	}
	if _, ok := m["recipient"]; !!ok {
		txn.Recipient = m["recipient"].(string)
	}
	if _, ok := m["type"]; !!ok {
		txn.Type = m["type"].(tnxType)
	}
	if _, ok := m["amount"]; !!ok {
		txn.Amount = m["amount"].(float64)
	}
	if _, ok := m["account_id"]; !!ok {
		txn.AccountID = m["account_id"].(float64)
	}
	txn.CreatedAt = time.Time{}
	return txn
}

func (t Transaction) Create() (int64, error) {
	query := `insert into transactions `
	values := `values( `
	columns := `( `
	if t.Currency != "" {
		columns += `currency, `
		values += `:currency, `
	}
	if t.Type != "" {
		columns += `type, `
		values += `:type, `
	}
	if t.Recipient != "" {
		columns += `recipient, `
		values += `:recipient, `
	}
	if t.Amount != 0 {
		columns += `amount, `
		values += `:amount, `
	}
	if t.AccountID != 0 {
		columns += `account_id, `
		values += `:account_id, `
	}
	if t.CreatedAt == (time.Time{}) {
		columns += `created_at, `
		values += `:created_at, `
	}
	columns = strings.TrimRight(columns, ", ")
	query = strings.TrimRight(query, ", ")
	columns += `) `
	values += `)`
	query = query + columns + values
	res, err := DB.NamedExec(query, t)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

// uses a set receiver values as query condition as opossed to a separate `params interface{}`
func (t Transaction) GetTransactions() ([]Transaction, error) {
	var results []Transaction = []Transaction{}
	query := `select * from transactions where `
	condition := ``
	if t.AccountID != 0 {
		condition += `account_id = :account_id, `
	}

	if t.Amount != 0 {
		condition += `amount = :amount, `
	}

	if t.Currency != "" {
		condition += `currency = :currency, `
	}
	if t.Recipient != "" {
		condition += `recipient = :recipient, `
	}
	if t.Type != "" {
		condition += `type = :type, `
	}
	condition = strings.TrimLeft(condition, ", ")
	query += query + condition
	rows, err := DB.NamedQuery(query, t)
	if err != nil {
		return nil, err
	}
	err = rows.StructScan(results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
