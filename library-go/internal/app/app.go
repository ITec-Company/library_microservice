package app

import (
	"github.com/rs/cors"
	"library-go/internal/handler/http"
	"library-go/internal/service"
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

	store := postgres.New(postgresDB.DB, logger)

	logger.Info("Initializing service...")

	service := service.New(store)

	logger.Info("Initializing router...")

	router := http.New(service)

	logger.Info("Initializing routes...")

	router.InitRoutes()

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

	handler := CORS.Handler(router.Router)

	logger.Info("Initializing server...")
	var server Server
	server.New(&config, router.Router, &handler, logger)
	if err := server.Start(); err != nil {
		logger.Fatal("Server falls")
	}
}
