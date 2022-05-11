package composite

import (
	"library-go/internal/handler"
	"library-go/internal/handler/http"
	"library-go/internal/service"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type TagComposite struct {
	logger  *logging.Logger
	Storage store.TagStorage
	Service service.TagService
	Handler handler.Handler
}

func (tc *TagComposite) New(storage store.TagStorage, logger *logging.Logger) {
	tc.logger = logger
	tc.Storage = storage
	tc.Service = service.NewTagService(tc.Storage, tc.logger)
	tc.Handler = http.NewTagHandler(tc.Service, tc.logger)
}
