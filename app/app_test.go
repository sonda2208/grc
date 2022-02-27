package app_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sonda2208/guardrails-challenge/app"
	"github.com/sonda2208/guardrails-challenge/jobs"
	"github.com/sonda2208/guardrails-challenge/store"
	"github.com/sonda2208/guardrails-challenge/store/mockstore"
)

type MockStore struct {
	AccountStore    *mockstore.MockAccountStore
	RepositoryStore *mockstore.MockRepositoryStore
	ScanStore       *mockstore.MockScanStore
}

func (m MockStore) Account() store.AccountStore {
	return m.AccountStore
}

func (m MockStore) Repository() store.RepositoryStore {
	return m.RepositoryStore
}

func (m MockStore) Scan() store.ScanStore {
	return m.ScanStore
}

func (m MockStore) MigrationDatabaseSchema() error {
	return nil
}

type TestHelper struct {
	a     *app.App
	store *MockStore
}

func Setup(t *testing.T) (*TestHelper, error) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := MockStore{
		AccountStore:    mockstore.NewMockAccountStore(ctrl),
		RepositoryStore: mockstore.NewMockRepositoryStore(ctrl),
		ScanStore:       mockstore.NewMockScanStore(ctrl),
	}

	js, err := jobs.NewJobServer(ms)
	if err != nil {
		return nil, err
	}

	a := app.App{
		Store:     ms,
		JobServer: js,
	}

	th := TestHelper{
		a:     &a,
		store: &ms,
	}
	return &th, nil
}
