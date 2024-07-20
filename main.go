package main

import (
	"log"
)

func main() {
	log.Println("Starting app...")

	handlePostgres()
	runServer()
}
