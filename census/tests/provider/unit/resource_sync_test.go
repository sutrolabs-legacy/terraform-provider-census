package unit_test

import (
	"reflect"
	"sort"
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
// Field Mapping Ordering Tests
// ============================================================================
// These tests verify that field mappings are sorted by position from the Census API,
// preventing spurious diffs when the API returns mappings in non-deterministic order.

func TestSortFieldMappingsByPosition_Basic(t *testing.T) {
	// Simulate API response with mappings in wrong order
	mappings := []client.FieldMapping{
		{
			From:     "last_name",
			To:       "LastName",
			Position: 2,
			Type:     "direct",
		},
		{
			From:     "email",
			To:       "Email",
			Position: 0,
			Type:     "direct",
		},
		{
			From:     "first_name",
			To:       "FirstName",
			Position: 1,
			Type:     "direct",
		},
	}

	// Sort by position
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].Position < mappings[j].Position
	})

	// Verify result is sorted by position
	if len(mappings) != 3 {
		t.Errorf("Expected 3 mappings, got %d", len(mappings))
	}

	expectedOrder := []string{"Email", "FirstName", "LastName"}
	expectedPositions := []int{0, 1, 2}
	for i, mapping := range mappings {
		if mapping.To != expectedOrder[i] {
			t.Errorf("Position %d: got %s, want %s", i, mapping.To, expectedOrder[i])
		}
		if mapping.Position != expectedPositions[i] {
			t.Errorf("Mapping %d: got position %d, want %d", i, mapping.Position, expectedPositions[i])
		}
	}
}

func TestSortFieldMappingsByPosition_AlreadyOrdered(t *testing.T) {
	// Simulate API response already in correct order
	mappings := []client.FieldMapping{
		{To: "Email", Position: 0},
		{To: "FirstName", Position: 1},
		{To: "LastName", Position: 2},
	}

	// Sort by position (should be stable)
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].Position < mappings[j].Position
	})

	// Verify order is unchanged
	expectedOrder := []string{"Email", "FirstName", "LastName"}
	for i, mapping := range mappings {
		if mapping.To != expectedOrder[i] {
			t.Errorf("Position %d: got %s, want %s", i, mapping.To, expectedOrder[i])
		}
	}
}

func TestSortFieldMappingsByPosition_MixedMappingTypes(t *testing.T) {
	// Test with various mapping types in wrong order
	mappings := []client.FieldMapping{
		{
			To:             "FullName",
			Position:       3,
			Type:           "liquid_template",
			LiquidTemplate: "{{ first_name }} {{ last_name }}",
		},
		{
			From:                "email",
			To:                  "Email",
			Position:            0,
			Type:                "direct",
			IsPrimaryIdentifier: true,
		},
		{
			To:              "SyncRunID",
			Position:        2,
			Type:            "sync_metadata",
			SyncMetadataKey: "sync_run_id",
		},
		{
			To:       "Source",
			Position: 1,
			Type:     "constant",
			Constant: "Website",
		},
	}

	// Sort by position
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].Position < mappings[j].Position
	})

	// Verify order matches positions
	expectedOrder := []string{"Email", "Source", "SyncRunID", "FullName"}
	for i, mapping := range mappings {
		if mapping.To != expectedOrder[i] {
			t.Errorf("Position %d: got %s, want %s", i, mapping.To, expectedOrder[i])
		}
	}

	// Verify types are preserved
	if mappings[0].Type != "direct" || !mappings[0].IsPrimaryIdentifier {
		t.Errorf("Email mapping should be direct type with primary identifier")
	}
	if mappings[1].Type != "constant" {
		t.Errorf("Source mapping should be constant type")
	}
	if mappings[2].Type != "sync_metadata" {
		t.Errorf("SyncRunID mapping should be sync_metadata type")
	}
	if mappings[3].Type != "liquid_template" {
		t.Errorf("FullName mapping should be liquid_template type")
	}
}

func TestSortFieldMappingsByPosition_EmptySlice(t *testing.T) {
	// Empty slice
	mappings := []client.FieldMapping{}

	// Sort should not panic
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].Position < mappings[j].Position
	})

	// Should still be empty
	if len(mappings) != 0 {
		t.Errorf("Expected 0 mappings, got %d", len(mappings))
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
	if result == nil {
		t.Errorf("ExpandAlerts(empty) should return empty slice, got nil")
	}
	if len(result) != 0 {
		t.Errorf("ExpandAlerts(empty) should return empty slice with 0 items, got %d items", len(result))
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

// ============================================================================
// State Migration Tests
// ============================================================================
// These tests verify that state migration from v0 (TypeSet) to v1 (TypeList)
// works correctly for the alert field.

func TestResourceSyncStateUpgradeV0(t *testing.T) {
	// Simulate v0 state with alerts stored as they would be in TypeSet format
	// (which is actually just a JSON array, same as TypeList)
	v0State := map[string]interface{}{
		"id":           "12345",
		"label":        "Test Sync",
		"workspace_id": "67890",
		"alert": []interface{}{
			map[string]interface{}{
				"id":                   754756,
				"type":                 "FailureAlertConfiguration",
				"send_for":             "first_time",
				"should_send_recovery": true,
				"options":              map[string]interface{}{},
			},
			map[string]interface{}{
				"id":                   754757,
				"type":                 "InvalidRecordPercentAlertConfiguration",
				"send_for":             "first_time",
				"should_send_recovery": true,
				"options": map[string]interface{}{
					"threshold": "75",
				},
			},
		},
	}

	// Call the upgrade function
	upgradedState, err := provider.ResourceSyncStateUpgradeV0(nil, v0State, nil)

	// Assertions
	if err != nil {
		t.Errorf("State upgrade should not return an error, got: %v", err)
	}
	if upgradedState == nil {
		t.Errorf("Upgraded state should not be nil")
	}

	// Verify that the state is unchanged (since JSON format is identical)
	if !reflect.DeepEqual(v0State, upgradedState) {
		t.Errorf("State should be unchanged - TypeSet and TypeList have same JSON format")
	}

	// Verify alert data is preserved
	alerts, ok := upgradedState["alert"].([]interface{})
	if !ok {
		t.Errorf("alert should be a []interface{}")
	}
	if len(alerts) != 2 {
		t.Errorf("Should have 2 alerts, got %d", len(alerts))
	}

	// Verify first alert
	alert1, ok := alerts[0].(map[string]interface{})
	if !ok {
		t.Errorf("alert[0] should be a map")
	}
	if alert1["type"] != "FailureAlertConfiguration" {
		t.Errorf("alert[0] type = %v, want FailureAlertConfiguration", alert1["type"])
	}
	if alert1["id"] != 754756 {
		t.Errorf("alert[0] id = %v, want 754756", alert1["id"])
	}

	// Verify second alert
	alert2, ok := alerts[1].(map[string]interface{})
	if !ok {
		t.Errorf("alert[1] should be a map")
	}
	if alert2["type"] != "InvalidRecordPercentAlertConfiguration" {
		t.Errorf("alert[1] type = %v, want InvalidRecordPercentAlertConfiguration", alert2["type"])
	}
	if alert2["id"] != 754757 {
		t.Errorf("alert[1] id = %v, want 754757", alert2["id"])
	}
}

func TestResourceSyncV0SchemaCompatibility(t *testing.T) {
	// Verify that v0 schema uses TypeSet for alerts
	v0Resource := provider.ResourceSyncV0()
	if v0Resource == nil {
		t.Errorf("v0 resource should not be nil")
		return
	}
	if v0Resource.Schema == nil {
		t.Errorf("v0 schema should not be nil")
		return
	}

	alertSchema, exists := v0Resource.Schema["alert"]
	if !exists {
		t.Errorf("alert field should exist in v0 schema")
		return
	}

	// Note: We can't directly test schema.TypeSet value since it's an internal constant
	// But we can verify the schema is valid and has the expected structure
	if alertSchema == nil {
		t.Errorf("alert schema should not be nil")
	}
	if !alertSchema.Optional {
		t.Errorf("alert should be optional")
	}
}

func TestResourceSyncV1SchemaCompatibility(t *testing.T) {
	// Verify that v1 schema (current) uses TypeList for alerts
	v1Resource := provider.ResourceSync()
	if v1Resource == nil {
		t.Errorf("v1 resource should not be nil")
		return
	}
	if v1Resource.SchemaVersion != 1 {
		t.Errorf("SchemaVersion = %d, want 1", v1Resource.SchemaVersion)
	}
	if len(v1Resource.StateUpgraders) != 1 {
		t.Errorf("Should have 1 StateUpgrader, got %d", len(v1Resource.StateUpgraders))
		return
	}

	upgrader := v1Resource.StateUpgraders[0]
	if upgrader.Version != 0 {
		t.Errorf("StateUpgrader Version = %d, want 0", upgrader.Version)
	}
	if upgrader.Upgrade == nil {
		t.Errorf("Upgrade function should be set")
	}
	// Note: Type is a cty.Type struct (not a pointer), so we can't check for nil
	// The existence of the StateUpgrader itself confirms it's properly configured
}

// ============================================================================
// Mapping Conversion Tests (FieldMapping → MappingAttributes)
// ============================================================================

func TestConvertFieldMappingsToMappingAttributes_Constant(t *testing.T) {
	tests := []struct {
		name     string
		input    []client.FieldMapping
		expected []client.MappingAttributes
	}{
		{
			name: "single constant string mapping",
			input: []client.FieldMapping{
				{
					To:       "Source",
					Type:     "constant",
					Constant: "Website",
				},
			},
			expected: []client.MappingAttributes{
				{
					From: client.MappingFrom{
						Type: "constant_value",
						Data: map[string]interface{}{
							"basic_type": "text",
							"value":      "Website",
						},
					},
					To:                  "Source",
					IsPrimaryIdentifier: false,
				},
			},
		},
		{
			name: "multiple constant mappings with different to fields",
			input: []client.FieldMapping{
				{
					To:       "AssistantName",
					Type:     "constant",
					Constant: "HERE is my constant value",
				},
				{
					To:       "Title",
					Type:     "constant",
					Constant: "HERE is my constant value 2",
				},
			},
			expected: []client.MappingAttributes{
				{
					From: client.MappingFrom{
						Type: "constant_value",
						Data: map[string]interface{}{
							"basic_type": "text",
							"value":      "HERE is my constant value",
						},
					},
					To:                  "AssistantName",
					IsPrimaryIdentifier: false,
				},
				{
					From: client.MappingFrom{
						Type: "constant_value",
						Data: map[string]interface{}{
							"basic_type": "text",
							"value":      "HERE is my constant value 2",
						},
					},
					To:                  "Title",
					IsPrimaryIdentifier: false,
				},
			},
		},
		{
			name: "constant numeric value",
			input: []client.FieldMapping{
				{
					To:       "Priority",
					Type:     "constant",
					Constant: 1,
				},
			},
			expected: []client.MappingAttributes{
				{
					From: client.MappingFrom{
						Type: "constant_value",
						Data: map[string]interface{}{
							"basic_type": "text",
							"value":      "1",
						},
					},
					To:                  "Priority",
					IsPrimaryIdentifier: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ConvertFieldMappingsToMappingAttributes(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("convertFieldMappingsToMappingAttributes() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestConvertFieldMappingsToMappingAttributes_Direct(t *testing.T) {
	tests := []struct {
		name     string
		input    []client.FieldMapping
		expected []client.MappingAttributes
	}{
		{
			name: "direct mapping",
			input: []client.FieldMapping{
				{
					From: "user_id",
					To:   "UserID",
					Type: "direct",
				},
			},
			expected: []client.MappingAttributes{
				{
					From: client.MappingFrom{
						Type: "column",
						Data: "user_id", // Census API expects just the column name as string
					},
					To:                  "UserID",
					IsPrimaryIdentifier: false,
				},
			},
		},
		{
			name: "direct mapping with primary identifier",
			input: []client.FieldMapping{
				{
					From:                "user_id",
					To:                  "UserID",
					Type:                "direct",
					IsPrimaryIdentifier: true,
				},
			},
			expected: []client.MappingAttributes{
				{
					From: client.MappingFrom{
						Type: "column",
						Data: "user_id", // Census API expects just the column name as string
					},
					To:                  "UserID",
					IsPrimaryIdentifier: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.ConvertFieldMappingsToMappingAttributes(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("convertFieldMappingsToMappingAttributes() got = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestConvertFieldMappingsToMappingAttributes_Mixed(t *testing.T) {
	// Test a realistic scenario with multiple mapping types
	input := []client.FieldMapping{
		{
			From:                "user_id",
			To:                  "UserID",
			Type:                "direct",
			IsPrimaryIdentifier: true,
		},
		{
			From: "email",
			To:   "Email",
			Type: "direct",
		},
		{
			To:       "Source",
			Type:     "constant",
			Constant: "Website",
		},
		{
			To:              "SyncTime",
			Type:            "sync_metadata",
			SyncMetadataKey: "sync_run_id",
		},
	}

	expected := []client.MappingAttributes{
		{
			From: client.MappingFrom{
				Type: "column",
				Data: "user_id", // Census API expects string for column
			},
			To:                  "UserID",
			IsPrimaryIdentifier: true,
		},
		{
			From: client.MappingFrom{
				Type: "column",
				Data: "email", // Census API expects string for column
			},
			To:                  "Email",
			IsPrimaryIdentifier: false,
		},
		{
			From: client.MappingFrom{
				Type: "constant_value",
				Data: map[string]interface{}{
					"basic_type": "text",
					"value":      "Website",
				},
			},
			To:                  "Source",
			IsPrimaryIdentifier: false,
		},
		{
			From: client.MappingFrom{
				Type: "sync_metadata",
				Data: "sync_run_id", // Census API expects string for sync_metadata
			},
			To:                  "SyncTime",
			IsPrimaryIdentifier: false,
		},
	}

	result := provider.ConvertFieldMappingsToMappingAttributes(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("convertFieldMappingsToMappingAttributes() got = %+v, want %+v", result, expected)
	}
}

func TestConvertFieldMappingsToMappingAttributes_LiquidTemplate(t *testing.T) {
	// Test liquid template format
	input := []client.FieldMapping{
		{
			To:             "FormattedDate",
			Type:           "liquid_template",
			LiquidTemplate: "{{ record['date'] | upcase }}",
		},
	}

	expected := []client.MappingAttributes{
		{
			From: client.MappingFrom{
				Type: "liquid_template",
				Data: "{{ record['date'] | upcase }}", // Census API expects string directly, not wrapped in hash
			},
			To:                  "FormattedDate",
			IsPrimaryIdentifier: false,
		},
	}

	result := provider.ConvertFieldMappingsToMappingAttributes(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("convertFieldMappingsToMappingAttributes() got = %+v, want %+v", result, expected)
	}
}

// ============================================================================
// Reverse Mapping Conversion Tests (MappingAttributes → FieldMapping)
// ============================================================================

func TestConvertMappingAttributesToFieldMappings_LiquidTemplate(t *testing.T) {
	// Test that we correctly extract liquid template from API response
	// The Census API returns liquid templates as: {"liquid_template": "..."}
	input := []client.MappingAttributes{
		{
			From: client.MappingFrom{
				Type: "liquid_template",
				Data: map[string]interface{}{
					"liquid_template": "{{ record['date'] | upcase }}",
				},
			},
			To:                  "FormattedDate",
			IsPrimaryIdentifier: false,
			SyncNullValues:      boolPtr(true),
		},
	}

	expected := []client.FieldMapping{
		{
			From:                "",
			To:                  "FormattedDate",
			Type:                "liquid_template",
			LiquidTemplate:      "{{ record['date'] | upcase }}",
			IsPrimaryIdentifier: false,
			SyncNullValues:      boolPtr(true),
		},
	}

	result := provider.ConvertMappingAttributesToFieldMappings(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ConvertMappingAttributesToFieldMappings() got = %+v, want %+v", result, expected)
	}
}

func TestConvertMappingAttributesToFieldMappings_Constant(t *testing.T) {
	// Test constant value extraction from API
	input := []client.MappingAttributes{
		{
			From: client.MappingFrom{
				Type: "constant_value",
				Data: map[string]interface{}{
					"value":      "Website",
					"basic_type": "text",
				},
			},
			To:                  "Source",
			IsPrimaryIdentifier: false,
			SyncNullValues:      boolPtr(true),
		},
	}

	expected := []client.FieldMapping{
		{
			From:                "",
			To:                  "Source",
			Type:                "constant",
			Constant:            "Website",
			IsPrimaryIdentifier: false,
			SyncNullValues:      boolPtr(true),
		},
	}

	result := provider.ConvertMappingAttributesToFieldMappings(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ConvertMappingAttributesToFieldMappings() got = %+v, want %+v", result, expected)
	}
}

func TestConvertMappingAttributesToFieldMappings_SyncMetadata(t *testing.T) {
	// Test sync_metadata extraction from API (returned as string)
	input := []client.MappingAttributes{
		{
			From: client.MappingFrom{
				Type: "sync_metadata",
				Data: "sync_run_id",
			},
			To:                  "SyncRunID",
			IsPrimaryIdentifier: false,
			SyncNullValues:      boolPtr(true),
		},
	}

	expected := []client.FieldMapping{
		{
			From:                "",
			To:                  "SyncRunID",
			Type:                "sync_metadata",
			SyncMetadataKey:     "sync_run_id",
			IsPrimaryIdentifier: false,
			SyncNullValues:      boolPtr(true),
		},
	}

	result := provider.ConvertMappingAttributesToFieldMappings(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ConvertMappingAttributesToFieldMappings() got = %+v, want %+v", result, expected)
	}
}

func TestConvertMappingAttributesToFieldMappings_SegmentMembership(t *testing.T) {
	// Test segment_membership extraction from API
	input := []client.MappingAttributes{
		{
			From: client.MappingFrom{
				Type: "segment_membership",
				Data: map[string]interface{}{
					"identify_by": "name",
				},
			},
			To:                  "SegmentField",
			IsPrimaryIdentifier: false,
			SyncNullValues:      boolPtr(true),
		},
	}

	expected := []client.FieldMapping{
		{
			From:                "",
			To:                  "SegmentField",
			Type:                "segment_membership",
			SegmentIdentifyBy:   "name",
			IsPrimaryIdentifier: false,
			SyncNullValues:      boolPtr(true),
		},
	}

	result := provider.ConvertMappingAttributesToFieldMappings(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ConvertMappingAttributesToFieldMappings() got = %+v, want %+v", result, expected)
	}
}

func TestConvertMappingAttributesToFieldMappings_RoundTrip(t *testing.T) {
	// Test full round-trip: FieldMapping → MappingAttributes → FieldMapping
	original := []client.FieldMapping{
		{
			From:                "email",
			To:                  "Email",
			Type:                "direct",
			IsPrimaryIdentifier: true,
			SyncNullValues:      boolPtr(true),
		},
		{
			To:             "FormattedDate",
			Type:           "liquid_template",
			LiquidTemplate: "{{ record['date'] | upcase }}",
			SyncNullValues: boolPtr(true),
		},
		{
			To:              "SyncTime",
			Type:            "sync_metadata",
			SyncMetadataKey: "sync_run_id",
			SyncNullValues:  boolPtr(true),
		},
		{
			To:             "Source",
			Type:           "constant",
			Constant:       "Website",
			SyncNullValues: boolPtr(true),
		},
	}

	// Convert to MappingAttributes (what we send to API)
	mappingAttrs := provider.ConvertFieldMappingsToMappingAttributes(original)

	// Convert back to FieldMapping (what we read from API)
	result := provider.ConvertMappingAttributesToFieldMappings(mappingAttrs)

	// Should match original (with "from" field empty for non-direct mappings)
	expected := []client.FieldMapping{
		{
			From:                "email",
			To:                  "Email",
			Type:                "direct",
			IsPrimaryIdentifier: true,
			SyncNullValues:      boolPtr(true),
		},
		{
			From:           "",
			To:             "FormattedDate",
			Type:           "liquid_template",
			LiquidTemplate: "{{ record['date'] | upcase }}",
			SyncNullValues: boolPtr(true),
		},
		{
			From:            "",
			To:              "SyncTime",
			Type:            "sync_metadata",
			SyncMetadataKey: "sync_run_id",
			SyncNullValues:  boolPtr(true),
		},
		{
			From:           "",
			To:             "Source",
			Type:           "constant",
			Constant:       "Website",
			SyncNullValues: boolPtr(true),
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Round-trip conversion failed.\nGot:      %+v\nExpected: %+v", result, expected)
	}
}

// Helper function for tests
func boolPtr(b bool) *bool {
	return &b
}

// ============================================================================
// Compound Key Mapping Tests
// ============================================================================
// These tests verify that Census-managed compound_key mappings are properly
// handled and filtered out to prevent drift in Terraform state.

func TestConvertMappingAttributesToFieldMappings_CompoundKey(t *testing.T) {
	// Test that compound_key type from API is properly converted
	// Census auto-generates these for destinations like Google Ads with compound keys
	input := []client.MappingAttributes{
		{
			From: client.MappingFrom{
				Type: "compound_key",
				Data: []interface{}{
					map[string]interface{}{
						"type": "column",
						"data": "ID",
					},
					map[string]interface{}{
						"type": "column",
						"data": "ID",
					},
					map[string]interface{}{
						"type": "column",
						"data": "ID",
					},
				},
			},
			To:                  "click_conversion._census_tracking_id",
			IsPrimaryIdentifier: true,
			SyncNullValues:      boolPtr(true),
		},
	}

	expected := []client.FieldMapping{
		{
			From:                "[map[data:ID type:column] map[data:ID type:column] map[data:ID type:column]]",
			To:                  "click_conversion._census_tracking_id",
			Type:                "compound_key",
			IsPrimaryIdentifier: true,
			SyncNullValues:      boolPtr(true),
		},
	}

	result := provider.ConvertMappingAttributesToFieldMappings(input)
	if len(result) != len(expected) {
		t.Errorf("Expected %d mappings, got %d", len(expected), len(result))
	}
	if result[0].Type != "compound_key" {
		t.Errorf("Expected Type='compound_key', got '%s'", result[0].Type)
	}
	if result[0].To != expected[0].To {
		t.Errorf("Expected To='%s', got '%s'", expected[0].To, result[0].To)
	}
	if result[0].IsPrimaryIdentifier != expected[0].IsPrimaryIdentifier {
		t.Errorf("Expected IsPrimaryIdentifier=%v, got %v", expected[0].IsPrimaryIdentifier, result[0].IsPrimaryIdentifier)
	}
}

func TestConvertMappingAttributesToFieldMappings_WithCompoundKey(t *testing.T) {
	// Test conversion of mixed mappings including a Census-managed compound_key
	// This simulates the real API response for a Google Ads sync
	input := []client.MappingAttributes{
		{
			From: client.MappingFrom{
				Type: "column",
				Data: "ID",
			},
			To:                  "click_conversion.conversion_date_time",
			IsPrimaryIdentifier: false,
			Position:            0,
		},
		{
			From: client.MappingFrom{
				Type: "column",
				Data: "ID",
			},
			To:                  "click_conversion.gclid",
			IsPrimaryIdentifier: false,
			Position:            1,
		},
		{
			From: client.MappingFrom{
				Type: "column",
				Data: "ID",
			},
			To:                  "click_conversion.conversion_action",
			IsPrimaryIdentifier: false,
			Position:            2,
		},
		{
			From: client.MappingFrom{
				Type: "compound_key",
				Data: []interface{}{
					map[string]interface{}{
						"type": "column",
						"data": "ID",
					},
					map[string]interface{}{
						"type": "column",
						"data": "ID",
					},
					map[string]interface{}{
						"type": "column",
						"data": "ID",
					},
				},
			},
			To:                  "click_conversion._census_tracking_id",
			IsPrimaryIdentifier: true,
			Position:            3,
		},
	}

	result := provider.ConvertMappingAttributesToFieldMappings(input)

	// Verify we have all 4 mappings
	if len(result) != 4 {
		t.Errorf("Expected 4 mappings, got %d", len(result))
	}

	// Verify the compound_key mapping is at the end
	compoundKeyMapping := result[3]
	if compoundKeyMapping.Type != "compound_key" {
		t.Errorf("Expected last mapping Type='compound_key', got '%s'", compoundKeyMapping.Type)
	}
	if compoundKeyMapping.To != "click_conversion._census_tracking_id" {
		t.Errorf("Expected To='click_conversion._census_tracking_id', got '%s'", compoundKeyMapping.To)
	}
	if !compoundKeyMapping.IsPrimaryIdentifier {
		t.Errorf("Expected compound_key mapping to be primary identifier")
	}

	// Verify the user-defined mappings are correct
	if result[0].Type != "direct" {
		t.Errorf("Expected mapping 0 Type='direct', got '%s'", result[0].Type)
	}
	if result[0].From != "ID" {
		t.Errorf("Expected mapping 0 From='ID', got '%s'", result[0].From)
	}
}

func TestFilterCompoundKeyMappings(t *testing.T) {
	// Test the filtering logic that removes Census-managed compound_key mappings
	// This simulates what happens in resourceSyncRead
	tests := []struct {
		name     string
		input    []client.FieldMapping
		expected int // Number of mappings after filtering
	}{
		{
			name: "filter out single compound_key mapping",
			input: []client.FieldMapping{
				{
					From:                "ID",
					To:                  "click_conversion.conversion_date_time",
					Type:                "direct",
					IsPrimaryIdentifier: false,
					Position:            0,
				},
				{
					From:                "ID",
					To:                  "click_conversion.gclid",
					Type:                "direct",
					IsPrimaryIdentifier: false,
					Position:            1,
				},
				{
					From:                "[map[data:ID type:column] map[data:ID type:column] map[data:ID type:column]]",
					To:                  "click_conversion._census_tracking_id",
					Type:                "compound_key",
					IsPrimaryIdentifier: true,
					Position:            2,
				},
			},
			expected: 2, // Only the 2 user-defined mappings
		},
		{
			name: "no compound_key mappings to filter",
			input: []client.FieldMapping{
				{
					From:                "email",
					To:                  "Email",
					Type:                "direct",
					IsPrimaryIdentifier: true,
					Position:            0,
				},
				{
					From:     "first_name",
					To:       "FirstName",
					Type:     "direct",
					Position: 1,
				},
			},
			expected: 2, // All mappings remain
		},
		{
			name:     "empty input",
			input:    []client.FieldMapping{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the filtering logic from resourceSyncRead
			userManagedMappings := make([]client.FieldMapping, 0, len(tt.input))
			for _, mapping := range tt.input {
				if mapping.Type != "compound_key" {
					userManagedMappings = append(userManagedMappings, mapping)
				}
			}

			if len(userManagedMappings) != tt.expected {
				t.Errorf("Expected %d mappings after filtering, got %d", tt.expected, len(userManagedMappings))
			}

			// Verify no compound_key mappings remain
			for i, mapping := range userManagedMappings {
				if mapping.Type == "compound_key" {
					t.Errorf("Mapping %d should not be compound_key type", i)
				}
			}
		})
	}
}

func TestFilterCompoundKeyMappings_PreservesOtherMappingTypes(t *testing.T) {
	// Verify that filtering compound_key doesn't affect other mapping types
	input := []client.FieldMapping{
		{
			From:                "email",
			To:                  "Email",
			Type:                "direct",
			IsPrimaryIdentifier: true,
			Position:            0,
		},
		{
			To:       "Source",
			Type:     "constant",
			Constant: "Website",
			Position: 1,
		},
		{
			To:              "SyncRunID",
			Type:            "sync_metadata",
			SyncMetadataKey: "sync_run_id",
			Position:        2,
		},
		{
			To:             "FullName",
			Type:           "liquid_template",
			LiquidTemplate: "{{ first_name }} {{ last_name }}",
			Position:       3,
		},
		{
			From:                "[map[data:ID type:column]]",
			To:                  "_census_tracking_id",
			Type:                "compound_key",
			IsPrimaryIdentifier: true,
			Position:            4,
		},
	}

	// Filter out compound_key mappings
	userManagedMappings := make([]client.FieldMapping, 0, len(input))
	for _, mapping := range input {
		if mapping.Type != "compound_key" {
			userManagedMappings = append(userManagedMappings, mapping)
		}
	}

	// Should have 4 mappings (all except compound_key)
	if len(userManagedMappings) != 4 {
		t.Errorf("Expected 4 mappings after filtering, got %d", len(userManagedMappings))
	}

	// Verify all mapping types are preserved
	expectedTypes := []string{"direct", "constant", "sync_metadata", "liquid_template"}
	for i, mapping := range userManagedMappings {
		if mapping.Type != expectedTypes[i] {
			t.Errorf("Mapping %d: expected type '%s', got '%s'", i, expectedTypes[i], mapping.Type)
		}
	}

	// Verify the primary identifier is preserved
	if !userManagedMappings[0].IsPrimaryIdentifier {
		t.Errorf("Expected first mapping to be primary identifier")
	}
}
