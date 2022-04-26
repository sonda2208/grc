package app

import (
	"github.com/sonda2208/grc/jobs"
	"github.com/sonda2208/grc/model"
	"github.com/sonda2208/grc/store"
	"github.com/sonda2208/grc/store/sqlstore"
)

type App struct {
	config *model.Config

	Store     store.Store
	JobServer *jobs.JobServer
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
		config:    conf,
		Store:     s,
		JobServer: js,
	}
	return app, nil
}
