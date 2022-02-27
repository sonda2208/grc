package app

import (
	"github.com/sonda2208/guardrails-challenge/model"
)

func (a App) ListScans(repoID int, opt *model.ListScansOption) ([]*model.Scan, int, error) {
	scans, cnt, err := a.Store.Scan().GetByRepository(repoID, opt)
	if err != nil {
		return nil, 0, err
	}

	return scans, cnt, nil
}

func (a App) GetScanByID(id int) (*model.Scan, error) {
	scan, err := a.Store.Scan().Get(id)
	if err != nil {
		return nil, err
	}

	return scan, nil
}

func (a App) CreateScan(repoID int, p *model.CreateScanPayload) (*model.Scan, error) {
	scan := model.Scan{
		RepoID: repoID,
		Type:   p.Type,
		Branch: p.Branch,
		Commit: p.Commit,
		Status: model.ScanStatusQueued,
	}
	rs, err := a.Store.Scan().Save(&scan)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func (a App) RerunScan(scanID int) error {
	scan, err := a.GetScanByID(scanID)
	if err != nil {
		return err
	}

	if scan.Status == model.ScanStatusInProgress {
		return model.NewError("app.scan.rerun_in_progress_scan", "Unable to re-run a running scan.")
	}

	err = a.js.SetScanStatus(scanID, model.ScanStatusQueued, "re-run")
	if err != nil {
		return err
	}

	return nil
}
