package app

import (
	"github.com/julienschmidt/httprouter"
	"library-go/internal/composite"
)

func ConfigureRouter(router *httprouter.Router, composites composite.Composites) {
	composites.Article.Handler.Register(router)
	composites.Audio.Handler.Register(router)
	composites.Author.Handler.Register(router)
	composites.Book.Handler.Register(router)
	composites.Direction.Handler.Register(router)
	composites.Review.Handler.Register(router)
	composites.Tag.Handler.Register(router)
	composites.Video.Handler.Register(router)

}
