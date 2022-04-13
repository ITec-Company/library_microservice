package model

import "time"

// TestReview ...
func TestReview() *Review {
	return &Review{
		ID:           1,
		UserID:       1,
		LiteratureID: 1,
		Ratig:        5,
		Text:         "simple review",
	}
}

// TestUser ...
func TestUser() *User {
	return &User{
		ID:       1,
		Email:    "example@mail.org",
		FullName: "Ivanov Ivan Ivanovich",
	}
}

// TestAuthor ...
func TestAuthor() *Author {
	return &Author{
		ID:       1,
		FullName: "Ivanov Ivan Ivanovich",
	}
}

// TestDevDirection ...
func TestDevDirection() *DevDirection {
	return &DevDirection{
		ID:        1,
		Direction: Backend,
	}
}

// TestDevDirectionDTO ...
func TestDevDirectionDTO() *DevDirectionDTO {
	return &DevDirectionDTO{
		ID:        1,
		Direction: string(Backend),
	}
}

// TestDevSubDirection ...
func TestDevSubDirection() *DevSubDirection {
	return &DevSubDirection{
		ID:           1,
		SubDirection: Golang,
	}
}

// TestDevSubDirectionDTO ...
func TestDevSubDirectionDTO() *DevSubDirectionDTO {
	return &DevSubDirectionDTO{
		ID:           1,
		SubDirection: string(Golang),
	}
}

// TestTag ...
func TestTag() *Tag {
	return &Tag{
		ID:   1,
		Name: "golang",
	}
}

// TestTagDTO ...
func TestTagDTO() *TagDTO {
	return &TagDTO{
		ID:   1,
		Name: "golang",
	}
}

// TestBook ...
func TestBook() *Book {
	return &Book{
		ID:            1,
		Author:        *TestAuthor(),
		SubDirection:  *TestDevSubDirection(),
		Title:         "Golang for beginners",
		EditionDate:   time.Date(2015, 05, 05, 0, 0, 0, 0, time.UTC),
		Diffuculty:    JuniorLevel,
		Rating:        8.2,
		Description:   "Description",
		Language:      "eng",
		URL:           "URL...",
		DownloadCount: 829,
	}
}

// TestBookDTO ...
func TestBookDTO() *BookDTO {
	return &BookDTO{
		ID:             1,
		AuthorID:       1,
		SubDirectionID: 1,
		Title:          "Golang for beginners",
		EditionDate:    time.Date(2015, 05, 05, 0, 0, 0, 0, time.UTC),
		Diffuculty:     string(JuniorLevel),
		Rating:         8.2,
		Description:    "Description",
		Language:       "eng",
		URL:            "URL...",
		DownloadCount:  829,
	}
}

// TestArticle ...
func TestArticle() *Article {
	return &Article{
		ID:            1,
		Author:        *TestAuthor(),
		SubDirection:  *TestDevSubDirection(),
		Title:         "Golang for beginners artcile",
		EditionDate:   time.Date(2015, 05, 05, 0, 0, 0, 0, time.UTC),
		Diffuculty:    JuniorLevel,
		Rating:        8.2,
		Description:   "Description",
		Language:      "eng",
		URL:           "URL...",
		DownloadCount: 829,
	}
}

// TestArticleDTO ...
func TestArticleDTO() *ArticleDTO {
	return &ArticleDTO{
		ID:             1,
		AuthorID:       1,
		SubDirectionID: 1,
		Title:          "Golang for beginners artcile",
		EditionDate:    time.Date(2015, 05, 05, 0, 0, 0, 0, time.UTC),
		Diffuculty:     string(JuniorLevel),
		Rating:         8.2,
		Description:    "Description",
		Language:       "eng",
		URL:            "URL...",
		DownloadCount:  829,
	}
}

// TestVideo ...
func TestVideo() *Video {
	return &Video{
		ID:           1,
		SubDirection: *TestDevSubDirection(),
		Diffuculty:   JuniorLevel,
		Description:  "Description",
		Language:     "eng",
		URL:          "URL...",
	}
}

// TestVideoDTO ...
func TestVideoDTO() *VideoDTO {
	return &VideoDTO{
		ID:             1,
		SubDirectionID: 1,
		Diffuculty:     string(JuniorLevel),
		Description:    "Description Video",
		Language:       "eng",
		URL:            "URL...",
	}
}

// TestAudio ...
func TestAudio() *Audio {
	return &Audio{
		ID:            1,
		SubDirection:  *TestDevSubDirection(),
		Title:         "Golang for beginners artcile",
		Diffuculty:    JuniorLevel,
		Rating:        8.2,
		Description:   "Description",
		Language:      "eng",
		URL:           "URL...",
		DownloadCount: 829,
	}
}

// TestAudioDTO ...
func TestAudioDTO() *AudioDTO {
	return &AudioDTO{
		ID:             1,
		SubDirectionID: 1,
		Title:          "Golang for beginners artcile",
		Diffuculty:     string(JuniorLevel),
		Rating:         8.2,
		Description:    "Description",
		Language:       "eng",
		URL:            "URL...",
		DownloadCount:  829,
	}
}
