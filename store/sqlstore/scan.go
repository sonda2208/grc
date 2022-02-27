package sqlstore

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/upper/db/v4"

	"github.com/jmoiron/sqlx"
	"github.com/sonda2208/guardrails-challenge/model"
	"github.com/sonda2208/guardrails-challenge/store"
)

const (
	scanTable = "scans"
)

type SQLScanStore struct {
	sqlStore *SQLStore
}

func NewSQLScanStore(s *SQLStore) store.ScanStore {
	return &SQLScanStore{
		sqlStore: s,
	}
}

func (s SQLScanStore) Save(sc *model.Scan) (*model.Scan, error) {
	err := sc.IsValid()
	if err != nil {
		return nil, err
	}

	iter := s.sqlStore.dbx.SQL().
		InsertInto(scanTable).
		Values(*sc).
		Returning("id", "created_at").
		Iterator()
	defer iter.Close()

	err = iter.NextScan(&sc.ID, &sc.CreatedAt)
	if err != nil {
		return nil, model.WithError("store.scan.save", err)
	}

	return sc, nil
}

func (s SQLScanStore) Update(sc *model.Scan) error {
	err := sc.IsValid()
	if err != nil {
		return err
	}

	_, err = s.sqlStore.dbx.SQL().Update(scanTable).Set(sc).Where("id = ?", sc.ID).Exec()
	if err != nil {
		return model.WithError("store.scan.update", err)
	}

	return nil
}

func (s SQLScanStore) Get(id int) (*model.Scan, error) {
	sc := &model.Scan{}
	err := s.sqlStore.dbx.SQL().SelectFrom(scanTable).Where("id = ?", id).One(&sc)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, model.NewNotFoundError(strconv.Itoa(id), "scan")
		}

		return nil, model.WithError("store.scan.get", err)
	}

	return sc, nil
}

func (s SQLScanStore) GetByRepository(repoID int, opt *model.ListScansOption) ([]*model.Scan, int, error) {
	orderByClause := ``
	searchCondition := ``
	limitClause := ``
	params := map[string]interface{}{}
	query := `
		select scans.*, count(*) over() as total_count
		from scans
		where
			repo_id = :repoID
			SEARCH_CONDITION
		order by ORDER_BY_CLAUSE
		LIMIT_CLAUSE
	`

	params["repoID"] = repoID

	if opt.PerPage > 0 && opt.Page > 0 {
		limitClause += `
		offset :offset
		limit :limit
		`
		params["limit"] = opt.PerPage
		params["offset"] = (opt.Page - 1) * opt.PerPage
	}

	if opt.Branch != "" {
		searchCondition += " and branch = :branch"
		params["branch"] = opt.Branch
	}

	if opt.Status != "" {
		searchCondition += " and status = :status"
		params["status"] = opt.Status
	}

	if opt.OrderBy != "" {
		cmd := " asc"
		if opt.IsDesc {
			cmd = " desc"
		}

		field := model.ScanSortableFields[opt.OrderBy]
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
		return nil, 0, model.WithError("store.scan.get_by_repo.execute_query", err)
	}

	var res []*struct {
		*model.Scan
		TotalCount int `db:"total_count"`
	}

	err = s.sqlStore.dbx.SQL().NewIterator(rows).All(&res)
	if err != nil {
		return nil, 0, model.WithError("store.scan.get_by_repo.scan", err)
	}

	totalCount := 0
	scans := make([]*model.Scan, len(res))
	for i, r := range res {
		scans[i] = r.Scan
		totalCount = r.TotalCount
	}

	return scans, totalCount, nil
}

func (s SQLScanStore) GetByStatus(status string) ([]*model.Scan, error) {
	var scans []*model.Scan
	err := s.sqlStore.dbx.SQL().
		SelectFrom(scanTable).
		Where("status = ?", status).
		OrderBy("id desc").
		All(&scans)
	if err != nil {
		return nil, model.WithError("store.scan.get_by_status", err)
	}

	return scans, nil
}
