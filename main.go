package main

import (
	"log"
)

func main() {
	log.Println("Starting app...")

	// create db of type sql.DB
	database := handlePostgres()

	// update global instance with it
	db_instance.self = database

	runServer()
}
