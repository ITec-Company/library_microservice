package domain

import "time"

// Article

type CreateArticleDTO struct {
	Title         string    `json:"title"`
	DirectionUUID string    `json:"direction_uuid"`
	AuthorUUID    string    `json:"author_uuid"`
	Difficulty    string    `json:"difficulty"`
	EditionDate   time.Time `json:"edition_date"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	Language      string    `json:"language"`
	TagsUUIDs     []string  `json:"tags_uuids"`
}

type UpdateArticleDTO struct {
	UUID          string    `json:"uuid"`
	Title         string    `json:"title,omitempty"`
	DirectionUUID string    `json:"direction_uuid,omitempty"`
	AuthorUUID    string    `json:"author_uuid,omitempty"`
	Difficulty    string    `json:"difficulty,omitempty"`
	EditionDate   time.Time `json:"edition_date,omitempty"`
	Rating        float32   `json:"rating,omitempty"`
	Description   string    `json:"description,omitempty"`
	URL           string    `json:"url,omitempty"`
	Language      string    `json:"language,omitempty"`
	TagsUUIDs     []string  `json:"tags_uuids,omitempty"`
	DownloadCount uint32    `json:"download_count,omitempty"`
}

// Audio

type CreateAudioDTO struct {
	Title         string   `json:"title"`
	DirectionUUID string   `json:"direction_uuid"`
	Language      string   `json:"language"`
	URL           string   `json:"url"`
	Difficulty    string   `json:"difficulty"`
	TagsUUIDs     []string `json:"tags_uuids"`
}

type UpdateAudioDTO struct {
	UUID          string   `json:"uuid"`
	Title         string   `json:"title,omitempty"`
	DirectionUUID string   `json:"direction_uuid,omitempty"`
	Difficulty    string   `json:"difficulty,omitempty"`
	Rating        float32  `json:"rating,omitempty"`
	URL           string   `json:"url,omitempty"`
	Language      string   `json:"language,omitempty"`
	TagsUUIDs     []string `json:"tags_uuids,omitempty"`
	DownloadCount uint32   `json:"download_count,omitempty"`
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
	Title         string    `json:"title"`
	DirectionUUID string    `json:"direction_uuid"`
	AuthorUUID    string    `json:"author_uuid"`
	Difficulty    string    `json:"difficulty"`
	EditionDate   time.Time `json:"edition_date"`
	Description   string    `json:"description"`
	URL           string    `json:"url"`
	Language      string    `json:"language"`
	TagsUUIDs     []string  `json:"tags_uuids"`
}

type UpdateBookDTO struct {
	UUID          string    `json:"uuid"`
	Title         string    `json:"title,omitempty"`
	DirectionUUID string    `json:"direction_uuid,omitempty"`
	AuthorUUID    string    `json:"author_uuid,omitempty"`
	Difficulty    string    `json:"difficulty,omitempty"`
	EditionDate   time.Time `json:"edition_date,omitempty"`
	Rating        float32   `json:"rating,omitempty"`
	Description   string    `json:"description,omitempty"`
	URL           string    `json:"url,omitempty"`
	Language      string    `json:"language,omitempty"`
	TagsUUIDs     []string  `json:"tags_uuids,omitempty"`
	DownloadCount uint32    `json:"download_count,omitempty"`
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
	URL           string   `json:"url"`
	Language      string   `json:"language"`
	TagsUUIDs     []string `json:"tags_uuids"`
}

type UpdateVideoDTO struct {
	UUID          string   `json:"uuid"`
	DirectionUUID string   `json:"direction_uuid,omitempty"`
	Title         string   `json:"title,omitempty"`
	Difficulty    string   `json:"difficulty,omitempty"`
	Rating        float32  `json:"rating,omitempty"`
	URL           string   `json:"url,omitempty"`
	Language      string   `json:"language,omitempty"`
	TagsUUIDs     []string `json:"tags_uuids,omitempty"`
	DownloadCount uint32   `json:"download_count,omitempty"`
}
