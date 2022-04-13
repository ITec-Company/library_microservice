package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Audio ...
type Audio struct {
	ID            int             `json:"id"`
	SubDirection  DevSubDirection `json:"subDirectionId"`
	Title         string          `json:"title"`
	Diffuculty    DiffucultyLevel `json:"diffuculty"`
	Rating        float32         `json:"rating"`
	Description   string          `json:"description"`
	Language      string          `json:"language"`
	URL           string          `json:"url"`
	DownloadCount uint32          `json:"downloadCount"`
}

// AudioDTO ...
type AudioDTO struct {
	ID             int     `json:"id"`
	SubDirectionID int     `json:"subDirectionId"`
	Title          string  `json:"title"`
	Diffuculty     string  `json:"diffuculty"`
	Rating         float32 `json:"rating"`
	Description    string  `json:"description"`
	Language       string  `json:"language"`
	URL            string  `json:"url"`
	DownloadCount  uint32  `json:"downloadCount"`
}

// Validate ...
func (a *AudioDTO) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.SubDirectionID, validation.Required),
		validation.Field(&a.Title, validation.Required, validation.Length(3, 100), validation.By(IsSQL), validation.By(IsLetterHyphenSpaces)),
		validation.Field(&a.Diffuculty, validation.Required, validation.By(IsDiffucultyLevel)),
		validation.Field(&a.Rating, validation.Required),
		validation.Field(&a.Description, validation.By(IsSQL)),
		validation.Field(&a.Language, validation.Required, validation.Length(2, 100), validation.By(IsSQL)),
		validation.Field(&a.URL, validation.Required, validation.Length(1, 100), validation.By(IsSQL)),
	)
}
