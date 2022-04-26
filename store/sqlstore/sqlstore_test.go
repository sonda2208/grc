package sqlstore_test

import (
	"os"

	"github.com/icrowley/fake"
	"github.com/sonda2208/grc/model"
	"github.com/sonda2208/grc/store/sqlstore"

	"github.com/joho/godotenv"
)

type TestHelper struct {
	Store *sqlstore.SQLStore

	SampleAccount    *model.Account
	SampleRepository *model.Repository
}

func Setup() (*TestHelper, error) {
	godotenv.Load()
	s, err := sqlstore.NewSQLStore(model.SQLSetting{
		DataSource:         os.Getenv("DATA_SOURCE"),
		MaxIdleConnections: 4,
		MaxOpenConnections: 4,
	})
	if err != nil {
		return nil, err
	}

	err = s.MigrationDatabaseSchema()
	if err != nil {
		return nil, err
	}

	th := &TestHelper{
		Store: s,
	}

	return th, nil
}

func (th *TestHelper) InitBasic() error {
	acc, err := th.CreateAccount()
	if err != nil {
		return err
	}

	repo, err := th.CreateRepository(acc.ID)
	if err != nil {
		return err
	}

	th.SampleAccount = acc
	th.SampleRepository = repo
	return nil
}

func (th *TestHelper) CreateAccount() (*model.Account, error) {
	acc := &model.Account{
		Email:     fake.EmailAddress(),
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
	}

	ru, err := th.Store.Account().Save(acc)
	if err != nil {
		return nil, err
	}

	return ru, nil
}

func (th *TestHelper) CreateRepository(accountID int) (*model.Repository, error) {
	repo := &model.Repository{
		AccountID: accountID,
		Name:      fake.ProductName(),
		URL:       "https://github.com/test",
	}

	rr, err := th.Store.Repository().Save(repo)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func (th *TestHelper) CreateScan(repoID int) (*model.Scan, error) {
	s := &model.Scan{
		RepoID:  repoID,
		Branch:  "main",
		Commit:  "d41727b",
		Status:  model.ScanStatusQueued,
		Message: "",
	}

	s, err := th.Store.Scan().Save(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (th *TestHelper) Teardown() {
	th.Store.DropAllRecords()
}
