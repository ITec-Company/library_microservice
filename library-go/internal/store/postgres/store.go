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

func New(db *sql.DB, logger *logging.Logger) *Store {
	return &Store{
		Logger:           logger,
		DB:               db,
		ArticleStorage:   NewArticleStorage(db, logger),
		AudioStorage:     NewAudioStorage(db, logger),
		AuthorStorage:    NewAuthorStorage(db, logger),
		BookStorage:      NewBookStorage(db, logger),
		ReviewStorage:    NewReviewStorage(db, logger),
		TagStorage:       NewTagStorage(db, logger),
		VideoStorage:     NewVideoStorage(db, logger),
		DirectionStorage: NewDirectionStorage(db, logger),
	}
}
