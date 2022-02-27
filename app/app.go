package app

import (
	"github.com/sonda2208/guardrails-challenge/jobs"
	"github.com/sonda2208/guardrails-challenge/model"
	"github.com/sonda2208/guardrails-challenge/store"
	"github.com/sonda2208/guardrails-challenge/store/sqlstore"
)

type App struct {
	config *model.Config

	Store store.Store
	js    *jobs.JobServer
}

func (a App) Config() model.Config {
	return *a.config
}

func New(conf *model.Config) (*App, error) {
	s, err := sqlstore.NewSQLStore(conf.SQLSetting)
	if err != nil {
		return nil, err
	}

	err = s.MigrationDatabaseSchema()
	if err != nil {
		return nil, err
	}

	js, err := jobs.NewJobServer(s)
	if err != nil {
		return nil, err
	}

	js.Start()

	app := &App{
		config: conf,
		Store:  s,
		js:     js,
	}
	return app, nil
}
