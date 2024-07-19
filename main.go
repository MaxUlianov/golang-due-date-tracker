package main

import (
	"log"
)

func main() {
	handlePostgres()

	// records, err := getRecords()
	// if err != nil {
	// 	log.Printf("error in get records func: %v", err)
	// }
	// fmt.Printf("Records found: \n%v\n", records)

	// var data = dataRecord{
	// 	Id:       "a15893d3-9e39-409f-b985-6094bca7c72e",
	// 	Title:    "Record 3",
	// 	Comment:  "This is the third (FIXED) record.",
	// 	LastDate: time.Now(),
	// }

	// log.Println(data)

	// res, err := createRecord(data)
	// res, err := updateRecord(data)
	// res, err := deleteRecord(data.Id)
	// log.Printf("res: %s, err: %v", res, err)

	log.Println("running server")
	runServer()
}
