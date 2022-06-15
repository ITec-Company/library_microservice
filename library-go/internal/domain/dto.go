package domain

import (
	"image"
	"io"
)

// Article

type CreateArticleDTO struct {
	Title         string   `json:"title"`
	DirectionUUID string   `json:"direction_uuid"`
	AuthorUUID    string   `json:"author_uuid"`
	Difficulty    string   `json:"difficulty"`
	EditionDate   uint     `json:"edition_date"`
	Description   string   `json:"description"`
	Text          string   `json:"text"`
	LocalURL      string   `json:"local_url"`
	WebURL        string   `json:"web_url,omitempty"`
	Language      string   `json:"language"`
	TagsUUIDs     []string `json:"tags_uuids"`
	ImageURL      string   `json:"image_url,omitempty"`
}

type UpdateArticleDTO struct {
	UUID          string   `json:"uuid"`
	Title         string   `json:"title,omitempty"`
	DirectionUUID string   `json:"direction_uuid,omitempty"`
	AuthorUUID    string   `json:"author_uuid,omitempty"`
	Difficulty    string   `json:"difficulty,omitempty"`
	EditionDate   uint     `json:"edition_date,omitempty"`
	Rating        float32  `json:"rating,omitempty"`
	Description   string   `json:"description,omitempty"`
	Text          string   `json:"text,omitempty"`
	LocalURL      string   `json:"local_url,omitempty"`
	WebURL        string   `json:"web_url,omitempty"`
	Language      string   `json:"language,omitempty"`
	TagsUUIDs     []string `json:"tags_uuids,omitempty"`
	DownloadCount uint32   `json:"download_count,omitempty"`
	ImageURL      string   `json:"image_url,omitempty"`
}

type UpdateArticleFileDTO struct {
	UUID        string    `json:"uuid"`
	NewFileName string    `json:"new-file-name"`
	OldFileName string    `json:"old-file-name"`
	File        io.Reader `json:"file"`
	LocalURL    string    `json:"local-url"`
	LocalPath   string    `json:"local-path"`
}

type UpdateArticleImageDTO struct {
	UUID      string      `json:"uuid"`
	Image     image.Image `json:"image"`
	LocalPath string      `json:"local-path"`
}

// Audio

type CreateAudioDTO struct {
	Title         string   `json:"title"`
	DirectionUUID string   `json:"direction_uuid"`
	Language      string   `json:"language"`
	LocalURL      string   `json:"local_url"`
	Difficulty    string   `json:"difficulty"`
	TagsUUIDs     []string `json:"tags_uuids"`
}

type UpdateAudioDTO struct {
	UUID          string   `json:"uuid"`
	Title         string   `json:"title,omitempty"`
	DirectionUUID string   `json:"direction_uuid,omitempty"`
	Difficulty    string   `json:"difficulty,omitempty"`
	Rating        float32  `json:"rating,omitempty"`
	LocalURL      string   `json:"local_url,omitempty"`
	Language      string   `json:"language,omitempty"`
	TagsUUIDs     []string `json:"tags_uuids,omitempty"`
	DownloadCount uint32   `json:"download_count,omitempty"`
}

type UpdateAudioFileDTO struct {
	UUID        string    `json:"uuid"`
	NewFileName string    `json:"new-file-name"`
	OldFileName string    `json:"old-file-name"`
	File        io.Reader `json:"file"`
	LocalURL    string    `json:"local-url"`
	LocalPath   string    `json:"local-path"`
}

// Author

type CreateAuthorDTO struct {
	FullName string `json:"full_name"`
}

type UpdateAuthorDTO struct {
	UUID     string `json:"uuid"`
	FullName string `json:"full_name,omitempty"`
}

// Book

type CreateBookDTO struct {
	Title         string   `json:"title"`
	DirectionUUID string   `json:"direction_uuid"`
	AuthorUUID    string   `json:"author_uuid"`
	Difficulty    string   `json:"difficulty"`
	EditionDate   uint     `json:"edition_date"`
	Description   string   `json:"description"`
	LocalURL      string   `json:"local_url"`
	Language      string   `json:"language"`
	TagsUUIDs     []string `json:"tags_uuids"`
	ImageURL      string   `json:"image_url,omitempty"`
}

type UpdateBookDTO struct {
	UUID          string   `json:"uuid"`
	Title         string   `json:"title,omitempty"`
	DirectionUUID string   `json:"direction_uuid,omitempty"`
	AuthorUUID    string   `json:"author_uuid,omitempty"`
	Difficulty    string   `json:"difficulty,omitempty"`
	EditionDate   uint     `json:"edition_date,omitempty"`
	Rating        float32  `json:"rating,omitempty"`
	Description   string   `json:"description,omitempty"`
	LocalURL      string   `json:"local_url,omitempty"`
	Language      string   `json:"language,omitempty"`
	TagsUUIDs     []string `json:"tags_uuids,omitempty"`
	DownloadCount uint32   `json:"download_count,omitempty"`
	ImageURL      string   `json:"image_url,omitempty"`
}

type UpdateBookFileDTO struct {
	UUID        string    `json:"uuid"`
	NewFileName string    `json:"new-file-name"`
	OldFileName string    `json:"old-file-name"`
	File        io.Reader `json:"file"`
	LocalURL    string    `json:"local-url"`
	LocalPath   string    `json:"local-path"`
}

type UpdateBookImageDTO struct {
	UUID      string      `json:"uuid"`
	Image     image.Image `json:"image"`
	LocalPath string      `json:"local-path"`
}

// Direction

type CreateDirectionDTO struct {
	Name string `json:"name"`
}

type UpdateDirectionDTO struct {
	UUID string `json:"uuid"`
	Name string `json:"name,omitempty"`
}

// Review

type CreateReviewDTO struct {
	Text           string `json:"text"`
	FullName       string `json:"full_name"`
	LiteratureUUID string `json:"literature_uuid"`
	Source         string `json:"source"`
}

type UpdateReviewDTO struct {
	UUID     string  `json:"uuid"`
	FullName string  `json:"full_name,omitempty"`
	Text     string  `json:"text,omitempty"`
	Rating   float32 `json:"rating,omitempty"`
}

// Tag

type CreateTagDTO struct {
	Name string `json:"name"`
}

type UpdateTagDTO struct {
	UUID string `json:"uuid"`
	Name string `json:"name,omitempty"`
}

// Video

type CreateVideoDTO struct {
	Title         string   `json:"title"`
	DirectionUUID string   `json:"direction_uuid"`
	Difficulty    string   `json:"difficulty"`
	LocalURL      string   `json:"local_url"`
	WebURL        string   `json:"web_url,omitempty"`
	Language      string   `json:"language"`
	TagsUUIDs     []string `json:"tags_uuids"`
}

type UpdateVideoDTO struct {
	UUID          string   `json:"uuid"`
	DirectionUUID string   `json:"direction_uuid,omitempty"`
	Title         string   `json:"title,omitempty"`
	Difficulty    string   `json:"difficulty,omitempty"`
	Rating        float32  `json:"rating,omitempty"`
	LocalURL      string   `json:"local_url,omitempty"`
	WebURL        string   `json:"web_url,omitempty"`
	Language      string   `json:"language,omitempty"`
	TagsUUIDs     []string `json:"tags_uuids,omitempty"`
	DownloadCount uint32   `json:"download_count,omitempty"`
}

type UpdateVideoFileDTO struct {
	UUID        string    `json:"uuid"`
	NewFileName string    `json:"new-file-name"`
	OldFileName string    `json:"old-file-name"`
	File        io.Reader `json:"file"`
	LocalURL    string    `json:"local-url"`
	LocalPath   string    `json:"local-path"`
}
