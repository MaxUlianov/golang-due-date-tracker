package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBInstance struct {
	self *sql.DB
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Params   string
}

var db_instance DBInstance

func NewDBConfig() (*DBConfig, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || name == "" {
		return nil, fmt.Errorf("missing environment variables for DB connection")
	}

	return &DBConfig{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Name:     name,
		Params:   "sslmode=disable",
	}, nil
}

func handlePostgres() *sql.DB {

	postgresConfig, configErr := NewDBConfig()
	if configErr != nil {
		log.Fatalf("error creating database configuration: %v", configErr)
	}

	connectString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.Name,
		postgresConfig.Params,
	)

	var err error
	database, err := sql.Open("postgres", connectString)

	if err != nil {
		log.Fatal(fmt.Errorf("database connection error: %w", err))
	}

	// check connect
	pingErr := database.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Database Connected!")

	return database
}
