package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sutrolabs/terraform-provider-census/census/client"
)

func TestRefreshDatasetColumns(t *testing.T) {
	tests := []struct {
		name           string
		datasetID      int
		responseStatus int
		responseBody   map[string]interface{}
		wantRefreshKey int
		wantErr        bool
	}{
		{
			name:           "successful refresh",
			datasetID:      123,
			responseStatus: http.StatusAccepted,
			responseBody: map[string]interface{}{
				"refresh_key": 1234567890,
			},
			wantRefreshKey: 1234567890,
			wantErr:        false,
		},
		{
			name:           "api error",
			datasetID:      123,
			responseStatus: http.StatusBadRequest,
			responseBody: map[string]interface{}{
				"error": "Invalid dataset",
			},
			wantRefreshKey: 0,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				expectedPath := "/datasets/" + string(rune(tt.datasetID)) + "/refresh_columns"
				if r.URL.Path != expectedPath {
					t.Logf("Got path: %s, expected to contain dataset ID", r.URL.Path)
				}

				// Send response
				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Create client
			config := &client.Config{
				BaseURL:              server.URL,
				PersonalAccessToken:  "test-token",
				WorkspaceAccessToken: "test-workspace-token",
			}
			c, err := client.NewClient(config)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Call method
			ctx := context.Background()
			refreshKey, err := c.RefreshDatasetColumnsWithToken(ctx, tt.datasetID, "test-workspace-token")

			// Check results
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshDatasetColumnsWithToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && refreshKey != tt.wantRefreshKey {
				t.Errorf("RefreshDatasetColumnsWithToken() refreshKey = %v, want %v", refreshKey, tt.wantRefreshKey)
			}
		})
	}
}

func TestGetDatasetRefreshStatus(t *testing.T) {
	tests := []struct {
		name           string
		datasetID      int
		refreshKey     int
		responseStatus int
		responseBody   map[string]interface{}
		wantStatus     string
		wantMessage    *string
		wantErr        bool
	}{
		{
			name:           "processing status",
			datasetID:      123,
			refreshKey:     1234567890,
			responseStatus: http.StatusOK,
			responseBody: map[string]interface{}{
				"status": "processing",
			},
			wantStatus:  "processing",
			wantMessage: nil,
			wantErr:     false,
		},
		{
			name:           "completed status",
			datasetID:      123,
			refreshKey:     1234567890,
			responseStatus: http.StatusOK,
			responseBody: map[string]interface{}{
				"status": "completed",
			},
			wantStatus:  "completed",
			wantMessage: nil,
			wantErr:     false,
		},
		{
			name:           "error status with message",
			datasetID:      123,
			refreshKey:     1234567890,
			responseStatus: http.StatusOK,
			responseBody: map[string]interface{}{
				"status":  "error",
				"message": "Failed to refresh source columns",
			},
			wantStatus: "error",
			wantMessage: func() *string {
				s := "Failed to refresh source columns"
				return &s
			}(),
			wantErr: false,
		},
		{
			name:           "api error",
			datasetID:      123,
			refreshKey:     1234567890,
			responseStatus: http.StatusNotFound,
			responseBody: map[string]interface{}{
				"error": "Dataset not found",
			},
			wantStatus:  "",
			wantMessage: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				// Verify query parameter
				if r.URL.Query().Get("refresh_key") == "" {
					t.Error("Expected refresh_key query parameter")
				}

				// Send response
				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Create client
			config := &client.Config{
				BaseURL:              server.URL,
				PersonalAccessToken:  "test-token",
				WorkspaceAccessToken: "test-workspace-token",
			}
			c, err := client.NewClient(config)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Call method
			ctx := context.Background()
			statusResp, err := c.GetDatasetRefreshStatusWithToken(ctx, tt.datasetID, tt.refreshKey, "test-workspace-token")

			// Check results
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDatasetRefreshStatusWithToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if statusResp.Status != tt.wantStatus {
					t.Errorf("GetDatasetRefreshStatusWithToken() status = %v, want %v", statusResp.Status, tt.wantStatus)
				}

				if tt.wantMessage != nil {
					if statusResp.Message == nil {
						t.Error("Expected message but got nil")
					} else if *statusResp.Message != *tt.wantMessage {
						t.Errorf("GetDatasetRefreshStatusWithToken() message = %v, want %v", *statusResp.Message, *tt.wantMessage)
					}
				} else if statusResp.Message != nil {
					t.Errorf("Expected nil message but got %v", *statusResp.Message)
				}
			}
		})
	}
}

func TestGetDatasetRefreshStatus_Polling(t *testing.T) {
	// Test simulates polling behavior: processing -> processing -> completed
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)

		var response map[string]interface{}
		if callCount < 3 {
			response = map[string]interface{}{
				"status": "processing",
			}
		} else {
			response = map[string]interface{}{
				"status": "completed",
			}
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &client.Config{
		BaseURL:              server.URL,
		PersonalAccessToken:  "test-token",
		WorkspaceAccessToken: "test-workspace-token",
	}
	c, err := client.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Poll until completed
	for i := 0; i < 5; i++ {
		statusResp, err := c.GetDatasetRefreshStatusWithToken(ctx, 123, 1234567890, "test-workspace-token")
		if err != nil {
			t.Fatalf("Poll %d failed: %v", i, err)
		}

		if statusResp.Status == "completed" {
			if i < 2 {
				t.Errorf("Expected at least 3 polls before completion, got %d", i+1)
			}
			return
		}

		if statusResp.Status != "processing" {
			t.Errorf("Poll %d: unexpected status %s", i, statusResp.Status)
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Error("Expected to reach 'completed' status but didn't")
}
