package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Book ...
type Book struct {
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

// BookDTO ...
type BookDTO struct {
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

// DiffucultyLevel ...
type DiffucultyLevel string

// DiffucultyLevel variatons
var (
	JuniorLevel DiffucultyLevel = "junior"
	MiddleLevel DiffucultyLevel = "middle"
	SeniorLevel DiffucultyLevel = "senior"
)

// Validate ...
func (b *BookDTO) Validate() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.AuthorID, validation.Required),
		validation.Field(&b.SubDirectionID, validation.Required),
		validation.Field(&b.Title, validation.Required, validation.Length(3, 100), validation.By(IsSQL), validation.By(IsLetterHyphenSpaces)),
		validation.Field(&b.EditionDate, validation.Required),
		validation.Field(&b.Diffuculty, validation.Required, validation.By(IsDiffucultyLevel)),
		validation.Field(&b.Rating, validation.Required),
		validation.Field(&b.Description, validation.By(IsSQL)),
		validation.Field(&b.Language, validation.Required, validation.Length(2, 100), validation.By(IsSQL)),
		validation.Field(&b.URL, validation.Required, validation.Length(1, 100), validation.By(IsSQL)),
	)
}
