package main

import (
	"fmt"
	"log"
)

func main() {
	handlePostgres()

	records, err := getRecords()
	if err != nil {
		log.Printf("error in get records func: %v", err)
	}

	fmt.Printf("Records:\n%v\n", records)

	log.Println("running server")
	fmt.Println("will launch server")
	runServer()
}
