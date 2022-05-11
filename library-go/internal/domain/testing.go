package domain

import "time"

func TestAuthor() *Author {
	return &Author{
		UUID:     "1",
		FullName: "test Author",
	}
}

func TestAuthorCreateDTO() *CreateAuthorDTO {
	return &CreateAuthorDTO{
		FullName: "test Author",
	}
}

func TestAuthorUpdateDTO() *UpdateAuthorDTO {
	return &UpdateAuthorDTO{
		UUID:     "1",
		FullName: "test Author",
	}
}

func TestDirection() *Direction {
	return &Direction{
		UUID: "1",
		Name: "test Direction",
	}
}

func TestDirectionCreateDTO() *CreateDirectionDTO {
	return &CreateDirectionDTO{
		Name: "test Direction",
	}
}

func TestDirectionUpdateDTO() *UpdateDirectionDTO {
	return &UpdateDirectionDTO{
		UUID: "1",
		Name: "test Direction",
	}
}

func TestTag() *Tag {
	return &Tag{
		UUID: "1",
		Name: "test Tag",
	}
}

func TestTagCreateDTO() *CreateTagDTO {
	return &CreateTagDTO{
		Name: "test Tag",
	}
}

func TestTagUpdateDTO() *UpdateTagDTO {
	return &UpdateTagDTO{
		UUID: "1",
		Name: "test Tag",
	}
}

func TestArticle() *Article {
	return &Article{
		UUID:          "1",
		Title:         "Test Title",
		Direction:     *TestDirection(),
		Difficulty:    "Test Difficulty",
		Author:        *TestAuthor(),
		EditionDate:   time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		Rating:        5.0,
		Description:   "Test Description",
		URL:           "Test URL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
	}
}

func TestArticleCreateDTO() *CreateArticleDTO {
	return &CreateArticleDTO{
		Title:         "Test Title",
		DirectionUUID: "1",
		AuthorUUID:    "1",
		Difficulty:    "Test Difficulty",
		EditionDate:   time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		Description:   "Test Description",
		URL:           "Test URL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
	}
}

func TestArticleUpdateDTO() *UpdateArticleDTO {
	return &UpdateArticleDTO{
		UUID:          "1",
		Title:         "Test Title",
		DirectionUUID: "1",
		AuthorUUID:    "1",
		Difficulty:    "Test Difficulty",
		EditionDate:   time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		Rating:        5.5,
		Description:   "Test Description",
		URL:           "Test URL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		DownloadCount: 10,
	}
}

func TestBook() *Book {
	return &Book{
		UUID:          "1",
		Title:         "Test Title",
		Direction:     *TestDirection(),
		Difficulty:    "Test Difficulty",
		Author:        *TestAuthor(),
		EditionDate:   time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		Rating:        5.0,
		Description:   "Test Description",
		URL:           "Test URL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
	}
}

func TestBookCreateDTO() *CreateBookDTO {
	return &CreateBookDTO{
		Title:         "Test Title",
		DirectionUUID: "1",
		AuthorUUID:    "1",
		Difficulty:    "Test Difficulty",
		EditionDate:   time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		Description:   "Test Description",
		URL:           "Test URL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
	}
}

func TestBookUpdateDTO() *UpdateBookDTO {
	return &UpdateBookDTO{
		UUID:          "1",
		Title:         "Test Title",
		DirectionUUID: "1",
		AuthorUUID:    "1",
		Difficulty:    "Test Difficulty",
		EditionDate:   time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		Rating:        5.5,
		Description:   "Test Description",
		URL:           "Test URL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		DownloadCount: 10,
	}
}

func TestAudio() *Audio {
	return &Audio{
		UUID:          "1",
		Title:         "Test Title",
		Direction:     *TestDirection(),
		Difficulty:    "Test Difficulty",
		Rating:        5.0,
		URL:           "Test URL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
	}
}

func TestAudioCreateDTO() *CreateAudioDTO {
	return &CreateAudioDTO{
		Title:         "Test Title",
		Difficulty:    "Test Difficulty",
		URL:           "Test URL",
		Language:      "Test Language",
		DirectionUUID: "1",
		TagsUUIDs:     []string{"1"},
	}
}

func TestAudioUpdateDTO() *UpdateAudioDTO {
	return &UpdateAudioDTO{
		UUID:          "1",
		Title:         "Test Title",
		DirectionUUID: "1",
		Difficulty:    "Test Difficulty",
		Rating:        5.5,
		URL:           "Test URL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		DownloadCount: 10,
	}
}

func TestVideo() *Video {
	return &Video{
		UUID:          "1",
		Title:         "Test Title",
		Direction:     *TestDirection(),
		Rating:        5.0,
		Difficulty:    "Test Difficulty",
		URL:           "Test URL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
	}
}

func TestVideoCreateDTO() *CreateVideoDTO {
	return &CreateVideoDTO{
		Title:         "Test Title",
		DirectionUUID: "1",
		Difficulty:    "Test Difficulty",
		URL:           "Test URL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
	}
}

func TestVideoUpdateDTO() *UpdateVideoDTO {
	return &UpdateVideoDTO{
		UUID:          "1",
		DirectionUUID: "1",
		Title:         "Test Title",
		Difficulty:    "Test Difficulty",
		Rating:        5.5,
		URL:           "Test URL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		DownloadCount: 10,
	}
}

func TestReview() *Review {
	return &Review{
		UUID:           "1",
		FullName:       "test Review",
		Text:           "test Text",
		Rating:         5.5,
		Date:           time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		Source:         "test Source",
		LiteratureUUID: "1",
	}
}

func TestReviewCreateDTO() *CreateReviewDTO {
	return &CreateReviewDTO{
		Text:           "test Text",
		FullName:       "test Review",
		LiteratureUUID: "1",
		Source:         "test Source",
	}
}

func TestReviewUpdateDTO() *UpdateReviewDTO {
	return &UpdateReviewDTO{
		UUID:     "1",
		FullName: "test Review",
		Text:     "test Text",
		Rating:   5.5,
	}
}
