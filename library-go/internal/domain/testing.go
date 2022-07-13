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
		Name: "test tag",
	}
}

func TestTagCreateDTO() *CreateTagDTO {
	return &CreateTagDTO{
		Name: "test tag",
	}
}

func TestTagUpdateDTO() *UpdateTagDTO {
	return &UpdateTagDTO{
		UUID: "1",
		Name: "test tag",
	}
}

func TestArticle() *Article {
	return &Article{
		UUID:          "1",
		Title:         "Test Title",
		Direction:     *TestDirection(),
		Difficulty:    "Test Difficulty",
		Text:          "Test Text",
		Author:        *TestAuthor(),
		EditionDate:   2000,
		Rating:        5.0,
		Description:   "Test Description",
		LocalURL:      "Test LocalURL",
		WebURL:        "Test WebURL",
		ImageURL:      "Test ImageURL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
		CreatedAt:     time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
	}
}

func TestArticleCreateDTO() *CreateArticleDTO {
	return &CreateArticleDTO{
		Title:         "Test_Title",
		DirectionUUID: "1",
		AuthorUUID:    "1",
		Difficulty:    "Test Difficulty",
		EditionDate:   2000,
		Description:   "Test Description",
		Text:          "text",
		LocalURL:      "/articles/|split|/author(1)-title(Test Title).docx",
		WebURL:        "Test URL",
		ImageURL:      "/articles/|split|/original.jpg",
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
		EditionDate:   2000,
		Rating:        5.5,
		Description:   "Test Description",
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
		EditionDate:   2000,
		Rating:        5.0,
		Description:   "Test Description",
		LocalURL:      "Test LocalURL",
		ImageURL:      "Test ImageURL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
		CreatedAt:     time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
	}
}

func TestBookCreateDTO() *CreateBookDTO {
	return &CreateBookDTO{
		Title:         "Test_Title",
		DirectionUUID: "1",
		AuthorUUID:    "1",
		Difficulty:    "Test Difficulty",
		EditionDate:   2000,
		Description:   "Test Description",
		LocalURL:      "/books/|split|/author(1)-title(Test_Title).pdf",
		ImageURL:      "/books/|split|/original.jpg",
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
		EditionDate:   2000,
		Description:   "Test Description",
		Language:      "Test Language",
		TagsUUIDs:     []string{"1"},
	}
}

func TestAudio() *Audio {
	return &Audio{
		UUID:          "1",
		Title:         "Test Title",
		Direction:     *TestDirection(),
		Difficulty:    "Test Difficulty",
		Rating:        5.0,
		LocalURL:      "Test URL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
		CreatedAt:     time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
	}
}

func TestAudioCreateDTO() *CreateAudioDTO {
	return &CreateAudioDTO{
		Title:         "Test_Title",
		Difficulty:    "Test Difficulty",
		LocalURL:      "/audios/|split|/title(Test_Title).mp3",
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
		Description:   "Test Description",
		LocalURL:      "Test LocalURL",
		WebURL:        "Test WebURL",
		Language:      "Test Language",
		Tags:          []Tag{*TestTag()},
		DownloadCount: 10,
		CreatedAt:     time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
	}
}

func TestVideoCreateDTO() *CreateVideoDTO {
	return &CreateVideoDTO{
		Title:         "Test_Title",
		DirectionUUID: "1",
		Difficulty:    "Test Difficulty",
		Description:   "Test Description",
		LocalURL:      "/videos/|split|/title(Test_Title).mp4",
		WebURL:        "text",
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
		Description:   "Test Description",
		WebURL:        "Test URL",
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
		CreatedAt:      time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC),
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
