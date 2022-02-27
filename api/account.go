package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sonda2208/guardrails-challenge/model"
)

func (a *API) InitAccount() {
	a.Routes.Accounts.POST("/signup", a.Signup)

	a.Routes.Account.GET("", a.GetAccount)
	a.Routes.Account.PATCH("", a.PatchAccount)
}

// Signup
// @Summary Signup
// @Tags Account
// @Accept json
// @Produce json
// @Param message body model.AccountSignupPayload true "Payload"
// @Success 200 {object} model.Account
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/accounts/signup [post]
func (a API) Signup(c echo.Context) error {
	p := model.AccountSignupPayload{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.WithError("api.account.signup.bind_fail", err))
	}

	acc, err := a.App.CreateAccountFromSignup(&p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, acc)
}

// GetAccount
// @Summary Return account info
// @Tags Account
// @Produce json
// @Success 200 {object} model.Account
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/accounts/:accountID [get]
func (a API) GetAccount(c echo.Context) error {
	accountID := paramFromContext(c).accountID
	acc, err := a.App.GetAccountByID(accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, acc)
}

// PatchAccount
// @Summary Update account info
// @Tags Account
// @Accept json
// @Produce json
// @Param message body model.AccountPatch true "Payload"
// @Success 200 {object} model.Account
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/accounts/:accountID [patch]
func (a API) PatchAccount(c echo.Context) error {
	accountID := paramFromContext(c).accountID
	p := model.AccountPatch{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.WithError("api.account.patch.bind_fail", err))
	}

	err = a.App.PatchAccount(accountID, &p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	acc, err := a.App.GetAccountByID(accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, acc)
}
