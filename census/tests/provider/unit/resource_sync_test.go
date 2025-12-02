package unit_test

import (
	"reflect"
	"testing"

	"github.com/sutrolabs/terraform-provider-census/census/client"
	"github.com/sutrolabs/terraform-provider-census/census/provider"
)

// Unit tests for sync resource helper functions
// These tests do NOT require API credentials or external dependencies

// ============================================================================
// Field Mapping Tests
// ============================================================================

func TestExpandFieldMappings_Direct(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []client.FieldMapping
	}{
		{
			name: "direct mapping with primary identifier",
			input: []interface{}{
				map[string]interface{}{
					"from":                  "email",
					"to":                    "Email",
					"type":                  "direct",
					"is_primary_identifier": true,
				},
			},
			expected: []client.FieldMapping{
				{
					From:                "email",
					To:                  "Email",
					Type:                "direct",
					IsPrimaryIdentifier: true,
				},
			},
		},
		{
			name: "direct mapping without primary identifier",
			input: []interface{}{
				map[string]interface{}{
					"from": "first_name",
					"to":   "FirstName",
				},
			},
			expected: []client.FieldMapping{
				{
					From: "first_name",
					To:   "FirstName",
					Type: "direct", // Default
				},
			},
		},
		{
			name: "multiple direct mappings",
			input: []interface{}{
				map[string]interface{}{
					"from":                  "email",
					"to":                    "Email",
					"is_primary_identifier": true,
				},
				map[string]interface{}{
					"from": "first_name",
					"to":   "FirstName",
				},
				map[string]interface{}{
					"from": "last_name",
					"to":   "LastName",
				},
			},
			expected: []client.FieldMapping{
				{
					From:                "email",
					To:                  "Email",
					Type:                "direct",
					IsPrimaryIdentifier: true,
				},
				{
					From: "first_name",
					To:   "FirstName",
					Type: "direct",
				},
				{
					From: "last_name",
					To:   "LastName",
					Type: "direct",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandFieldMappings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandFieldMappings() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestExpandFieldMappings_Constant(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []client.FieldMapping
	}{
		{
			name: "constant string value",
			input: []interface{}{
				map[string]interface{}{
					"to":       "Source",
					"type":     "constant",
					"constant": "Website",
				},
			},
			expected: []client.FieldMapping{
				{
					To:       "Source",
					Type:     "constant",
					Constant: "Website",
				},
			},
		},
		{
			name: "constant numeric value",
			input: []interface{}{
				map[string]interface{}{
					"to":       "Priority",
					"type":     "constant",
					"constant": 1,
				},
			},
			expected: []client.FieldMapping{
				{
					To:       "Priority",
					Type:     "constant",
					Constant: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandFieldMappings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandFieldMappings() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestExpandFieldMappings_LiquidTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []client.FieldMapping
	}{
		{
			name: "liquid template mapping",
			input: []interface{}{
				map[string]interface{}{
					"to":              "FullName",
					"type":            "liquid_template",
					"liquid_template": "{{ first_name }} {{ last_name }}",
				},
			},
			expected: []client.FieldMapping{
				{
					To:             "FullName",
					Type:           "liquid_template",
					LiquidTemplate: "{{ first_name }} {{ last_name }}",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandFieldMappings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandFieldMappings() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestExpandFieldMappings_SyncMetadata(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []client.FieldMapping
	}{
		{
			name: "sync metadata mapping",
			input: []interface{}{
				map[string]interface{}{
					"to":                "Sync_Run_ID",
					"type":              "sync_metadata",
					"sync_metadata_key": "sync_run_id",
				},
			},
			expected: []client.FieldMapping{
				{
					To:              "Sync_Run_ID",
					Type:            "sync_metadata",
					SyncMetadataKey: "sync_run_id",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandFieldMappings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandFieldMappings() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestExpandFieldMappings_Empty(t *testing.T) {
	result := provider.ExpandFieldMappings([]interface{}{})
	if len(result) != 0 {
		t.Errorf("ExpandFieldMappings(empty) should return empty slice, got %d items", len(result))
	}
}

func TestFlattenFieldMappings_Direct(t *testing.T) {
	tests := []struct {
		name     string
		input    []client.FieldMapping
		expected int // We'll check length and specific fields
	}{
		{
			name: "direct mapping with primary identifier",
			input: []client.FieldMapping{
				{
					From:                "email",
					To:                  "Email",
					Type:                "direct",
					IsPrimaryIdentifier: true,
				},
			},
			expected: 1,
		},
		{
			name: "multiple direct mappings",
			input: []client.FieldMapping{
				{
					From:                "email",
					To:                  "Email",
					Type:                "direct",
					IsPrimaryIdentifier: true,
				},
				{
					From: "first_name",
					To:   "FirstName",
					Type: "direct",
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.FlattenFieldMappings(tt.input)
			if len(result) != tt.expected {
				t.Errorf("FlattenFieldMappings() returned %d items, want %d", len(result), tt.expected)
			}

			// Check first mapping
			if len(result) > 0 {
				firstMapping := result[0].(map[string]interface{})
				if firstMapping["from"] != tt.input[0].From {
					t.Errorf("FlattenFieldMappings() first mapping from = %v, want %v", firstMapping["from"], tt.input[0].From)
				}
				if firstMapping["to"] != tt.input[0].To {
					t.Errorf("FlattenFieldMappings() first mapping to = %v, want %v", firstMapping["to"], tt.input[0].To)
				}
			}
		})
	}
}

func TestFlattenFieldMappings_Empty(t *testing.T) {
	result := provider.FlattenFieldMappings([]client.FieldMapping{})
	if len(result) != 0 {
		t.Errorf("FlattenFieldMappings(empty) should return empty slice, got %d items", len(result))
	}
}

// ============================================================================
// Alert Tests
// ============================================================================

func TestExpandAlerts_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected int // Check length since AlertAttribute has Options map
	}{
		{
			name: "basic alert",
			input: []interface{}{
				map[string]interface{}{
					"type":                 "email",
					"send_for":             "failure",
					"should_send_recovery": true,
					"emails":               []interface{}{"admin@example.com"},
				},
			},
			expected: 1,
		},
		{
			name: "multiple alerts",
			input: []interface{}{
				map[string]interface{}{
					"type":                 "email",
					"send_for":             "failure",
					"should_send_recovery": true,
					"emails":               []interface{}{"failure@example.com"},
				},
				map[string]interface{}{
					"type":                 "email",
					"send_for":             "success",
					"should_send_recovery": false,
					"emails":               []interface{}{"success@example.com"},
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandAlerts(tt.input)
			if len(result) != tt.expected {
				t.Errorf("ExpandAlerts() returned %d items, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestExpandAlerts_Empty(t *testing.T) {
	result := provider.ExpandAlerts([]interface{}{})
	if result != nil {
		t.Errorf("ExpandAlerts(empty) should return nil, got %d items", len(result))
	}
}

// TestExpandAlerts_EmptyStrings tests that alerts with empty string values
// are properly skipped, preventing invalid API payloads (regression test)
func TestExpandAlerts_EmptyStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected int
		checkFn  func(*testing.T, []client.AlertAttribute)
	}{
		{
			name: "alert with empty type is skipped",
			input: []interface{}{
				map[string]interface{}{
					"type":                 "", // Empty string
					"send_for":             "first_time",
					"should_send_recovery": false,
					"options":              map[string]interface{}{},
				},
			},
			expected: 0, // Should be skipped
		},
		{
			name: "alert with empty send_for uses default",
			input: []interface{}{
				map[string]interface{}{
					"type":                 "FailureAlertConfiguration",
					"send_for":             "", // Empty string should trigger default
					"should_send_recovery": true,
					"options":              map[string]interface{}{},
				},
			},
			expected: 1,
			checkFn: func(t *testing.T, alerts []client.AlertAttribute) {
				if alerts[0].SendFor != "first_time" {
					t.Errorf("Expected SendFor to be 'first_time' (default), got '%s'", alerts[0].SendFor)
				}
			},
		},
		{
			name: "mixed valid and invalid alerts",
			input: []interface{}{
				map[string]interface{}{
					"type":                 "", // Empty - should be skipped
					"send_for":             "first_time",
					"should_send_recovery": false,
				},
				map[string]interface{}{
					"type":                 "FailureAlertConfiguration", // Valid
					"send_for":             "every_time",
					"should_send_recovery": true,
				},
				map[string]interface{}{
					"type":                 "", // Empty - should be skipped
					"send_for":             "",
					"should_send_recovery": false,
				},
			},
			expected: 1, // Only the valid alert should remain
			checkFn: func(t *testing.T, alerts []client.AlertAttribute) {
				if alerts[0].Type != "FailureAlertConfiguration" {
					t.Errorf("Expected Type to be 'FailureAlertConfiguration', got '%s'", alerts[0].Type)
				}
				if alerts[0].SendFor != "every_time" {
					t.Errorf("Expected SendFor to be 'every_time', got '%s'", alerts[0].SendFor)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandAlerts(tt.input)
			if len(result) != tt.expected {
				t.Errorf("ExpandAlerts() returned %d items, want %d", len(result), tt.expected)
			}
			if tt.checkFn != nil {
				tt.checkFn(t, result)
			}
		})
	}
}

// ============================================================================
// Schedule Tests
// ============================================================================

func TestExpandSyncSchedule_Hourly(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		wantFreq string
	}{
		{
			name: "hourly schedule with minute",
			input: []interface{}{
				map[string]interface{}{
					"frequency": "hourly",
					"minute":    30,
				},
			},
			wantFreq: "hourly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandSyncSchedule(tt.input)
			if result == nil {
				t.Errorf("ExpandSyncSchedule() returned nil")
				return
			}
			if result.Frequency != tt.wantFreq {
				t.Errorf("ExpandSyncSchedule() frequency = %v, want %v", result.Frequency, tt.wantFreq)
			}
		})
	}
}

func TestExpandSyncSchedule_Daily(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		wantFreq string
	}{
		{
			name: "daily schedule at 9am",
			input: []interface{}{
				map[string]interface{}{
					"frequency": "daily",
					"hour":      9,
					"minute":    0,
				},
			},
			wantFreq: "daily",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandSyncSchedule(tt.input)
			if result == nil {
				t.Errorf("ExpandSyncSchedule() returned nil")
				return
			}
			if result.Frequency != tt.wantFreq {
				t.Errorf("ExpandSyncSchedule() frequency = %v, want %v", result.Frequency, tt.wantFreq)
			}
		})
	}
}

func TestExpandSyncSchedule_Nil(t *testing.T) {
	result := provider.ExpandSyncSchedule([]interface{}{})
	if result != nil {
		t.Errorf("ExpandSyncSchedule(empty) should return nil, got %+v", result)
	}
}

// ============================================================================
// Utility Helper Tests
// ============================================================================

func TestExpandStringList(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []string
	}{
		{
			name:     "empty list",
			input:    []interface{}{},
			expected: []string{},
		},
		{
			name:     "single item",
			input:    []interface{}{"item1"},
			expected: []string{"item1"},
		},
		{
			name:     "multiple items",
			input:    []interface{}{"item1", "item2", "item3"},
			expected: []string{"item1", "item2", "item3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandStringList(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandStringList() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestExpandStringMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name: "simple map",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "map with empty string values - kept as-is",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "",
				"key3": "value3",
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "",
				"key3": "value3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandStringMap(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandStringMap() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestCleanEmptyStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name: "map with empty strings - removed",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "",
				"key3": "value3",
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key3": "value3",
			},
		},
		{
			name: "map with zero cohort_id - removed",
			input: map[string]interface{}{
				"key1":      "value1",
				"cohort_id": 0,
				"key3":      "value3",
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key3": "value3",
			},
		},
		{
			name: "map with non-zero integers - kept",
			input: map[string]interface{}{
				"key1":  "value1",
				"count": 0,
				"key3":  "value3",
			},
			expected: map[string]interface{}{
				"key1":  "value1",
				"count": 0,
				"key3":  "value3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.CleanEmptyStrings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CleanEmptyStrings() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Primary Identifier Validation Tests
// ============================================================================
// These tests verify that client-side validation of primary identifiers has been removed.
// The Census API now handles primary identifier validation, allowing destinations like
// Google Sheets that don't require explicit primary identifiers.

func TestExpandFieldMappings_NoPrimaryIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []client.FieldMapping
	}{
		{
			name: "mappings without primary identifier - allowed for destinations like Google Sheets",
			input: []interface{}{
				map[string]interface{}{
					"from": "city",
					"to":   "City",
					"type": "direct",
				},
				map[string]interface{}{
					"from": "state",
					"to":   "State",
					"type": "direct",
				},
			},
			expected: []client.FieldMapping{
				{
					From:                "city",
					To:                  "City",
					Type:                "direct",
					IsPrimaryIdentifier: false,
				},
				{
					From:                "state",
					To:                  "State",
					Type:                "direct",
					IsPrimaryIdentifier: false,
				},
			},
		},
		{
			name: "single mapping without primary identifier",
			input: []interface{}{
				map[string]interface{}{
					"from": "email",
					"to":   "Email",
				},
			},
			expected: []client.FieldMapping{
				{
					From:                "email",
					To:                  "Email",
					Type:                "direct",
					IsPrimaryIdentifier: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandFieldMappings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandFieldMappings() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestExpandFieldMappings_GoogleSheetsScenario(t *testing.T) {
	// Google Sheets syncs don't require primary identifiers
	// This test verifies that field mappings work correctly for such destinations
	tests := []struct {
		name     string
		input    []interface{}
		expected int // Expected number of mappings
	}{
		{
			name: "Google Sheets sync with replace operation - no primary identifier needed",
			input: []interface{}{
				map[string]interface{}{
					"from": "city",
					"to":   "City",
				},
				map[string]interface{}{
					"from": "state",
					"to":   "State",
				},
				map[string]interface{}{
					"from": "population",
					"to":   "Population",
				},
			},
			expected: 3,
		},
		{
			name: "Google Sheets sync with constant values - no primary identifier",
			input: []interface{}{
				map[string]interface{}{
					"from": "name",
					"to":   "Name",
				},
				map[string]interface{}{
					"type":     "constant",
					"constant": "2024",
					"to":       "Year",
				},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ExpandFieldMappings(tt.input)
			if len(result) != tt.expected {
				t.Errorf("ExpandFieldMappings() returned %d mappings, want %d", len(result), tt.expected)
			}

			// Verify all mappings are properly expanded
			for i, mapping := range result {
				if mapping.To == "" {
					t.Errorf("ExpandFieldMappings() mapping %d has empty 'to' field", i)
				}
			}

			// Verify no primary identifiers are set
			for i, mapping := range result {
				if mapping.IsPrimaryIdentifier {
					t.Errorf("ExpandFieldMappings() mapping %d unexpectedly has IsPrimaryIdentifier=true", i)
				}
			}
		})
	}
}
