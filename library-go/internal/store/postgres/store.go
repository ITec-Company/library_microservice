package postgres

import (
	"database/sql"
	"errors"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type Store struct {
	Logger           *logging.Logger
	DB               *sql.DB
	ArticleStorage   store.ArticleStorage
	AudioStorage     store.AudioStorage
	AuthorStorage    store.AuthorStorage
	BookStorage      store.BookStorage
	ReviewStorage    store.ReviewStorage
	TagStorage       store.TagStorage
	VideoStorage     store.VideoStorage
	DirectionStorage store.DirectionStorage
}

var (
	// ErrNoRowsAffected ...
	ErrNoRowsAffected = errors.New("no rows affected")
)

func (s *Store) NewDB(db *sql.DB, logger *logging.Logger) {
	s.Logger = logger
	s.DB = db
	s.ArticleStorage = NewArticleStorage(s.DB, s.Logger)
	s.AudioStorage = NewAudioStorage(s.DB, s.Logger)
	s.AuthorStorage = NewAuthorStorage(s.DB, s.Logger)
	s.BookStorage = NewBookStorage(s.DB, s.Logger)
	s.ReviewStorage = NewReviewStorage(s.DB, s.Logger)
	s.TagStorage = NewTagStorage(s.DB, s.Logger)
	s.VideoStorage = NewVideoStorage(s.DB, s.Logger)
	s.DirectionStorage = NewDirectionStorage(s.DB, s.Logger)
}
