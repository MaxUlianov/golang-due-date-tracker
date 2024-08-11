package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type dataRecord struct {
	Id       string
	Title    string
	Comment  string
	LastDate time.Time
}

type dataRecordExtended struct {
	dataRecord
	TimeLeft int64
	TimeCode int64
}

type Config struct {
	DefaultInterval string `json:"defaultInterval"`
}

func LoadInterval(filename string) time.Duration {
	var defaultInterval = 365 * 24 * time.Hour

	var config Config

	// read the config file
	data, err := os.ReadFile(filename)
	if err != nil {
		return defaultInterval
	}

	// get JSON data into the config struct
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultInterval
	}

	// rarse the config interval string into time.Duration
	interval, err := time.ParseDuration(config.DefaultInterval)
	if err != nil {
		return defaultInterval
	}

	return interval
}

func checkTimeInterval(t time.Time, interval time.Duration) (int64, int64) {
	// add interval to t time to add the expiration time
	t = t.Add(interval)

	now := time.Now()
	duration := t.Sub(now)

	threshold := interval / 12

	if duration < 0 {
		return -1 * int64(duration.Hours()/24), 2 // due
	}
	if duration > 0 && duration <= threshold {
		return int64(duration.Hours() / 24), 1 // soon
	}

	return int64(duration.Hours() / 24), 0 // a lot of time left
}

// define the interval
var defaultInterval time.Duration = LoadInterval("config.json")

// func to get all records
func getRecords(db DBInstance, userId string) ([]dataRecordExtended, error) {
	var records []dataRecordExtended
	var recUserId string

	rows, err := db.self.Query("SELECT * FROM records WHERE created_by = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("error getting records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var rec dataRecordExtended

		if err := rows.Scan(&rec.Id, &rec.Title, &rec.Comment, &rec.LastDate, &recUserId); err != nil {
			return nil, fmt.Errorf("error reading records: %v", err)
		}

		rec.TimeLeft, rec.TimeCode = checkTimeInterval(rec.LastDate, defaultInterval)

		records = append(records, rec)
	}

	return records, nil
}

func getRecordById(db DBInstance, recordId string, userId string) (dataRecordExtended, error) {
	var rec dataRecordExtended
	var recUserId string

	row := db.self.QueryRow("SELECT * FROM records WHERE id = $1 AND created_by = $2", recordId, userId)
	if err := row.Scan(&rec.Id, &rec.Title, &rec.Comment, &rec.LastDate, &recUserId); err != nil {
		if err == sql.ErrNoRows {
			return rec, fmt.Errorf("recordById %s: no such record", recordId)
		}
		return rec, fmt.Errorf("recordById %s: %v", recordId, err)
	}

	rec.TimeLeft, rec.TimeCode = checkTimeInterval(rec.LastDate, defaultInterval)

	return rec, nil
}

func createRecord(db DBInstance, record dataRecord, userId string) (string, error) {

	result, err := db.self.Exec("INSERT INTO records (title, comment, last_date, created_by) VALUES ($1, $2, to_timestamp($3, 'YYYY-MM-DD'), $4)", record.Title, record.Comment, record.LastDate, userId)
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	return fmt.Sprintf("%d", id), nil
}

func updateRecord(db DBInstance, record dataRecord, userId string) (string, error) {

	result, err := db.self.Exec("UPDATE records SET title = $1, comment = $2, last_date = to_timestamp($3, 'YYYY-MM-DD') WHERE id = $4 AND created_by = $5", record.Title, record.Comment, record.LastDate, record.Id, userId)
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}
	return fmt.Sprintf("%d", id), nil
}

func deleteRecord(db DBInstance, recordId string, userId string) (string, error) {
	_, err := db.self.Exec("DELETE FROM records WHERE id = $1 AND created_by = $2", recordId, userId)

	if err != nil {
		return "0", fmt.Errorf("addRecord: %v", err)
	}

	return "1", nil
}
