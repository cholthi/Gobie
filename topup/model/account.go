package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           int       `json:"account_id"`
	Balance      float64   `json:"airtime_balance"`
	Organisation string    `json:"organisation"`
	CreateAt     time.Time `json:"created_at"`
	UserID       int       `json:"user_id"`
	LastDeposit  uuid.UUID `json:last_deposit`
	LastTopUp    uuid.UUID `json:"last_topup"`
}

func GetAccountByUserID(id int) (Account, error) {
	query := `select * from accounts where accounts.user_id = ?`
	acc := Account{}
	DB.Select(&acc, query, id)
	if acc == (Account{}) {
		return (Account{}), errors.New("Database Error: could not populate Account")
	}
	return acc, nil
}

func AccountCommitTransaction(params map[string]interface{}) error {
	//createdAt := time.Now().Format(time.RFC3339)
	acc, err := GetAccountByUserID(params["user_id"].(int))
	if err != nil {
		return err
	}
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	tnx := NewTransaction(params)
	insertID, err := tnx.Create()
	if err != nil {
		tx.Rollback()
		return err
	}

	accountquery := `update accounts set balance=? ,last_topup=? where user_id = ?`
	top64, _ := params["amount"].(float64)
	bal := acc.Balance - top64 // be carefull!
	userid, ok := params["user_id"].(int)
	if !ok {
		panic(ok)
	}
	res := DB.MustExec(accountquery, bal, insertID, userid)
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func GetAccountIDByUserID(id int64) (*Account, error) {
	acc := &Account{}
	query := `select id from accounts where user_id = ?`
	err := DB.Get(acc, query, id)
	if err != nil {
		return nil, err
	}
	return acc, nil

}
