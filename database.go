package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Params   string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		Params:   "sslmode=disable",
	}
}

func handlePostgres() {
	postgresConfig := NewDBConfig()

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
	db, err = sql.Open("postgres", connectString)

	if err != nil {
		log.Fatal(err)
	}

	// check connect
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Database Connected!")
}
