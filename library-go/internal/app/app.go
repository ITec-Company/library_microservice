package app

import (
	"github.com/julienschmidt/httprouter"
	"library-go/internal/composite"
	"library-go/internal/store/postgres"
	"library-go/pkg/db"
	"library-go/pkg/logging"
)

func Run() {
	logger := logging.GetLogger()

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

	logger.Info("Initializing postgres composites...")
	var postgresComposites composite.Composites
	postgresComposites.NewPostgres(store)

	logger.Info("Initializing httprouter...")
	router := httprouter.New()
	ConfigureRouter(router, postgresComposites)

	logger.Info("Initializing server...")
	var server Server
	server.New(&config, router, logger)
	if err := server.Start(); err != nil {
		logger.Fatal("Server falls")
	}
}
