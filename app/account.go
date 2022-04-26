package app

import (
	"github.com/sonda2208/grc/model"
)

func (a App) CreateAccountFromSignup(p *model.AccountSignupPayload) (*model.Account, error) {
	acc := model.Account{
		Email:     p.Email,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}
	ra, err := a.Store.Account().Save(&acc)
	if err != nil {
		return nil, err
	}

	return ra, nil
}

func (a App) GetAccountByID(id int) (*model.Account, error) {
	ra, err := a.Store.Account().Get(id)
	if err != nil {
		return nil, err
	}

	return ra, nil
}

func (a App) PatchAccount(id int, p *model.AccountPatch) error {
	acc, err := a.GetAccountByID(id)
	if err != nil {
		return err
	}

	acc.Patch(p)
	err = a.Store.Account().Update(acc)
	if err != nil {
		return err
	}

	return nil
}
