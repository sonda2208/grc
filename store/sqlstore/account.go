package sqlstore

import (
	"github.com/upper/db/v4"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sonda2208/grc/model"
	"github.com/sonda2208/grc/store"
)

const (
	accountTable = "accounts"
)

type SQLAccountStore struct {
	sqlStore *SQLStore
}

func NewSQLAccountStore(s *SQLStore) store.AccountStore {
	return &SQLAccountStore{
		sqlStore: s,
	}
}

func (s SQLAccountStore) Save(a *model.Account) (*model.Account, error) {
	a.PreSave()
	err := a.IsValid()
	if err != nil {
		return nil, err
	}

	iter := s.sqlStore.dbx.SQL().InsertInto(accountTable).Values(*a).Returning("id", "created_at", "updated_at").Iterator()
	defer iter.Close()

	err = iter.NextScan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, model.WithError("store.account.save", err)
	}

	return a, nil
}

func (s SQLAccountStore) Update(a *model.Account) error {
	a.PreUpdate()
	err := a.IsValid()
	if err != nil {
		return err
	}

	_, err = s.sqlStore.dbx.SQL().Update(accountTable).Set(a).Where("id = ?", a.ID).Exec()
	if err != nil {
		return model.WithError("store.account.update", err)
	}

	return nil
}

func (s SQLAccountStore) Get(id int) (*model.Account, error) {
	a := &model.Account{}
	err := s.sqlStore.dbx.SQL().
		SelectFrom(accountTable).
		Where("id = ? and deleted_at is null", id).
		One(&a)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, model.NewNotFoundError(strconv.Itoa(id), "account")
		}

		return nil, model.WithError("store.account.get", err)
	}

	return a, nil
}

func (s SQLAccountStore) GetByEmail(email string) (*model.Account, error) {
	email = model.NormalizeEmail(email)
	a := &model.Account{}
	err := s.sqlStore.dbx.SQL().
		SelectFrom(accountTable).
		Where("email = ? and deleted_at is null", email).
		One(&a)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, model.NewNotFoundError(email, "account")
		}

		return nil, model.WithError("store.account.get_by_email", err)
	}

	return a, nil
}
