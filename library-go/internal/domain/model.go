package domain

import (
	"time"
)

type Difficulty string
type Source string
type Order string

const (
	Junior Difficulty = "junior"
	Middle Difficulty = "middle"
	Senior Difficulty = "senior"

	BookSrc    Source = "book"
	ArticleSrc Source = "article"
	VideoSrc   Source = "video"
	AudioSrc   Source = "audio"

	OrderASC  Order = "asc"
	OrderDESC Order = "desc"
)

type Article struct {
	UUID          string     `json:"uuid,omitempty"`
	Title         string     `json:"title,omitempty"`
	Direction     Direction  `json:"direction,omitempty"`
	Difficulty    Difficulty `json:"difficulty,omitempty"`
	Author        Author     `json:"author,omitempty"`
	EditionDate   time.Time  `json:"edition_date"`
	Rating        float32    `json:"rating,omitempty"`
	Description   string     `json:"description,omitempty"`
	LocalURL      string     `json:"local_url,omitempty"`
	ImageURL      string     `json:"image_url,omitempty"`
	WebURL        string     `json:"web_url,omitempty"`
	Language      string     `json:"language,omitempty"`
	Tags          []Tag      `json:"tags,omitempty"`
	DownloadCount uint32     `json:"download_count,omitempty"`
}

type Audio struct {
	UUID          string     `json:"uuid,omitempty"`
	Title         string     `json:"title,omitempty"`
	Direction     Direction  `json:"direction,omitempty"`
	Difficulty    Difficulty `json:"difficulty,omitempty"`
	Rating        float32    `json:"rating,omitempty"`
	LocalURL      string     `json:"local_url,omitempty"`
	Language      string     `json:"language,omitempty"`
	Tags          []Tag      `json:"tags,omitempty"`
	DownloadCount uint32     `json:"download_count,omitempty"`
}

type Author struct {
	UUID     string `json:"uuid,omitempty"`
	FullName string `json:"full_name,omitempty"`
}

type Book struct {
	UUID          string     `json:"uuid,omitempty"`
	Title         string     `json:"title,omitempty"`
	Direction     Direction  `json:"direction,omitempty"`
	Author        Author     `json:"author,omitempty"`
	Difficulty    Difficulty `json:"difficulty,omitempty"`
	EditionDate   time.Time  `json:"edition_date"`
	Rating        float32    `json:"rating,omitempty"`
	Description   string     `json:"description,omitempty"`
	LocalURL      string     `json:"local_url,omitempty"`
	Language      string     `json:"language,omitempty"`
	Tags          []Tag      `json:"tags,omitempty"`
	DownloadCount uint32     `json:"download_count,omitempty"`
	ImageURL      string     `json:"image_url,omitempty"`
}

type Direction struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

type Review struct {
	UUID           string    `json:"uuid,omitempty"`
	FullName       string    `json:"full_name,omitempty"`
	Text           string    `json:"text,omitempty"`
	Rating         float32   `json:"rating,omitempty"`
	Date           time.Time `json:"date"`
	Source         Source    `json:"source,omitempty"`
	LiteratureUUID string    `json:"literature_uuid,omitempty"`
}

type Tag struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

type Video struct {
	UUID          string     `json:"uuid,omitempty"`
	Title         string     `json:"title,omitempty"`
	Direction     Direction  `json:"direction,omitempty"`
	Rating        float32    `json:"rating,omitempty"`
	Difficulty    Difficulty `json:"difficulty,omitempty"`
	LocalURL      string     `json:"local_url,omitempty"`
	WebURL        string     `json:"web_url,omitempty"`
	Language      string     `json:"language,omitempty"`
	Tags          []Tag      `json:"tags,omitempty"`
	DownloadCount uint32     `json:"download_count,omitempty"`
}

type SortFilterPagination struct {
	SortBy         string                 `json:"sort_by,omitempty"`
	Order          Order                  `json:"order,omitempty"`
	FiltersAndArgs map[string]interface{} `json:"filters_and_args,omitempty"`
	Limit          uint64                 `json:"limit"`
	Page           uint64                 `json:"page"`
}
