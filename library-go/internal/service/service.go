package service

import (
	"library-go/internal/store/postgres"
	"library-go/pkg/logging"
)

type Service struct {
	Logger    *logging.Logger
	Store     *postgres.Store
	Article   ArticleService
	Audio     AudioService
	Author    AuthorService
	Book      BookService
	Review    ReviewService
	Tag       TagService
	Video     VideoService
	Direction DirectionService
}

func New(store *postgres.Store) *Service {
	return &Service{
		Logger:    store.Logger,
		Store:     store,
		Article:   NewArticleService(&store.ArticleStorage, store.Logger),
		Audio:     NewAudioService(&store.AudioStorage, store.Logger),
		Author:    NewAuthorService(&store.AuthorStorage, store.Logger),
		Book:      NewBookService(&store.BookStorage, store.Logger),
		Review:    NewReviewService(&store.ReviewStorage, store.Logger),
		Tag:       NewTagService(&store.TagStorage, store.Logger),
		Video:     NewVideoService(&store.VideoStorage, store.Logger),
		Direction: NewDirectionService(&store.DirectionStorage, store.Logger),
	}
}
