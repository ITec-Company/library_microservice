package composite

import (
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

func (c *Composites) NewPostgres(store postgres.Store) {
	c.Article.New(store.ArticleStorage, store.Logger)
	c.Audio.New(store.AudioStorage, store.Logger)
	c.Author.New(store.AuthorStorage, store.Logger)
	c.Book.New(store.BookStorage, store.Logger)
	c.Direction.New(store.DirectionStorage, store.Logger)
	c.Review.New(store.ReviewStorage, store.Logger)
	c.Tag.New(store.TagStorage, store.Logger)
	c.Video.New(store.VideoStorage, store.Logger)
}
