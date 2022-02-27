package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ScanStatusQueued     = "Queued"
	ScanStatusInProgress = "In Progress"
	ScanStatusSuccess    = "Success"
	ScanStatusFailure    = "Failure"
)

var ScanSortableFields = map[string]string{
	"id":        "id",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
	"branch":    "branch",
	"status":    "status",
}

type FindingMetadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type PositionIndex struct {
	Line int `json:"line"`
}

type FindingPosition struct {
	Begin PositionIndex `json:"begin"`
}

type FindingLocation struct {
	Path      string          `json:"path"`
	Positions FindingPosition `json:"positions"`
}

type Finding struct {
	Type     string          `json:"type"`
	Location FindingLocation `json:"location"`
	Metadata FindingMetadata `json:"metadata"`
}

type Findings []*Finding

func (f Findings) Value() (driver.Value, error) {
	res, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (f *Findings) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return NewError("model.scan.finding_db_scan", "type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	return nil
}

type Scan struct {
	ID         int        `json:"id" db:"id,omitempty"`
	RepoID     int        `json:"repoId" db:"repo_id"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at,omitempty"`
	ScanningAt *time.Time `json:"scanningAt,omitempty" db:"scanning_at,omitempty"`
	FinishedAt *time.Time `json:"finishedAt,omitempty" db:"finished_at,omitempty"`
	Branch     string     `json:"branch" db:"branch"`
	Commit     string     `json:"commit" db:"commit"`
	Status     string     `json:"status" db:"status"`
	Message    string     `json:"message" db:"message"`
	Findings   Findings   `json:"findings" db:"findings"`
}

func (s Scan) IsValid() error {
	var options []ErrorOption

	if s.Branch == "" {
		options = append(options, WithValidateError("branch", ErrInvalidValue))
	}

	if s.Commit == "" {
		options = append(options, WithValidateError("commit", ErrInvalidValue))
	}

	err := IsValidScanStatus(s.Status)
	if err != nil {
		options = append(options, WithValidateError("status", ErrInvalidValue))
	}

	if len(options) > 0 {
		return NewError("model.scan.invalid", "Invalid scan", options...)
	}

	return nil
}

func (s Scan) RepoLockName() string {
	return fmt.Sprintf("%d-%s-%s", s.RepoID, s.Branch, s.Commit)
}

func IsValidScanStatus(status string) error {
	switch status {
	case ScanStatusQueued:
	case ScanStatusInProgress:
	case ScanStatusSuccess:
	case ScanStatusFailure:
	default:
		return NewError("model.scan.invalid_status", "Invalid status")
	}

	return nil
}

type ListScansOption struct {
	Page    int    `query:"page"`
	PerPage int    `query:"perPage"`
	Branch  string `query:"branch"`
	Status  string `query:"status"`
	OrderBy string `query:"orderBy"`
	IsDesc  bool   `query:"isDesc"`
}

func (o ListScansOption) IsValid() error {
	var options []ErrorOption

	if o.OrderBy != "" {
		if _, ok := ScanSortableFields[o.OrderBy]; !ok {
			options = append(options, WithValidateError("orderBy", ErrInvalidValue))
		}
	}

	return nil
}

type CreateScanPayload struct {
	Branch string `json:"branch"`
	Commit string `json:"commit"`
	Type   string `json:"type"`
}

type ScanResult struct {
	Findings Findings
	Error    error
}
