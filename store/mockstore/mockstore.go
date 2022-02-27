// Code generated by MockGen. DO NOT EDIT.
// Source: store/store.go

// Package mockstore is a generated GoMock package.
package mockstore

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/sonda2208/guardrails-challenge/model"
	store "github.com/sonda2208/guardrails-challenge/store"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Account mocks base method.
func (m *MockStore) Account() store.AccountStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Account")
	ret0, _ := ret[0].(store.AccountStore)
	return ret0
}

// Account indicates an expected call of Account.
func (mr *MockStoreMockRecorder) Account() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Account", reflect.TypeOf((*MockStore)(nil).Account))
}

// MigrationDatabaseSchema mocks base method.
func (m *MockStore) MigrationDatabaseSchema() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MigrationDatabaseSchema")
	ret0, _ := ret[0].(error)
	return ret0
}

// MigrationDatabaseSchema indicates an expected call of MigrationDatabaseSchema.
func (mr *MockStoreMockRecorder) MigrationDatabaseSchema() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MigrationDatabaseSchema", reflect.TypeOf((*MockStore)(nil).MigrationDatabaseSchema))
}

// Repository mocks base method.
func (m *MockStore) Repository() store.RepositoryStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Repository")
	ret0, _ := ret[0].(store.RepositoryStore)
	return ret0
}

// Repository indicates an expected call of Repository.
func (mr *MockStoreMockRecorder) Repository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Repository", reflect.TypeOf((*MockStore)(nil).Repository))
}

// Scan mocks base method.
func (m *MockStore) Scan() store.ScanStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan")
	ret0, _ := ret[0].(store.ScanStore)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockStoreMockRecorder) Scan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockStore)(nil).Scan))
}

// MockAccountStore is a mock of AccountStore interface.
type MockAccountStore struct {
	ctrl     *gomock.Controller
	recorder *MockAccountStoreMockRecorder
}

// MockAccountStoreMockRecorder is the mock recorder for MockAccountStore.
type MockAccountStoreMockRecorder struct {
	mock *MockAccountStore
}

// NewMockAccountStore creates a new mock instance.
func NewMockAccountStore(ctrl *gomock.Controller) *MockAccountStore {
	mock := &MockAccountStore{ctrl: ctrl}
	mock.recorder = &MockAccountStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountStore) EXPECT() *MockAccountStoreMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockAccountStore) Get(id int) (*model.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*model.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAccountStoreMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAccountStore)(nil).Get), id)
}

// GetByEmail mocks base method.
func (m *MockAccountStore) GetByEmail(email string) (*model.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(*model.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockAccountStoreMockRecorder) GetByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockAccountStore)(nil).GetByEmail), email)
}

// Save mocks base method.
func (m *MockAccountStore) Save(a *model.Account) (*model.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", a)
	ret0, _ := ret[0].(*model.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockAccountStoreMockRecorder) Save(a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockAccountStore)(nil).Save), a)
}

// Update mocks base method.
func (m *MockAccountStore) Update(a *model.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", a)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAccountStoreMockRecorder) Update(a interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAccountStore)(nil).Update), a)
}

// MockRepositoryStore is a mock of RepositoryStore interface.
type MockRepositoryStore struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryStoreMockRecorder
}

// MockRepositoryStoreMockRecorder is the mock recorder for MockRepositoryStore.
type MockRepositoryStoreMockRecorder struct {
	mock *MockRepositoryStore
}

// NewMockRepositoryStore creates a new mock instance.
func NewMockRepositoryStore(ctrl *gomock.Controller) *MockRepositoryStore {
	mock := &MockRepositoryStore{ctrl: ctrl}
	mock.recorder = &MockRepositoryStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryStore) EXPECT() *MockRepositoryStoreMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRepositoryStore) Delete(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryStoreMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepositoryStore)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockRepositoryStore) Get(id int) (*model.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*model.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryStoreMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepositoryStore)(nil).Get), id)
}

// GetByAccount mocks base method.
func (m *MockRepositoryStore) GetByAccount(accountID int, opt *model.ListRepositoriesOption) ([]*model.Repository, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAccount", accountID, opt)
	ret0, _ := ret[0].([]*model.Repository)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByAccount indicates an expected call of GetByAccount.
func (mr *MockRepositoryStoreMockRecorder) GetByAccount(accountID, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAccount", reflect.TypeOf((*MockRepositoryStore)(nil).GetByAccount), accountID, opt)
}

// Save mocks base method.
func (m *MockRepositoryStore) Save(r *model.Repository) (*model.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", r)
	ret0, _ := ret[0].(*model.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockRepositoryStoreMockRecorder) Save(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepositoryStore)(nil).Save), r)
}

// Update mocks base method.
func (m *MockRepositoryStore) Update(r *model.Repository) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", r)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryStoreMockRecorder) Update(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepositoryStore)(nil).Update), r)
}

// MockScanStore is a mock of ScanStore interface.
type MockScanStore struct {
	ctrl     *gomock.Controller
	recorder *MockScanStoreMockRecorder
}

// MockScanStoreMockRecorder is the mock recorder for MockScanStore.
type MockScanStoreMockRecorder struct {
	mock *MockScanStore
}

// NewMockScanStore creates a new mock instance.
func NewMockScanStore(ctrl *gomock.Controller) *MockScanStore {
	mock := &MockScanStore{ctrl: ctrl}
	mock.recorder = &MockScanStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScanStore) EXPECT() *MockScanStoreMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockScanStore) Get(id int) (*model.Scan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*model.Scan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockScanStoreMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockScanStore)(nil).Get), id)
}

// GetByRepository mocks base method.
func (m *MockScanStore) GetByRepository(repoID int, opt *model.ListScansOption) ([]*model.Scan, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRepository", repoID, opt)
	ret0, _ := ret[0].([]*model.Scan)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByRepository indicates an expected call of GetByRepository.
func (mr *MockScanStoreMockRecorder) GetByRepository(repoID, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRepository", reflect.TypeOf((*MockScanStore)(nil).GetByRepository), repoID, opt)
}

// GetByStatus mocks base method.
func (m *MockScanStore) GetByStatus(status string) ([]*model.Scan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByStatus", status)
	ret0, _ := ret[0].([]*model.Scan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByStatus indicates an expected call of GetByStatus.
func (mr *MockScanStoreMockRecorder) GetByStatus(status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByStatus", reflect.TypeOf((*MockScanStore)(nil).GetByStatus), status)
}

// Save mocks base method.
func (m *MockScanStore) Save(s *model.Scan) (*model.Scan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", s)
	ret0, _ := ret[0].(*model.Scan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockScanStoreMockRecorder) Save(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockScanStore)(nil).Save), s)
}

// Update mocks base method.
func (m *MockScanStore) Update(s *model.Scan) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", s)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockScanStoreMockRecorder) Update(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockScanStore)(nil).Update), s)
}