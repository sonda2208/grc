package model

import "time"

type Account struct {
	ID        int        `json:"id" db:"id,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty" db:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at,omitempty"`
	Email     string     `json:"email" db:"email"`
	FirstName string     `json:"firstName" db:"first_name"`
	LastName  string     `json:"lastName" db:"last_name"`
}

func (a Account) IsValid() error {
	var options []ErrorOption

	if !IsValidEmail(a.Email) {
		options = append(options, WithValidateError("email", ErrInvalidValue))
	}

	if len(options) > 0 {
		return NewError("model.account.invalid", "Invalid account", options...)
	}

	return nil
}

func (a *Account) PreSave() {
	a.Email = NormalizeEmail(a.Email)
}

func (a *Account) PreUpdate() {
	a.UpdatedAt = time.Now().UTC()
}

func (a *Account) Patch(p *AccountPatch) {
	if p.FirstName != nil {
		a.FirstName = *p.FirstName
	}

	if p.LastName != nil {
		a.LastName = *p.LastName
	}
}

type AccountSignupPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type AccountPatch struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}
