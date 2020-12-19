package model

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/netlify/gotrue/models"
	"github.com/netlify/gotrue/storage"
	"github.com/pkg/errors"
)

type Transactions []Transaction
type Status int

var AccountNotFound error = errors.New("Account not found")

const (
	Active  Status = 1
	Disable Status = 0
)

type Account struct {
	ID           uuid.UUID    `json:"account_id" db:"id"`
	Balance      float64      `json:"airtime_balance" db:"balance"`
	Organization string       `json:"organization" db:"organization"`
	User         *models.User `belongs_to:"user"`
	Status       Status       `json:"-" db:"status"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	UserID       uuid.UUID    `json:"-" db:"user_id"`
	LastDeposit  uuid.UUID    `json:"last_deposit" db:"last_deposit"`
	LastTopUp    uuid.UUID    `json:"last_topup" db:"last_topup"`
}

func NewAccount(organization string, userid uuid.UUID) (*Account, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	obj := &Account{
		ID:           id,
		Organization: organization,
		UserID:       userid,
		Balance:      0,
	}
	return obj, nil
}

func (*Account) TableName() string {
	return "accounts"
}

func (a *Account) SetStatus(s Status) {
	a.Status = s
}
func (a *Account) Save(excludedColumns ...string) error {
	return DB.Create(a, excludedColumns...)
}

func (a *Account) SetBalance(bal float64) {
	a.Balance = bal
}

func (a *Account) Update(includedColumns ...string) error {
	return DB.UpdateColumns(a, includedColumns...)
}

func findAccount(query string, params ...interface{}) (*Account, error) {
	obj := &Account{}
	if err := DB.Q().Where(query, params...).First(obj); err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows {
			return nil, AccountNotFound
		}
		return nil, errors.Wrap(err, "error finding account")
	}

	return obj, nil

}

func FindAccountByUserID(userid uuid.UUID) (*Account, error) {
	return findAccount("user_id = ?", userid)
}

func FindAccountByOrganization(org string) (*Account, error) {
	return findAccount("organization = ?", org)
}

func (a *Account) DoTransaction(tnx *Transaction) error {
	err := DB.Transaction(func(tx *pop.Connection) error {
		newbalance := a.Balance - tnx.Amount
		a.SetBalance(newbalance)
		err := a.Update()
		if err != nil {
			return err
		}
		tnx.AccountID = a.ID
		err = tnx.Save()
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func CreateAccount(email string, organization string) (*Account, error) {
	instance_id, err := uuid.FromString("00000000-0000-0000-0000-000000000000")
	if err != nil {
		return nil, err
	}

	tx := &storage.Connection{DB}
	user, err := models.FindUserByEmailAndAudience(tx, instance_id, email, "jedco")
	if err != nil {
		return nil, err
	}

	acct, err := NewAccount(organization, user.ID)
	if err != nil {
		return nil, err
	}

	acct.SetBalance(float64(0.0))
	acct.SetStatus(Active)
	return acct, nil
}

func FindUserByID(userid uuid.UUID) (*models.User, error) {
	tx := &storage.Connection{DB}
	return models.FindUserByID(tx, userid)
}
