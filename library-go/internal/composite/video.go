package composite

import (
	"library-go/internal/handler"
	"library-go/internal/handler/http"
	"library-go/internal/service"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type VideoComposite struct {
	logger  *logging.Logger
	Storage store.VideoStorage
	Service service.VideoService
	Handler handler.Handler
}

func (vc *VideoComposite) New(storage store.VideoStorage, logger *logging.Logger, middleware *http.Middleware) {
	vc.logger = logger
	vc.Storage = storage
	vc.Service = service.NewService(vc.Storage, vc.logger)
	vc.Handler = http.NewVideoHandler(vc.Service, vc.logger, middleware)
}
