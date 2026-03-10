package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	host        string
	port        int
	databaseUrl string
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Dbname   string
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
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	portNum, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return &Config{
		host: host,
		port: portNum,
	}
}

func LoadConfig(cf ...ConfigFunc) *Config {
	config := defaultConfig()
	if config == nil {
		return nil
	}
	for _, f := range cf {
		f(config)
	}
	return config
}

func LoadDatabaseConfig() *DBConfig {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	dbname := os.Getenv("DATABASE_DBNAME")
	config := DBConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Dbname:   dbname,
	}

	return &config
}
