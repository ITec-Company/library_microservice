package composite

import (
	"library-go/internal/handler"
	"library-go/internal/handler/http"
	"library-go/internal/service"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type ReviewComposite struct {
	logger  *logging.Logger
	Storage store.ReviewStorage
	Service service.ReviewService
	Handler handler.Handler
}

func (rc *ReviewComposite) New(storage store.ReviewStorage, logger *logging.Logger) {
	rc.logger = logger
	rc.Storage = storage
	rc.Service = service.NewReviewService(rc.Storage, rc.logger)
	rc.Handler = http.NewReviewHandler(rc.Service, rc.logger)
}
