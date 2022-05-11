package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Host     string
		Port     int
		LogLevel string
		Store    string
	}
	PostgresDB struct {
		Host     string
		Port     int
		Username string
		Password string
		DbName   string
		SSLMODE  string
	}
}

func init() {
	err := godotenv.Load(" .env")
	if err != nil {
		log.Print(err)
	}
}

func (c *Config) NewConfig() {
	c.Server.Host = os.Getenv("SERVER_HOST")
	c.Server.Port, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	c.Server.LogLevel = os.Getenv("SEVER_LOG_LEVEL")
	c.Server.Store = os.Getenv("SERVER_STORE")

	c.PostgresDB.Host = os.Getenv("POSTGRES_HOST")
	c.PostgresDB.DbName = os.Getenv("POSTGRES_DB")
	c.PostgresDB.Username = os.Getenv("POSTGRES_USER")
	c.PostgresDB.Password = os.Getenv("POSTGRES_PASSWORD")
	c.PostgresDB.Port, _ = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	c.PostgresDB.SSLMODE = os.Getenv("POSTGRES_SSLMODE")
}

func (c *Config) PgDataSourceName() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresDB.Host,
		c.PostgresDB.Port,
		c.PostgresDB.Username,
		c.PostgresDB.Password,
		c.PostgresDB.DbName,
		c.PostgresDB.SSLMODE,
	)
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
