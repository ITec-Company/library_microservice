package http

import (
	"github.com/julienschmidt/httprouter"
	"library-go/internal/service"
	"library-go/pkg/logging"
)

type Router struct {
	Logger     *logging.Logger
	Service    *service.Service
	Router     *httprouter.Router
	Middleware Middleware
	Article    ArticleHandler
	Audio      AudioHandler
	Author     AuthorHandler
	Book       BookHandler
	Review     ReviewHandler
	Tag        TagHandler
	Video      VideoHandler
	Direction  DirectionHandler
}

func New(service *service.Service) *Router {
	return &Router{
		Logger:     service.Logger,
		Service:    service,
		Router:     httprouter.New(),
		Middleware: NewMiddlewares(service.Logger, service),
		Article:    NewArticleHandler(service.Article, service.Logger),
		Audio:      NewAudioHandler(service.Audio, service.Logger),
		Author:     NewAuthorHandler(service.Author, service.Logger),
		Book:       NewBookHandler(service.Book, service.Logger),
		Review:     NewReviewHandler(service.Review, service.Logger),
		Tag:        NewTagHandler(service.Tag, service.Logger),
		Video:      NewVideoHandler(service.Video, service.Logger),
		Direction:  NewDirectionHandler(service.Direction, service.Logger),
	}
}

func (r *Router) InitRoutes() {
	r.Article.Register(r.Router, &r.Middleware)
	r.Audio.Register(r.Router, &r.Middleware)
	r.Author.Register(r.Router)
	r.Book.Register(r.Router, &r.Middleware)
	r.Review.Register(r.Router)
	r.Tag.Register(r.Router)
	r.Video.Register(r.Router, &r.Middleware)
	r.Direction.Register(r.Router)
}
