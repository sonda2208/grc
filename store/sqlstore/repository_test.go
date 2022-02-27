package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/icrowley/fake"

	"github.com/sonda2208/guardrails-challenge/model"

	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	th, err := Setup()
	require.NoError(t, err)
	defer th.Teardown()

	err = th.InitBasic()
	require.NoError(t, err)

	t.Run("create repository", func(t *testing.T) {
		r := model.Repository{
			AccountID: th.SampleAccount.ID,
			Name:      fake.ProductName(),
			URL:       "https://git.com/repo",
		}
		rr, err := th.Store.Repository().Save(&r)
		require.NoError(t, err)
		assert.NotZero(t, rr.ID)
		assert.Equal(t, r.AccountID, rr.AccountID)
		assert.Equal(t, r.Name, rr.Name)
		assert.Equal(t, r.URL, rr.URL)
		assert.NotZero(t, r.CreatedAt, rr.CreatedAt)
		assert.NotNil(t, r.UpdatedAt, rr.UpdatedAt)
	})

	t.Run("create repository with missing fields", func(t *testing.T) {
		r := model.Repository{
			AccountID: th.SampleAccount.ID,
		}
		_, err := th.Store.Repository().Save(&r)
		assert.Error(t, err)

		appErr, ok := err.(*model.Error)
		require.True(t, ok)
		assert.NotEmpty(t, appErr.ValidationErrors)
	})

	t.Run("update repository", func(t *testing.T) {
		r, err := th.CreateRepository(th.SampleAccount.ID)
		require.NoError(t, err)

		expectedUpdateTime := r.UpdatedAt
		expectedName := fake.ProductName()

		r.Name = expectedName
		err = th.Store.Repository().Update(r)
		require.NoError(t, err)

		rr, err := th.Store.Repository().Get(r.ID)
		require.NoError(t, err)
		assert.Equal(t, rr.Name, expectedName)
		assert.True(t, !rr.UpdatedAt.Before(expectedUpdateTime))
	})

	t.Run("get repo by id", func(t *testing.T) {
		r, err := th.CreateRepository(th.SampleAccount.ID)
		require.NoError(t, err)

		rr, err := th.Store.Repository().Get(r.ID)
		require.NoError(t, err)
		assert.Equal(t, rr.ID, r.ID)
		assert.Equal(t, rr.Name, r.Name)
		assert.Equal(t, rr.URL, r.URL)
	})

	t.Run("get non-existing repo by id", func(t *testing.T) {
		_, err := th.Store.Repository().Get(0)
		assert.Error(t, err)
	})

	t.Run("list repos for account", func(t *testing.T) {
		a, err := th.CreateAccount()
		require.NoError(t, err)

		const n = 3
		for i := 0; i < n; i++ {
			_, err := th.CreateRepository(a.ID)
			require.NoError(t, err)
		}

		repos, cnt, err := th.Store.Repository().GetByAccount(a.ID, &model.ListRepositoriesOption{
			Page:    1,
			PerPage: n,
		})
		require.NoError(t, err)
		assert.Equal(t, n, cnt)
		assert.Equal(t, n, len(repos))
	})

	t.Run("delete repo", func(t *testing.T) {
		r, err := th.CreateRepository(th.SampleAccount.ID)
		require.NoError(t, err)

		err = th.Store.Repository().Delete(r.ID)
		require.NoError(t, err)

		_, err = th.Store.Repository().Get(r.ID)
		assert.Error(t, err)
	})
}
