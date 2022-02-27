package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sonda2208/guardrails-challenge/model"
)

func (a *API) InitScan() {
	a.Routes.ScansForRepo.GET("", a.ListScans)
	a.Routes.ScansForRepo.POST("", a.CreateScan)

	a.Routes.Scan.PUT("", a.RerunScan)
}

// ListScans
// @Summary Return list of scans of a repository
// @Tags Scan
// @Produce json
// @Param q query model.ListScansOption false "Options"
// @Success 200 {array} model.Scan
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/repos/:repoID/scans [get]
func (a API) ListScans(c echo.Context) error {
	repoID := paramFromContext(c).repoID
	opt := model.ListScansOption{}
	err := c.Bind(&opt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.WithError("api.scan.list.bind_fail", err))
	}

	err = opt.IsValid()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	scans, cnt, err := a.App.ListScans(repoID, &opt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	res := model.ListingResult{
		Data:       scans,
		TotalCount: cnt,
		Page:       opt.Page,
		PerPage:    opt.PerPage,
	}
	return c.JSON(http.StatusOK, res)
}

// CreateScan
// @Summary Add a repository
// @Tags Scan
// @Accept json
// @Produce json
// @Param message body model.CreateScanPayload true "Payload"
// @Success 200 {object} model.Scan
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/repos/:repoID/scans [post]
func (a API) CreateScan(c echo.Context) error {
	repoID := paramFromContext(c).repoID
	p := model.CreateScanPayload{}
	err := c.Bind(&p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.WithError("api.scan.add.bind_fail", err))
	}

	scan, err := a.App.CreateScan(repoID, &p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, scan)
}

// RerunScan
// @Summary Re-run a scan
// @Tags Scan
// @Produce json
// @Success 200 {object} model.Scan
// @Failure 400 {object} model.Error
// @Failure 500 {object} model.Error
// @Router /api/scans/:scanID [put]
func (a API) RerunScan(c echo.Context) error {
	scanID := paramFromContext(c).scanID
	err := a.App.RerunScan(scanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	scan, err := a.App.GetScanByID(scanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, scan)
}
