package main

import (
	"fmt"
	"time"
)

type dataRecord struct {
	Id       int64
	Title    string
	Comment  string
	LastDate time.Time
}

// func to get all records
func getRecords() ([]dataRecord, error) {
	var records []dataRecord

	rows, err := db.Query("SELECT * FROM records")
	if err != nil {
		return nil, fmt.Errorf("error getting records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var rec dataRecord

		if err := rows.Scan(&rec.Id, &rec.Title, &rec.Comment, &rec.LastDate); err != nil {
			return nil, fmt.Errorf("error reading records: %v", err)
		}

		records = append(records, rec)
	}

	return records, nil
}
