package composite

import (
	"library-go/internal/handler"
	"library-go/internal/handler/http"
	"library-go/internal/service"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type DirectionComposite struct {
	logger  *logging.Logger
	Storage store.DirectionStorage
	Service service.DirectionService
	Handler handler.Handler
}

func (dc *DirectionComposite) New(storage store.DirectionStorage, logger *logging.Logger, middleware *http.Middleware) {
	dc.logger = logger
	dc.Storage = storage
	dc.Service = service.NewDirectionService(dc.Storage, dc.logger)
	dc.Handler = http.NewDirectionHandler(dc.Service, dc.logger, middleware)
}
