package composite

import (
	"library-go/internal/handler"
	"library-go/internal/handler/http"
	"library-go/internal/service"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type BookComposite struct {
	logger  *logging.Logger
	Storage store.BookStorage
	Service service.BookService
	Handler handler.Handler
}

func (bc *BookComposite) New(storage store.BookStorage, logger *logging.Logger, middleware *http.Middleware) {
	bc.logger = logger
	bc.Storage = storage
	bc.Service = service.NewBookService(bc.Storage, bc.logger)
	bc.Handler = http.NewBookHandler(bc.Service, bc.logger, middleware)
}
