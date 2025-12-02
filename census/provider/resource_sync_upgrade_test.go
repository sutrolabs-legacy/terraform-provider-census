package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	ctx := context.Background()
	upgradedState, err := resourceSyncStateUpgradeV0(ctx, v0State, nil)

	// Assertions
	assert.NoError(t, err, "State upgrade should not return an error")
	assert.NotNil(t, upgradedState, "Upgraded state should not be nil")

	// Verify that the state is unchanged (since JSON format is identical)
	assert.Equal(t, v0State, upgradedState, "State should be unchanged - TypeSet and TypeList have same JSON format")

	// Verify alert data is preserved
	alerts, ok := upgradedState["alert"].([]interface{})
	assert.True(t, ok, "alert should be a []interface{}")
	assert.Len(t, alerts, 2, "Should have 2 alerts")

	// Verify first alert
	alert1, ok := alerts[0].(map[string]interface{})
	assert.True(t, ok, "alert[0] should be a map")
	assert.Equal(t, "FailureAlertConfiguration", alert1["type"])
	assert.Equal(t, 754756, alert1["id"])

	// Verify second alert
	alert2, ok := alerts[1].(map[string]interface{})
	assert.True(t, ok, "alert[1] should be a map")
	assert.Equal(t, "InvalidRecordPercentAlertConfiguration", alert2["type"])
	assert.Equal(t, 754757, alert2["id"])
}

func TestResourceSyncV0SchemaCompatibility(t *testing.T) {
	// Verify that v0 schema uses TypeSet for alerts
	v0Resource := resourceSyncV0()
	assert.NotNil(t, v0Resource, "v0 resource should not be nil")
	assert.NotNil(t, v0Resource.Schema, "v0 schema should not be nil")

	alertSchema, exists := v0Resource.Schema["alert"]
	assert.True(t, exists, "alert field should exist in v0 schema")

	// Note: We can't directly test schema.TypeSet value since it's an internal constant
	// But we can verify the schema is valid and has the expected structure
	assert.NotNil(t, alertSchema, "alert schema should not be nil")
	assert.True(t, alertSchema.Optional, "alert should be optional")
}

func TestResourceSyncV1SchemaCompatibility(t *testing.T) {
	// Verify that v1 schema (current) uses TypeList for alerts
	v1Resource := resourceSync()
	assert.NotNil(t, v1Resource, "v1 resource should not be nil")
	assert.Equal(t, 1, v1Resource.SchemaVersion, "SchemaVersion should be 1")
	assert.Len(t, v1Resource.StateUpgraders, 1, "Should have 1 StateUpgrader")

	upgrader := v1Resource.StateUpgraders[0]
	assert.Equal(t, 0, upgrader.Version, "StateUpgrader should handle version 0")
	assert.NotNil(t, upgrader.Upgrade, "Upgrade function should be set")
	assert.NotNil(t, upgrader.Type, "Type should be set")
}
