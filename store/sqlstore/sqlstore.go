package sqlstore

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/rs/zerolog/log"
	"github.com/sonda2208/guardrails-challenge/model"
	"github.com/sonda2208/guardrails-challenge/store"
	"github.com/sonda2208/guardrails-challenge/util"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"

	_ "github.com/lib/pq"
)

const (
	DBPingAttempts  = 6
	DBPingTimeoutMs = 10000 // 10 seconds
)

type SQLStore struct {
	settings model.SQLSetting
	db       *sqlx.DB
	dbx      db.Session

	account store.AccountStore
	repo    store.RepositoryStore
	scan    store.ScanStore
}

func NewSQLStore(setting model.SQLSetting) (*SQLStore, error) {
	s := &SQLStore{
		settings: setting,
	}

	err := s.initConnection()
	if err != nil {
		return nil, err
	}

	s.account = NewSQLAccountStore(s)
	s.repo = NewSQLRepositoryStore(s)
	s.scan = NewSQLScanStore(s)

	return s, nil
}

func (s *SQLStore) initConnection() error {
	db, err := sqlx.Open("postgres", s.settings.DataSource)
	if err != nil {
		return err
	}

	for i := 0; i < DBPingAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), DBPingTimeoutMs*time.Millisecond)
		defer cancel()

		_, err := db.QueryContext(ctx, "select 1")
		if err != nil {
			if i == DBPingAttempts-1 {
				return err
			} else {
				log.Info().Err(err).Msgf("Failed to ping database server, retrying in 1s")
				time.Sleep(1 * time.Second)
			}
		} else {
			break
		}
	}

	db.SetMaxIdleConns(s.settings.MaxIdleConnections)
	db.SetMaxOpenConns(s.settings.MaxOpenConnections)

	s.db = db
	s.dbx, err = postgresql.New(db.DB)
	if err != nil {
		return err
	}

	return nil
}

func (s SQLStore) Account() store.AccountStore {
	return s.account
}

func (s SQLStore) Repository() store.RepositoryStore {
	return s.repo
}

func (s SQLStore) Scan() store.ScanStore {
	return s.scan
}

func (s SQLStore) MigrationDatabaseSchema() error {
	dir, _ := util.FindDir("migration")
	err := goose.Up(s.db.DB, dir)
	if err != nil {
		return err
	}

	return nil
}

func (s SQLStore) DropAllRecords() {
	s.db.Exec("truncate table " + accountTable + " cascade")
	s.db.Exec("truncate table " + repositoryTable + " cascade")
	s.db.Exec("truncate table " + scanTable + " cascade")
}
