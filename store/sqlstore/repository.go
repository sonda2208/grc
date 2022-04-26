package sqlstore

import (
	"github.com/upper/db/v4"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"
	"github.com/sonda2208/grc/model"
	"github.com/sonda2208/grc/store"
)

const (
	repositoryTable = "repositories"
)

type SQLRepositoryStore struct {
	sqlStore *SQLStore
}

func NewSQLRepositoryStore(s *SQLStore) store.RepositoryStore {
	return &SQLRepositoryStore{
		sqlStore: s,
	}
}

func (s SQLRepositoryStore) Save(r *model.Repository) (*model.Repository, error) {
	err := r.IsValid()
	if err != nil {
		return nil, err
	}

	iter := s.sqlStore.dbx.SQL().
		InsertInto(repositoryTable).
		Values(*r).
		Returning("id", "created_at", "updated_at").
		Iterator()
	defer iter.Close()

	err = iter.NextScan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, model.WithError("store.repository.save", err)
	}

	return r, nil
}

func (s SQLRepositoryStore) Update(r *model.Repository) error {
	r.PreUpdate()
	err := r.IsValid()
	if err != nil {
		return err
	}

	_, err = s.sqlStore.dbx.SQL().Update(repositoryTable).Set(r).Where("id = ?", r.ID).Exec()
	if err != nil {
		return model.WithError("store.repository.update", err)
	}

	return nil
}

func (s SQLRepositoryStore) Get(id int) (*model.Repository, error) {
	r := &model.Repository{}
	err := s.sqlStore.dbx.SQL().
		SelectFrom(repositoryTable).
		Where("id = ? and deleted_at is null", id).
		One(&r)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, model.NewNotFoundError(strconv.Itoa(id), "account")
		}

		return nil, model.WithError("store.repository.get", err)
	}

	return r, nil
}

func (s SQLRepositoryStore) GetByAccount(accountID int, opt *model.ListRepositoriesOption) ([]*model.Repository, int, error) {
	orderByClause := ``
	searchCondition := ``
	limitClause := ``
	params := map[string]interface{}{}
	query := `
		select repositories.*, count(*) over() as total_count
		from repositories
		where
		    account_id = :accountID
			and deleted_at is null
			SEARCH_CONDITION
		order by ORDER_BY_CLAUSE
		LIMIT_CLAUSE
	`

	params["accountID"] = accountID

	if opt.PerPage > 0 && opt.Page > 0 {
		limitClause += `
		offset :offset
		limit :limit
		`
		params["limit"] = opt.PerPage
		params["offset"] = (opt.Page - 1) * opt.PerPage
	}

	if opt.OrderBy != "" {
		cmd := " asc"
		if opt.IsDesc {
			cmd = " desc"
		}

		field := model.RepositorySortableFields[opt.OrderBy]
		orderByClause += field + " " + cmd
	} else {
		orderByClause = "id desc"
	}

	query = strings.Replace(query, "SEARCH_CONDITION", searchCondition, 1)
	query = strings.Replace(query, "ORDER_BY_CLAUSE", orderByClause, 1)
	query = strings.Replace(query, "LIMIT_CLAUSE", limitClause, 1)
	q, args, _ := sqlx.Named(query, params)
	q, args, _ = sqlx.In(q, args...)
	q = sqlx.Rebind(sqlx.DOLLAR, q)

	rows, err := s.sqlStore.dbx.SQL().Query(q, args...)
	if err != nil {
		return nil, 0, model.WithError("store.repository.get_by_account.execute_query", err)
	}

	var res []*struct {
		*model.Repository
		TotalCount int `db:"total_count"`
	}

	err = s.sqlStore.dbx.SQL().NewIterator(rows).All(&res)
	if err != nil {
		return nil, 0, model.WithError("store.repository.get_by_account.scan", err)
	}

	totalCount := 0
	repos := make([]*model.Repository, len(res))
	for i, r := range res {
		repos[i] = r.Repository
		totalCount = r.TotalCount
	}

	return repos, totalCount, nil
}

func (s SQLRepositoryStore) Delete(id int) error {
	deletedAt := time.Now().UTC()
	_, err := s.sqlStore.dbx.SQL().
		Update(repositoryTable).
		Set("deleted_at", deletedAt).
		Where("id = ?", id).
		Exec()
	if err != nil {
		return model.WithError("store.repository.delete", err)
	}

	return nil
}
