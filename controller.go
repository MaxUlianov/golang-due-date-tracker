package main

import (
	"database/sql"
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
func getRecords(userId string) ([]dataRecord, error) {
	var records []dataRecord
	var recUserId string

	rows, err := db.Query("SELECT * FROM records WHERE created_by = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("error getting records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var rec dataRecord

		if err := rows.Scan(&rec.Id, &rec.Title, &rec.Comment, &rec.LastDate, &recUserId); err != nil {
			return nil, fmt.Errorf("error reading records: %v", err)
		}

		records = append(records, rec)
	}

	return records, nil
}

func getRecordById(recordId string, userId string) (dataRecord, error) {
	var rec dataRecord
	var recUserId string

	row := db.QueryRow("SELECT * FROM records WHERE id = $1 AND created_by = $2", recordId, userId)
	if err := row.Scan(&rec.Id, &rec.Title, &rec.Comment, &rec.LastDate, &recUserId); err != nil {
		if err == sql.ErrNoRows {
			return rec, fmt.Errorf("recordById %s: no such record", recordId)
		}
		return rec, fmt.Errorf("recordById %s: %v", recordId, err)
	}
	return rec, nil
}

func createRecord(record dataRecord, userId string) (string, error) {

	result, err := db.Exec("INSERT INTO records (title, comment, last_date, created_by) VALUES ($1, $2, to_timestamp($3, 'YYYY-MM-DD'), $4)", record.Title, record.Comment, record.LastDate, userId)
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	return string(id), nil
}

func updateRecord(record dataRecord, userId string) (string, error) {

	result, err := db.Exec("UPDATE records SET title = $1, comment = $2, last_date = to_timestamp($3, 'YYYY-MM-DD') WHERE id = $4 AND created_by = $5", record.Title, record.Comment, record.LastDate, record.Id, userId)
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	return string(id), nil
}

func deleteRecord(recordId string, userId string) (string, error) {
	_, err := db.Exec("DELETE FROM records WHERE id = $1 AND created_by = $2", recordId, userId)

	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}

	return "1", nil
}
