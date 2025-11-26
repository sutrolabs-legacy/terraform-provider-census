package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sutrolabs/terraform-provider-census/census/client"
)

// TestValidateSourceCredentials_SnowflakePassword tests Snowflake password authentication validation
func TestValidateSourceCredentials_SnowflakePassword(t *testing.T) {
	// Mock source types response
	sourceTypes := map[string]interface{}{
		"status": "success",
		"data": []interface{}{
			map[string]interface{}{
				"service_name": "snowflake",
				"configuration_fields": map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"id":    "account",
							"label": "Account",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
						map[string]interface{}{
							"id":    "username",
							"label": "Username",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
						map[string]interface{}{
							"id":                     "password",
							"label":                  "Password",
							"type":                   "string",
							"rules":                  []interface{}{"required"},
							"is_password_type_field": true,
							"show": map[string]interface{}{
								"unless": map[string]interface{}{
									"use_keypair": map[string]interface{}{
										"eq": true,
									},
								},
							},
						},
						map[string]interface{}{
							"id":    "database",
							"label": "Database",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
						map[string]interface{}{
							"id":    "warehouse",
							"label": "Warehouse",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/source_types" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sourceTypes)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	c, err := client.NewClient(&client.Config{
		PersonalAccessToken: "test-token",
		BaseURL:             server.URL,
		Region:              "us",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test valid password-based credentials
	credentials := map[string]interface{}{
		"account":   "test.us-east-1",
		"username":  "census",
		"password":  "secret123",
		"database":  "PRODUCTION",
		"warehouse": "COMPUTE_WH",
	}

	err = c.ValidateSourceCredentials(context.Background(), "snowflake", credentials, "test-workspace-token")
	if err != nil {
		t.Errorf("Expected no error for valid password credentials, got: %v", err)
	}

	// Test missing password (should fail)
	invalidCredentials := map[string]interface{}{
		"account":   "test.us-east-1",
		"username":  "census",
		"database":  "PRODUCTION",
		"warehouse": "COMPUTE_WH",
	}

	err = c.ValidateSourceCredentials(context.Background(), "snowflake", invalidCredentials, "test-workspace-token")
	if err == nil {
		t.Error("Expected error for missing password, got nil")
	}
}

// TestValidateSourceCredentials_SnowflakeKeypair tests Snowflake keypair authentication validation
func TestValidateSourceCredentials_SnowflakeKeypair(t *testing.T) {
	// Mock source types response with unless condition
	sourceTypes := map[string]interface{}{
		"status": "success",
		"data": []interface{}{
			map[string]interface{}{
				"service_name": "snowflake",
				"configuration_fields": map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"id":    "account",
							"label": "Account",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
						map[string]interface{}{
							"id":    "username",
							"label": "Username",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
						map[string]interface{}{
							"id":                     "password",
							"label":                  "Password",
							"type":                   "string",
							"rules":                  []interface{}{"required"},
							"is_password_type_field": true,
							"show": map[string]interface{}{
								"unless": map[string]interface{}{
									"use_keypair": map[string]interface{}{
										"eq": true,
									},
								},
							},
						},
						map[string]interface{}{
							"id":    "use_keypair",
							"label": "Use Keypair",
							"type":  "boolean",
						},
						map[string]interface{}{
							"id":    "private_key_pkcs8",
							"label": "Private Key (PKCS8)",
							"type":  "string",
						},
						map[string]interface{}{
							"id":    "database",
							"label": "Database",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
						map[string]interface{}{
							"id":    "warehouse",
							"label": "Warehouse",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/source_types" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sourceTypes)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	c, err := client.NewClient(&client.Config{
		PersonalAccessToken: "test-token",
		BaseURL:             server.URL,
		Region:              "us",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test valid keypair credentials (use_keypair=true, no password required)
	credentials := map[string]interface{}{
		"account":           "test.us-east-1",
		"username":          "census",
		"use_keypair":       true,
		"private_key_pkcs8": "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----",
		"database":          "PRODUCTION",
		"warehouse":         "COMPUTE_WH",
	}

	err = c.ValidateSourceCredentials(context.Background(), "snowflake", credentials, "test-workspace-token")
	if err != nil {
		t.Errorf("Expected no error for valid keypair credentials, got: %v", err)
	}

	// Test keypair with boolean as string
	credentialsStringBool := map[string]interface{}{
		"account":           "test.us-east-1",
		"username":          "census",
		"use_keypair":       "true",
		"private_key_pkcs8": "-----BEGIN PRIVATE KEY-----\ntest\n-----END PRIVATE KEY-----",
		"database":          "PRODUCTION",
		"warehouse":         "COMPUTE_WH",
	}

	err = c.ValidateSourceCredentials(context.Background(), "snowflake", credentialsStringBool, "test-workspace-token")
	if err != nil {
		t.Errorf("Expected no error for valid keypair credentials with string boolean, got: %v", err)
	}

	// Test keypair=false should still require password
	credentialsKeypairFalse := map[string]interface{}{
		"account":     "test.us-east-1",
		"username":    "census",
		"use_keypair": false,
		"database":    "PRODUCTION",
		"warehouse":   "COMPUTE_WH",
	}

	err = c.ValidateSourceCredentials(context.Background(), "snowflake", credentialsKeypairFalse, "test-workspace-token")
	if err == nil {
		t.Error("Expected error when use_keypair=false and password missing, got nil")
	}
}

// TestValidateSourceCredentials_ShowIfCondition tests show.if conditional logic
func TestValidateSourceCredentials_ShowIfCondition(t *testing.T) {
	// Mock source types with "show.if" condition
	sourceTypes := map[string]interface{}{
		"status": "success",
		"data": []interface{}{
			map[string]interface{}{
				"service_name": "postgres",
				"configuration_fields": map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"id":    "host",
							"label": "Host",
							"type":  "string",
							"rules": []interface{}{"required"},
						},
						map[string]interface{}{
							"id":    "ssh_tunnel_enabled",
							"label": "SSH Tunnel Enabled",
							"type":  "boolean",
						},
						map[string]interface{}{
							"id":    "ssh_host",
							"label": "SSH Host",
							"type":  "string",
							"rules": []interface{}{"required"},
							"show": map[string]interface{}{
								"if": map[string]interface{}{
									"ssh_tunnel_enabled": true,
								},
							},
						},
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/source_types" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sourceTypes)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	c, err := client.NewClient(&client.Config{
		PersonalAccessToken: "test-token",
		BaseURL:             server.URL,
		Region:              "us",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with SSH tunnel disabled - ssh_host not required
	credentialsNoSSH := map[string]interface{}{
		"host":               "db.example.com",
		"ssh_tunnel_enabled": false,
	}

	err = c.ValidateSourceCredentials(context.Background(), "postgres", credentialsNoSSH, "test-workspace-token")
	if err != nil {
		t.Errorf("Expected no error when SSH tunnel disabled, got: %v", err)
	}

	// Test with SSH tunnel enabled - ssh_host required
	credentialsWithSSH := map[string]interface{}{
		"host":               "db.example.com",
		"ssh_tunnel_enabled": true,
		"ssh_host":           "tunnel.example.com",
	}

	err = c.ValidateSourceCredentials(context.Background(), "postgres", credentialsWithSSH, "test-workspace-token")
	if err != nil {
		t.Errorf("Expected no error when SSH tunnel enabled with ssh_host, got: %v", err)
	}

	// Test with SSH tunnel enabled but missing ssh_host - should fail
	credentialsSSHNoHost := map[string]interface{}{
		"host":               "db.example.com",
		"ssh_tunnel_enabled": true,
	}

	err = c.ValidateSourceCredentials(context.Background(), "postgres", credentialsSSHNoHost, "test-workspace-token")
	if err == nil {
		t.Error("Expected error when SSH tunnel enabled but ssh_host missing, got nil")
	}
}

// TestValidateSourceCredentials_UnknownSourceType tests validation with unknown source type
func TestValidateSourceCredentials_UnknownSourceType(t *testing.T) {
	sourceTypes := map[string]interface{}{
		"status": "success",
		"data": []interface{}{
			map[string]interface{}{
				"service_name": "snowflake",
				"configuration_fields": map[string]interface{}{
					"fields": []interface{}{},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/source_types" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sourceTypes)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	c, err := client.NewClient(&client.Config{
		PersonalAccessToken: "test-token",
		BaseURL:             server.URL,
		Region:              "us",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	credentials := map[string]interface{}{}

	err = c.ValidateSourceCredentials(context.Background(), "unknown_type", credentials, "test-workspace-token")
	if err == nil {
		t.Error("Expected error for unknown source type, got nil")
	}
	if err != nil && err.Error() != "unknown source type: unknown_type" {
		t.Errorf("Expected 'unknown source type' error, got: %v", err)
	}
}
