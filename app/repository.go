package app

import (
	"github.com/sonda2208/grc/model"
)

func (a App) AddRepository(accountID int, p *model.AddRepositoryPayload) (*model.Repository, error) {
	r := model.Repository{
		AccountID: accountID,
		Name:      p.Name,
		URL:       p.URL,
	}
	rr, err := a.Store.Repository().Save(&r)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func (a App) GetRepository(repoID int) (*model.Repository, error) {
	rr, err := a.Store.Repository().Get(repoID)
	if err != nil {
		return nil, err
	}

	return rr, nil
}

func (a App) GetRepositories(accountID int, opt *model.ListRepositoriesOption) ([]*model.Repository, int, error) {
	repos, cnt, err := a.Store.Repository().GetByAccount(accountID, opt)
	if err != nil {
		return nil, 0, err
	}

	return repos, cnt, nil
}

func (a App) PatchRepository(repoID int, p *model.RepositoryPatch) error {
	r, err := a.GetRepository(repoID)
	if err != nil {
		return err
	}

	r.Patch(p)
	err = a.Store.Repository().Update(r)
	if err != nil {
		return err
	}

	return nil
}

func (a App) DeleteRepository(repoID int) error {
	err := a.Store.Repository().Delete(repoID)
	if err != nil {
		return err
	}

	return nil
}
