package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Video ...
type Video struct {
	ID           int             `json:"id"`
	SubDirection DevSubDirection `json:"subDirectionId"`
	Diffuculty   DiffucultyLevel `json:"diffuculty"`
	Description  string          `json:"description"`
	Language     string          `json:"language"`
	URL          string          `json:"url"`
}

// VideoDTO ...
type VideoDTO struct {
	ID             int    `json:"id"`
	SubDirectionID int    `json:"subDirectionId"`
	Diffuculty     string `json:"diffuculty"`
	Description    string `json:"description"`
	Language       string `json:"language"`
	URL            string `json:"url"`
}

// Validate ...
func (v *VideoDTO) Validate() error {
	return validation.ValidateStruct(
		v,
		validation.Field(&v.SubDirectionID, validation.Required),
		validation.Field(&v.Diffuculty, validation.Required, validation.By(IsDiffucultyLevel)),
		validation.Field(&v.Description, validation.By(IsSQL)),
		validation.Field(&v.Language, validation.Required, validation.Length(2, 100), validation.By(IsSQL)),
		validation.Field(&v.URL, validation.Required, validation.Length(1, 100), validation.By(IsSQL)),
	)
}
