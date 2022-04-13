package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Article ...
type Article struct {
	ID            int             `json:"id"`
	Author        Author          `json:"authorId"`
	SubDirection  DevSubDirection `json:"subDirectionId"`
	Title         string          `json:"title"`
	EditionDate   time.Time       `json:"editionDate"`
	Diffuculty    DiffucultyLevel `json:"diffuculty"`
	Rating        float32         `json:"rating"`
	Description   string          `json:"description"`
	Language      string          `json:"language"`
	URL           string          `json:"url"`
	DownloadCount uint32          `json:"downloadCount"`
}

// ArticleDTO ...
type ArticleDTO struct {
	ID             int       `json:"id"`
	AuthorID       int       `json:"authorId"`
	SubDirectionID int       `json:"subDirectionId"`
	Title          string    `json:"title"`
	EditionDate    time.Time `json:"editionDate"`
	Diffuculty     string    `json:"diffuculty"`
	Rating         float32   `json:"rating"`
	Description    string    `json:"description"`
	Language       string    `json:"language"`
	URL            string    `json:"url"`
	DownloadCount  uint32    `json:"downloadCount"`
}

// Validate ...
func (a *ArticleDTO) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.AuthorID, validation.Required),
		validation.Field(&a.SubDirectionID, validation.Required),
		validation.Field(&a.Title, validation.Required, validation.Length(3, 100), validation.By(IsSQL), validation.By(IsLetterHyphenSpaces)),
		validation.Field(&a.EditionDate, validation.Required),
		validation.Field(&a.Diffuculty, validation.Required, validation.By(IsDiffucultyLevel)),
		validation.Field(&a.Rating, validation.Required),
		validation.Field(&a.Description, validation.By(IsSQL)),
		validation.Field(&a.Language, validation.Required, validation.Length(2, 100), validation.By(IsSQL)),
		validation.Field(&a.URL, validation.Required, validation.Length(1, 100), validation.By(IsSQL)),
	)
}
