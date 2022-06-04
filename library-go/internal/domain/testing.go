package domain

import "time"

func TestAuthor() *Author {
	return &Author{
		UUID:     "1",
		FullName: "Test Author",
	}
}

func TestAuthorCreateDTO() *CreateAuthorDTO {
	return &CreateAuthorDTO{
		FullName: "Test Author",
	}
}

func TestAuthorUpdateDTO() *UpdateAuthorDTO {
	return &UpdateAuthorDTO{
		UUID:     "1",
		FullName: "Test Author",
	}
}

func TestDirection() *Direction {
	return &Direction{
		UUID: "1",
		Name: "Test Direction",
	}
}

func TestDirectionCreateDTO() *CreateDirectionDTO {
	return &CreateDirectionDTO{
		Name: "Test Direction",
	}
}

func TestDirectionUpdateDTO() *UpdateDirectionDTO {
	return &UpdateDirectionDTO{
		UUID: "1",
		Name: "Test Direction",
	}
}

func TestTag() *Tag {
	return &Tag{
		UUID: "1",
		Name: "Test Tag",
	}
}

func TestTagCreateDTO() *CreateTagDTO {
	return &CreateTagDTO{
		Name: "Test Tag",
	}
}

func TestTagUpdateDTO() *UpdateTagDTO {
	return &UpdateTagDTO{
		UUID: "1",
		Name: "Test Tag",
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
		LocalURL:      "Test LocalURL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
		ImageURL:      "imageURL",
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
		LocalURL:      "Test LocalURL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		ImageURL:      "imageURL",
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
		LocalURL:      "Test LocalURL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		DownloadCount: 10,
		ImageURL:      "imageURL",
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
		LocalURL:      "Test LocalURL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
		ImageURL:      "imageURL",
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
		LocalURL:      "Test LocalURL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		ImageURL:      "imageURL",
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
		LocalURL:      "Test LocalURL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		DownloadCount: 10,
		ImageURL:      "imageURL",
	}
}

func TestAudio() *Audio {
	return &Audio{
		UUID:          "1",
		Title:         "Test Title",
		Direction:     *TestDirection(),
		Difficulty:    "Test Difficulty",
		Rating:        5.0,
		LocalURL:      "Test LocalURL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
	}
}

func TestAudioCreateDTO() *CreateAudioDTO {
	return &CreateAudioDTO{
		Title:         "Test Title",
		Difficulty:    "Test Difficulty",
		LocalURL:      "Test LocalURL",
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
		LocalURL:      "Test LocalURL",
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
		LocalURL:      "Test LocalURL",
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
		LocalURL:      "Test LocalURL",
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
		LocalURL:      "Test LocalURL",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
		DownloadCount: 10,
	}
}

func TestReview() *Review {
	return &Review{
		UUID:           "1",
		FullName:       "Test Review",
		Text:           "Test Text",
		Rating:         5.5,
		Date:           time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
		Source:         "Test Source",
		LiteratureUUID: "1",
	}
}

func TestReviewCreateDTO() *CreateReviewDTO {
	return &CreateReviewDTO{
		Text:           "Test Text",
		FullName:       "Test Review",
		LiteratureUUID: "1",
		Source:         "Test Source",
	}
}

func TestReviewUpdateDTO() *UpdateReviewDTO {
	return &UpdateReviewDTO{
		UUID:     "1",
		FullName: "Test Review",
		Text:     "Test Text",
		Rating:   5.5,
	}
}
