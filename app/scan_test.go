package app_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/icrowley/fake"
	"github.com/sonda2208/guardrails-challenge/model"

	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	th, err := Setup(t)
	require.NoError(t, err)

	th.store.ScanStore.EXPECT().Get(1).Return(nil, sql.ErrNoRows)
	th.store.ScanStore.EXPECT().Get(2).Return(&model.Scan{
		ID:     2,
		RepoID: 1,
		Branch: fake.Characters(),
		Commit: fake.Characters(),
		Status: model.ScanStatusQueued,
	}, nil).Times(2)
	th.store.ScanStore.EXPECT().Get(3).Return(&model.Scan{
		ID:     3,
		RepoID: 1,
		Branch: fake.Characters(),
		Commit: fake.Characters(),
		Status: model.ScanStatusInProgress,
	}, nil)
	th.store.ScanStore.EXPECT().Update(gomock.Any()).Return(nil)

	testCases := []struct {
		scanID  int
		isError bool
	}{
		{
			scanID:  1,
			isError: true,
		},
		{
			scanID:  2,
			isError: false,
		},
		{
			scanID:  3,
			isError: true,
		},
	}
	for _, tc := range testCases {
		err := th.a.RerunScan(tc.scanID)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

}
