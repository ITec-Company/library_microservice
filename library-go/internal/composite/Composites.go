package composite

import (
	"library-go/internal/handler/http"
	"library-go/internal/store/postgres"
)

type Composites struct {
	Article   ArticleComposite
	Audio     AudioComposite
	Author    AuthorComposite
	Book      BookComposite
	Direction DirectionComposite
	Review    ReviewComposite
	Tag       TagComposite
	Video     VideoComposite
}

func (c *Composites) NewPostgres(store postgres.Store, middleware *http.Middleware) {
	c.Article.New(store.ArticleStorage, store.Logger, middleware)
	c.Audio.New(store.AudioStorage, store.Logger, middleware)
	c.Author.New(store.AuthorStorage, store.Logger, middleware)
	c.Book.New(store.BookStorage, store.Logger, middleware)
	c.Direction.New(store.DirectionStorage, store.Logger, middleware)
	c.Review.New(store.ReviewStorage, store.Logger, middleware)
	c.Tag.New(store.TagStorage, store.Logger, middleware)
	c.Video.New(store.VideoStorage, store.Logger, middleware)
}
