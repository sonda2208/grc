package api

import (
	"strconv"

	"github.com/labstack/echo"
)

type PathParam struct {
	accountID int
	repoID    int
	scanID    int
}

func paramFromContext(c echo.Context) PathParam {
	pp := PathParam{}
	accountID, err := strconv.Atoi(c.Param("accountID"))
	if err == nil {
		pp.accountID = accountID
	}

	repoID, err := strconv.Atoi(c.Param("repoID"))
	if err == nil {
		pp.repoID = repoID
	}

	scanID, err := strconv.Atoi(c.Param("scanID"))
	if err == nil {
		pp.scanID = scanID
	}

	return pp
}
