package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestLoadInterval(t *testing.T) {
	// create a temporary config file
	configData := Config{DefaultInterval: "48h"}
	data, _ := json.Marshal(configData)
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	// delete temporary file
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Ttst LoadInterval function
	interval := LoadInterval(tmpFile.Name())
	expected := 48 * time.Hour
	if interval != expected {
		t.Errorf("Expected %v, got %v", expected, interval)
	}
}

func TestCheckTimeInterval(t *testing.T) {
	now := time.Now()
	interval := 12 * 24 * time.Hour // 12 days

	tests := []struct {
		name     string
		lastDate time.Time
		expected int64
		code     int64
	}{
		{"Due", now.Add(-13 * 24 * time.Hour), 1, 2},
		{"Soon", now.Add(-11 * 24 * time.Hour), 0, 1},
		{"A Lot Of Time Left", now.Add(-1 * 24 * time.Hour), 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days, code := checkTimeInterval(tt.lastDate, interval)
			if days != tt.expected || code != tt.code {
				t.Errorf("Expected (%d, %d), got (%d, %d)", tt.expected, tt.code, days, code)
			}
		})
	}
}
