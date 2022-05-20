package composite

import (
	"library-go/internal/handler"
	"library-go/internal/handler/http"
	"library-go/internal/service"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type AuthorComposite struct {
	logger  *logging.Logger
	Storage store.AuthorStorage
	Service service.AuthorService
	Handler handler.Handler
}

func (ac *AuthorComposite) New(storage store.AuthorStorage, logger *logging.Logger, middleware *http.Middleware) {
	ac.logger = logger
	ac.Storage = storage
	ac.Service = service.NewAuthorService(ac.Storage, ac.logger)
	ac.Handler = http.NewAuthorHandler(ac.Service, ac.logger, middleware)
}
