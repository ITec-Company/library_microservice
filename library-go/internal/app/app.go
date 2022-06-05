package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"library-go/internal/composite"
	"library-go/internal/handler/http"
	"library-go/internal/store/postgres"
	"library-go/pkg/db"
	"library-go/pkg/logging"
	netHTTP "net/http"
)

func Run() {
	logger := logging.GetLogger("../../logs", "all.log")

	logger.Info("Initializing config...")
	config := Config{}
	config.NewConfig()

	logger.Infof("Initializing postgres DB at %s...", config.PgDataSourceName())
	postgresDB := db.Postgres{}
	if err := postgresDB.NewDB(config.PgDataSourceName(), logger); err != nil {
		return
	}

	logger.Info("Initializing store...")
	var store postgres.Store
	store.NewDB(postgresDB.DB, logger)

	middlewares := http.NewMiddlewares(logger)

	logger.Info("Initializing postgres composites...")
	var postgresComposites composite.Composites
	postgresComposites.NewPostgres(store, &middlewares)

	logger.Info("Initializing httprouter...")
	router := httprouter.New()
	ConfigureRouter(router, postgresComposites)

	logger.Info("Initializing CORS...")
	CORS := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{
			netHTTP.MethodHead,
			netHTTP.MethodGet,
			netHTTP.MethodPost,
			netHTTP.MethodPut,
			netHTTP.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Access-Token"},
	})

	handler := CORS.Handler(router)

	logger.Info("Initializing server...")
	var server Server
	server.New(&config, router, &handler, logger)
	if err := server.Start(); err != nil {
		logger.Fatal("Server falls")
	}
}
