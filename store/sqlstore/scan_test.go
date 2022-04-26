package sqlstore_test

import (
	"testing"
	"time"

	"github.com/icrowley/fake"

	"github.com/stretchr/testify/assert"

	"github.com/sonda2208/grc/model"

	"github.com/stretchr/testify/require"
)

func TestScan(t *testing.T) {
	th, err := Setup()
	require.NoError(t, err)
	defer th.Teardown()

	err = th.InitBasic()
	require.NoError(t, err)

	t.Run("create scan", func(t *testing.T) {
		s := model.Scan{
			RepoID: th.SampleRepository.ID,
			Branch: "dev",
			Commit: "d41727b",
			Status: model.ScanStatusQueued,
		}

		rs, err := th.Store.Scan().Save(&s)
		require.NoError(t, err)
		assert.NotZero(t, rs.ID)
		assert.Equal(t, th.SampleRepository.ID, rs.RepoID)
		assert.Equal(t, s.Branch, rs.Branch)
		assert.Equal(t, s.Commit, rs.Commit)
		assert.Equal(t, s.Status, rs.Status)
		assert.NotNil(t, rs.CreatedAt)
	})

	t.Run("create scan with missing fields", func(t *testing.T) {
		s := model.Scan{
			RepoID: th.SampleRepository.ID,
			Status: model.ScanStatusQueued,
		}

		_, err := th.Store.Scan().Save(&s)
		assert.Error(t, err)

		appErr, ok := err.(*model.Error)
		require.True(t, ok)
		assert.NotEmpty(t, appErr.ValidationErrors)
	})

	t.Run("create scan with invalid status", func(t *testing.T) {
		s := model.Scan{
			RepoID: th.SampleRepository.ID,
			Branch: "dev",
			Commit: "d41727b",
		}

		_, err := th.Store.Scan().Save(&s)
		assert.Error(t, err)
	})

	t.Run("update scan status", func(t *testing.T) {
		s, err := th.CreateScan(th.SampleRepository.ID)
		require.NoError(t, err)

		s.Status = model.ScanStatusSuccess
		err = th.Store.Scan().Update(s)
		require.NoError(t, err)

		rs, err := th.Store.Scan().Get(s.ID)
		require.NoError(t, err)
		assert.Equal(t, s.Status, rs.Status)
	})

	t.Run("update enqueue time", func(t *testing.T) {
		s, err := th.CreateScan(th.SampleRepository.ID)
		require.NoError(t, err)

		now := time.Now().UTC()
		s.FinishedAt = &now
		err = th.Store.Scan().Update(s)
		require.NoError(t, err)

		rs, err := th.Store.Scan().Get(s.ID)
		require.NoError(t, err)

		assert.True(t, now.Round(time.Microsecond).Equal(*rs.FinishedAt))
	})

	t.Run("add findings", func(t *testing.T) {
		s, err := th.CreateScan(th.SampleRepository.ID)
		require.NoError(t, err)

		s.Findings = []*model.Finding{
			{
				Type:     "test",
				Location: model.FindingLocation{},
				Metadata: model.FindingMetadata{},
			},
		}
		err = th.Store.Scan().Update(s)
		require.NoError(t, err)

		rs, err := th.Store.Scan().Get(s.ID)
		require.NoError(t, err)
		assert.EqualValues(t, s.Findings, rs.Findings)
	})

	t.Run("get by id", func(t *testing.T) {
		s, err := th.CreateScan(th.SampleRepository.ID)
		require.NoError(t, err)

		rs, err := th.Store.Scan().Get(s.ID)
		require.NoError(t, err)

		assert.Equal(t, rs.ID, s.ID)
		assert.Equal(t, rs.RepoID, s.RepoID)
		assert.Equal(t, rs.Branch, s.Branch)
		assert.Equal(t, rs.Commit, s.Commit)
		assert.Equal(t, rs.Status, s.Status)
	})

	t.Run("get non-existing scan", func(t *testing.T) {
		_, err := th.Store.Scan().Get(0)
		assert.Error(t, err)
	})

	t.Run("list scans by repo", func(t *testing.T) {
		repo, err := th.CreateRepository(th.SampleAccount.ID)
		require.NoError(t, err)

		const n = 3
		for i := 0; i < n; i++ {
			_, err := th.CreateScan(repo.ID)
			require.NoError(t, err)
		}

		scans, cnt, err := th.Store.Scan().GetByRepository(repo.ID, &model.ListScansOption{
			Page:    1,
			PerPage: n,
		})
		require.NoError(t, err)
		assert.Equal(t, n, cnt)
		assert.Equal(t, n, len(scans))
	})

	t.Run("list scans by branch for repo", func(t *testing.T) {
		repo, err := th.CreateRepository(th.SampleAccount.ID)
		require.NoError(t, err)

		const expectedBranch = "xxx"
		const n = 3
		const m = 5
		for i := 0; i < n; i++ {
			s := model.Scan{
				RepoID:    repo.ID,
				CreatedAt: time.Time{},
				Branch:    expectedBranch,
				Commit:    fake.HexColor(),
				Status:    model.ScanStatusQueued,
			}
			_, err := th.Store.Scan().Save(&s)
			require.NoError(t, err)
		}
		for i := 0; i < m; i++ {
			_, err := th.CreateScan(repo.ID)
			require.NoError(t, err)
		}

		scans, cnt, err := th.Store.Scan().GetByRepository(repo.ID, &model.ListScansOption{
			Page:    1,
			PerPage: n,
			Branch:  expectedBranch,
		})
		require.NoError(t, err)
		assert.Equal(t, n, cnt)
		assert.Equal(t, n, len(scans))
	})

	t.Run("list all scans by status", func(t *testing.T) {
		scans, err := th.Store.Scan().GetByStatus(model.ScanStatusQueued)
		require.NoError(t, err)

		for _, s := range scans {
			assert.Equal(t, model.ScanStatusQueued, s.Status)
		}
	})
}
