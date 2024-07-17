package main

import (
	"fmt"
	"time"
)

type dataRecord struct {
	Id       string
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

func createRecord(record dataRecord) (string, error) {

	result, err := db.Exec("INSERT INTO records (title, comment, last_date) VALUES ($1, $2, to_timestamp($3, 'YYYY-MM-DD'))", record.Title, record.Comment, record.LastDate)
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	return string(id), nil
}

func updateRecord(record dataRecord) (string, error) {

	result, err := db.Exec("UPDATE records SET title = ?, comment = ?, last_date = ?", record.Title, record.Comment, record.LastDate)
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	return string(id), nil
}

func deleteRecord(record dataRecord) (string, error) {
	_, err := db.Exec("DELETE FROM records WHERE id = ?", record.Id)

	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}

	return "1", nil
}
