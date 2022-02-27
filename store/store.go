package store

import "github.com/sonda2208/guardrails-challenge/model"

type Store interface {
	Account() AccountStore
	Repository() RepositoryStore
	Scan() ScanStore

	MigrationDatabaseSchema() error
}

type AccountStore interface {
	Save(a *model.Account) (*model.Account, error)
	Update(a *model.Account) error
	Get(id int) (*model.Account, error)
	GetByEmail(email string) (*model.Account, error)
}

type RepositoryStore interface {
	Save(r *model.Repository) (*model.Repository, error)
	Update(r *model.Repository) error
	Get(id int) (*model.Repository, error)
	GetByAccount(accountID int, opt *model.ListRepositoriesOption) ([]*model.Repository, int, error)
	Delete(id int) error
}

type ScanStore interface {
	Save(s *model.Scan) (*model.Scan, error)
	Update(s *model.Scan) error
	Get(id int) (*model.Scan, error)
	GetByRepository(repoID int, opt *model.ListScansOption) ([]*model.Scan, int, error)
	GetByStatus(status string) ([]*model.Scan, error)
}
