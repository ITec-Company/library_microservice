package store

import (
	"database/sql"
	"errors"
	"fmt"
	"library/internal/config"
	"library/pkg/logging"

	_ "github.com/lib/pq" // ...
)

var (
	// ErrNoRowsAffected ...
	ErrNoRowsAffected = errors.New("No rows affected")
)

// Store ...
type Store struct {
	Config *config.Config
	Db     *sql.DB
	Logger *logging.Logger
}

// New ...
func New(config *config.Config) *Store {
	return &Store{
		Config: config,
		Logger: logging.GetLogger(),
	}
}

// Open ...
func (s *Store) Open() error {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		s.Config.DataBase.Host,
		s.Config.DataBase.Port,
		s.Config.DataBase.Username,
		s.Config.DataBase.Password,
		s.Config.DataBase.DbName,
		s.Config.DataBase.Sslmode,
	)
	s.Logger.Infof("Library store opening. Source: %s", dataSourceName)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		s.Logger.Errorf("Can't open DB. Err msg: %v", err)
		return err
	}

	if err := db.Ping(); err != nil {
		s.Logger.Errorf("Can't ping DB. Err msg: %v", err)
		return err
	}

	s.Db = db

	return nil
}

// Close ...
func (s *Store) Close() {
	s.Db.Close()
}
