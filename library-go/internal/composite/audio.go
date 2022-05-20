package composite

import (
	"library-go/internal/handler"
	"library-go/internal/handler/http"
	"library-go/internal/service"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type AudioComposite struct {
	logger  *logging.Logger
	Storage store.AudioStorage
	Service service.AudioService
	Handler handler.Handler
}

func (ac *AudioComposite) New(storage store.AudioStorage, logger *logging.Logger, middleware *http.Middleware) {
	ac.logger = logger
	ac.Storage = storage
	ac.Service = service.NewAudioService(ac.Storage, ac.logger)
	ac.Handler = http.NewAudioHandler(ac.Service, ac.logger, middleware)
}
