package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	host        string
	port        int
	databaseUrl string
}

type ConfigFunc func(c *Config)

func (c *Config) GetHost() string {
	return c.host
}

func (c *Config) GetPort() int {
	return c.port
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func defaultConfig() *Config {
	return &Config{
		host: "localhost",
		port: 8080,
	}
}

func LoadConfig(cf ...ConfigFunc) *Config {
	config := defaultConfig()
	for _, f := range cf {
		f(config)
	}
	return config
}

func LoadDatabaseConfig() (*pgxpool.Config, error) {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	dbname := os.Getenv("DATABASE_DBNAME")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		user,
		password,
		host,
		port,
		dbname,
	)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	// Good defaults to avoid too many connections on Supabase
	config.MaxConns = 5
	config.MinConns = 1
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 30 * time.Second
	return config, nil
}
