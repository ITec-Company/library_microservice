package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"library-go/pkg/logging"
)

type Postgres struct {
	DB     *sql.DB
	Logger *logging.Logger
}

func (p *Postgres) NewDB(dbSourceName string, logger *logging.Logger) error {
	p.Logger = logger

	db, err := sql.Open("postgres", dbSourceName)
	if err != nil {
		p.Logger.Errorf("error occurred while opening DB. err: %v", err)
		return err
	}
	p.DB = db

	if err := p.Ping(); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) Ping() error {
	err := p.DB.Ping()
	if err != nil {
		p.Logger.Errorf("error occurred while pinging DB. err: %v", err)
	}
	return nil
}

func (p *Postgres) Close() error {
	err := p.Close()
	if err != nil {
		p.Logger.Errorf("error occurred while closing DB. err: %v", err)
	}
	return nil
}
