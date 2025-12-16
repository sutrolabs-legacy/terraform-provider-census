package client

import (
	"encoding/json"
	"testing"
)

// TestSyncScheduleDayUnmarshal tests that schedule_day can be unmarshaled as a string
// This reproduces the bug where Census API returns schedule_day as a string (e.g. "Monday")
// but the struct expected *int, causing JSON unmarshal errors
func TestSyncScheduleDayUnmarshal(t *testing.T) {
	// This is the actual JSON response format from Census API
	// schedule_day is a string, not an int
	jsonResponse := `{
		"id": 123,
		"label": "Test Sync",
		"status": "success",
		"schedule_frequency": "weekly",
		"schedule_day": "Monday",
		"schedule_hour": 6,
		"schedule_minute": 0,
		"cron_expression": "0 6 * * 1",
		"operation": "mirror",
		"paused": false
	}`

	var sync Sync
	err := json.Unmarshal([]byte(jsonResponse), &sync)

	// Before the fix: This would fail with error:
	// "json: cannot unmarshal string into Go struct field Sync.schedule_day of type int"

	// After the fix: This should succeed
	if err != nil {
		t.Fatalf("Failed to unmarshal sync with schedule_day as string: %v", err)
	}

	// Verify the fields were unmarshaled correctly
	if sync.ID != 123 {
		t.Errorf("Expected ID 123, got %d", sync.ID)
	}

	if sync.ScheduleDay == nil {
		t.Fatal("Expected schedule_day to be set, got nil")
	}

	if *sync.ScheduleDay != "Monday" {
		t.Errorf("Expected schedule_day 'Monday', got '%s'", *sync.ScheduleDay)
	}

	if sync.ScheduleHour == nil {
		t.Fatal("Expected schedule_hour to be set, got nil")
	}

	if *sync.ScheduleHour != 6 {
		t.Errorf("Expected schedule_hour 6, got %d", *sync.ScheduleHour)
	}

	if sync.ScheduleMinute == nil {
		t.Fatal("Expected schedule_minute to be set, got nil")
	}

	if *sync.ScheduleMinute != 0 {
		t.Errorf("Expected schedule_minute 0, got %d", *sync.ScheduleMinute)
	}

	if sync.CronExpression != "0 6 * * 1" {
		t.Errorf("Expected cron_expression '0 6 * * 1', got '%s'", sync.CronExpression)
	}
}

// TestSyncScheduleDayNull tests that null schedule_day is handled correctly
func TestSyncScheduleDayNull(t *testing.T) {
	jsonResponse := `{
		"id": 456,
		"label": "Test Sync 2",
		"status": "success",
		"schedule_day": null,
		"schedule_hour": null,
		"schedule_minute": null,
		"operation": "upsert",
		"paused": false
	}`

	var sync Sync
	err := json.Unmarshal([]byte(jsonResponse), &sync)

	if err != nil {
		t.Fatalf("Failed to unmarshal sync with null schedule fields: %v", err)
	}

	// Verify null values are handled correctly
	if sync.ScheduleDay != nil {
		t.Errorf("Expected schedule_day to be nil, got '%s'", *sync.ScheduleDay)
	}

	if sync.ScheduleHour != nil {
		t.Errorf("Expected schedule_hour to be nil, got %d", *sync.ScheduleHour)
	}

	if sync.ScheduleMinute != nil {
		t.Errorf("Expected schedule_minute to be nil, got %d", *sync.ScheduleMinute)
	}
}

// TestSyncScheduleDayOmitted tests that omitted schedule fields are handled correctly
func TestSyncScheduleDayOmitted(t *testing.T) {
	jsonResponse := `{
		"id": 789,
		"label": "Test Sync 3",
		"status": "success",
		"operation": "append",
		"paused": true
	}`

	var sync Sync
	err := json.Unmarshal([]byte(jsonResponse), &sync)

	if err != nil {
		t.Fatalf("Failed to unmarshal sync with omitted schedule fields: %v", err)
	}

	// Verify omitted values are nil
	if sync.ScheduleDay != nil {
		t.Errorf("Expected schedule_day to be nil when omitted, got '%s'", *sync.ScheduleDay)
	}

	if sync.ScheduleHour != nil {
		t.Errorf("Expected schedule_hour to be nil when omitted, got %d", *sync.ScheduleHour)
	}

	if sync.ScheduleMinute != nil {
		t.Errorf("Expected schedule_minute to be nil when omitted, got %d", *sync.ScheduleMinute)
	}

	if sync.Paused != true {
		t.Errorf("Expected paused to be true, got %v", sync.Paused)
	}
}
