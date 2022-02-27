package sqlstore_test

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/sonda2208/guardrails-challenge/model"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	th, err := Setup()
	require.NoError(t, err)
	defer th.Teardown()

	t.Run("create account", func(t *testing.T) {
		acc := model.Account{
			Email:     fake.EmailAddress(),
			FirstName: fake.FirstName(),
			LastName:  fake.LastName(),
		}
		ra, err := th.Store.Account().Save(&acc)
		require.NoError(t, err)
		assert.NotZero(t, ra.ID)
		assert.Equal(t, acc.Email, ra.Email)
		assert.Equal(t, acc.FirstName, ra.FirstName)
		assert.Equal(t, acc.LastName, ra.LastName)
		assert.NotZero(t, acc.CreatedAt)
		assert.NotNil(t, acc.UpdatedAt)
	})

	t.Run("create account with empty email", func(t *testing.T) {
		acc := model.Account{
			FirstName: fake.FirstName(),
			LastName:  fake.LastName(),
		}
		_, err := th.Store.Account().Save(&acc)
		assert.Error(t, err)
	})

	t.Run("update account", func(t *testing.T) {
		a, err := th.CreateAccount()
		require.NoError(t, err)

		expected := fake.LastName()
		expectedUpdateTime := a.UpdatedAt

		a.LastName = expected
		err = th.Store.Account().Update(a)
		require.NoError(t, err)

		ra, err := th.Store.Account().Get(a.ID)
		require.NoError(t, err)
		assert.Equal(t, expected, ra.LastName)
		assert.True(t, !ra.UpdatedAt.Before(expectedUpdateTime))
	})

	t.Run("get by id", func(t *testing.T) {
		a, err := th.CreateAccount()
		require.NoError(t, err)

		ra, err := th.Store.Account().Get(a.ID)
		require.NoError(t, err)
		assert.Equal(t, a.Email, ra.Email)
		assert.Equal(t, a.FirstName, ra.FirstName)
		assert.Equal(t, a.LastName, ra.LastName)
	})

	t.Run("get non-existing account by id", func(t *testing.T) {
		_, err := th.Store.Account().Get(0)
		assert.Error(t, err)
	})

	t.Run("get by email", func(t *testing.T) {
		a, err := th.CreateAccount()
		require.NoError(t, err)

		ra, err := th.Store.Account().GetByEmail(a.Email)
		require.NoError(t, err)
		assert.Equal(t, a.Email, ra.Email)
		assert.Equal(t, a.FirstName, ra.FirstName)
		assert.Equal(t, a.LastName, ra.LastName)
	})

	t.Run("get non-existing account by email", func(t *testing.T) {
		_, err := th.Store.Account().GetByEmail("xxx@domain.com")
		assert.Error(t, err)
	})
}
