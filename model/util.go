package model

import (
	"net/mail"
	"strings"
)

type ListingResult struct {
	Data       interface{} `json:"data"`
	TotalCount int         `json:"totalCount"`
	Page       int         `json:"page"`
	PerPage    int         `json:"perPage"`
}

func IsLower(s string) bool {
	return strings.ToLower(s) == s
}

func IsValidEmail(email string) bool {
	if !IsLower(email) {
		return false
	}

	if addr, err := mail.ParseAddress(email); err != nil {
		return false
	} else if addr.Name != "" {
		return false
	}

	return true
}

func NormalizeEmail(email string) string {
	email = strings.TrimSpace(email)
	return strings.ToLower(email)
}
