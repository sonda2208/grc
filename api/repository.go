package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sonda2208/guardrails-challenge/model"
)

func (a *API) InitRepository() {
	a.Routes.RepositoriesForAccount.POST("", a.AddRepository)
	a.Routes.RepositoriesForAccount.GET("", a.GetRepositories)

	a.Routes.Repository.PATCH("", a.PatchRepository)
	a.Routes.Repository.DELETE("", a.DeleteRepository)
}

// AddRepository
// @Summary Add a repository
// @Tags Repository
// @Accept json
// @Produce json
// @Param message body model.AddRepositoryPayload true "Payload"
// @Success 200 {object} model.Repository
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/accounts/:accountID/repos [post]
func (a API) AddRepository(c echo.Context) error {
	accountID := paramFromContext(c).accountID
	p := model.AddRepositoryPayload{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.WithError("api.repo.add.bind_fail", err))
	}

	repo, err := a.App.AddRepository(accountID, &p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, repo)
}

// GetRepositories
// @Summary Return list of repositories for an account
// @Tags Repository
// @Produce json
// @Param q query model.ListRepositoriesOption false "Options"
// @Success 200 {array} model.Repository
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/accounts/:accountID/repos [get]
func (a API) GetRepositories(c echo.Context) error {
	accountID := paramFromContext(c).accountID
	opt := model.ListRepositoriesOption{}
	err := c.Bind(&opt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.WithError("api.repo.list.bind_fail", err))
	}

	err = opt.IsValid()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	repos, cnt, err := a.App.GetRepositories(accountID, &opt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	res := model.ListingResult{
		Data:       repos,
		TotalCount: cnt,
		Page:       opt.Page,
		PerPage:    opt.PerPage,
	}
	return c.JSON(http.StatusOK, res)

}

// PatchRepository
// @Summary Update repository info
// @Tags Repository
// @Accept json
// @Produce json
// @Param message body model.RepositoryPatch true "Payload"
// @Success 200 {object} model.Repository
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/repos/:repoID [patch]
func (a API) PatchRepository(c echo.Context) error {
	repoID := paramFromContext(c).repoID
	p := model.RepositoryPatch{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.WithError("api.repo.patch.bind_fail", err))
	}

	err = a.App.PatchRepository(repoID, &p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	repo, err := a.App.GetRepository(repoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repo)
}

// DeleteRepository
// @Summary Delete a repository
// @Tags Repository
// @Produce json
// @Success 200
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/repos/:repoID [delete]
func (a API) DeleteRepository(c echo.Context) error {
	repoID := paramFromContext(c).repoID
	err := a.App.DeleteRepository(repoID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
