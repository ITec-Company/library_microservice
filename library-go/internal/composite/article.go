package composite

import (
	"library-go/internal/handler"
	"library-go/internal/handler/http"
	"library-go/internal/service"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type ArticleComposite struct {
	logger  *logging.Logger
	Storage store.ArticleStorage
	Service service.ArticleService
	Handler handler.Handler
}

func (ac *ArticleComposite) New(storage store.ArticleStorage, logger *logging.Logger, middleware *http.Middleware) {
	ac.logger = logger
	ac.Storage = storage
	ac.Service = service.NewArticleService(ac.Storage, ac.logger)
	ac.Handler = http.NewArticleHandler(ac.Service, ac.logger, middleware)
}
