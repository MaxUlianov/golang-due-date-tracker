package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	handlePostgres()

	records, err := getRecords()
	if err != nil {
		log.Printf("error in get records func: %v", err)
	}
	fmt.Printf("Records found: \n%v\n", records)

	var data = dataRecord{
		Title:    "Record 1",
		Comment:  "This is the first record.",
		LastDate: time.Now(),
	}

	log.Println(data)

	res, err := createRecord(data)
	log.Printf("res: %s, err: %v", res, err)

	// log.Println("running server")
	// runServer()
}
