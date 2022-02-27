package model

import "time"

var RepositorySortableFields = map[string]string{
	"id":        "id",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
	"name":      "name",
}

type Repository struct {
	ID        int        `json:"id" db:"id,omitempty"`
	AccountID int        `json:"accountId" db:"account_id"`
	CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty" db:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at,omitempty"`
	Name      string     `json:"name" db:"name"`
	URL       string     `json:"url" db:"url"`
}

func (r Repository) IsValid() error {
	var options []ErrorOption

	if r.Name == "" {
		options = append(options, WithValidateError("name", ErrInvalidValue))
	}

	/*
		TODO:
		- Validate git clone URL
	*/
	if r.URL == "" {
		options = append(options, WithValidateError("url", ErrInvalidValue))
	}

	if len(options) > 0 {
		return NewError("model.repository.invalid", "Invalid repository", options...)
	}

	return nil
}

func (r *Repository) PreUpdate() {
	r.UpdatedAt = time.Now().UTC()
}

func (r *Repository) Patch(p *RepositoryPatch) {
	if p.Name != nil {
		r.Name = *p.Name
	}

	if p.URL != nil {
		r.URL = *p.URL
	}
}

type ListRepositoriesOption struct {
	Page    int    `query:"page"`
	PerPage int    `query:"perPage"`
	OrderBy string `query:"orderBy"`
	IsDesc  bool   `query:"isDesc"`
}

func (o ListRepositoriesOption) IsValid() error {
	var options []ErrorOption

	if o.OrderBy != "" {
		if _, ok := RepositorySortableFields[o.OrderBy]; !ok {
			options = append(options, WithValidateError("orderBy", ErrInvalidValue))
		}
	}

	if len(options) > 0 {
		return NewError("model.repository.invalid_listing_option", "Invalid listing option", options...)
	}

	return nil
}

type AddRepositoryPayload struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type RepositoryPatch struct {
	Name *string `json:"name"`
	URL  *string `json:"url"`
}
