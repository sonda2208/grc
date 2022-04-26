package api

import (
	"net/http"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"github.com/sonda2208/grc/app"
)

type Routes struct {
	Root *echo.Echo
	API  *echo.Group // "/api"

	Accounts *echo.Group // "/api/accounts"
	Account  *echo.Group // "/api/accounts/:accountID"

	RepositoriesForAccount *echo.Group // "/api/accounts/:accountID/repos"
	Repository             *echo.Group // "/api/repos/:repoID"

	ScansForRepo *echo.Group // "/api/repos/:repoID/scans"
	Scan         *echo.Group // "/api/scans/:scanID"
}

type API struct {
	App    *app.App
	Routes *Routes
}

func New(app *app.App, e *echo.Echo) (*API, error) {
	a := &API{
		App: app,
		Routes: &Routes{
			Root: e,
		},
	}

	a.Routes.Root.Use(middleware.RequestID())
	a.Routes.Root.GET("/health", a.Health)

	a.Routes.API = a.Routes.Root.Group("/api")

	a.Routes.Accounts = a.Routes.API.Group("/accounts")
	a.Routes.Account = a.Routes.API.Group("/accounts/:accountID")

	a.Routes.RepositoriesForAccount = a.Routes.Account.Group("/repos")
	a.Routes.Repository = a.Routes.API.Group("/repos/:repoID")

	a.Routes.ScansForRepo = a.Routes.Repository.Group("/scans")
	a.Routes.Scan = a.Routes.API.Group("/scans/:scanID")

	a.InitAccount()
	a.InitRepository()
	a.InitScan()

	return a, nil
}

func (a API) Health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
