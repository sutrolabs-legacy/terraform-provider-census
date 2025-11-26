package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Source represents a Census data source
type Source struct {
	ID          int                    `json:"id"`
	WorkspaceID string                 `json:"workspace_id,omitempty"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Connection  map[string]interface{} `json:"connection"`
	Status      string                 `json:"status,omitempty"`
	TestStatus  string                 `json:"test_status,omitempty"`
	LastTested  *time.Time             `json:"last_tested,omitempty"`
}

// CreateSourceRequest represents the request to create a source
type CreateSourceRequest struct {
	Connection SourceConnection `json:"connection"`
}

// SourceConnection represents the connection configuration for a source
type SourceConnection struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type,omitempty"`
	SyncEngine  string                 `json:"sync_engine,omitempty"`
	Credentials map[string]interface{} `json:"credentials"`
}

// UpdateSourceRequest represents the request to update a source
type UpdateSourceRequest struct {
	Connection SourceConnection `json:"connection,omitempty"`
}

// SourceResponse represents a single source response
type SourceResponse struct {
	Status string  `json:"status"`
	Data   *Source `json:"data"`
}

// SourceListResponse represents a paginated source list response
type SourceListResponse struct {
	Status     string         `json:"status"`
	Pagination PaginationInfo `json:"pagination"`
	Data       []Source       `json:"data"`
}

// SourceObject represents an object within a source (table, model, etc.)
type SourceObject struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	FullName string `json:"full_name,omitempty"`
}

// SourceObjectsResponse represents the response for source objects
type SourceObjectsResponse struct {
	Status string         `json:"status"`
	Data   []SourceObject `json:"data"`
}

// ConnectLink represents a connection link for OAuth/reauth
type ConnectLink struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ConnectLinkResponse represents the connect link response
type ConnectLinkResponse struct {
	Status string      `json:"status"`
	Data   ConnectLink `json:"data"`
}

// RefreshStatus represents the status of a refresh operation
type RefreshStatus struct {
	Status     string `json:"status"`
	InProgress bool   `json:"in_progress"`
}

// RefreshStatusResponse represents the response for refresh status
type RefreshStatusResponse struct {
	Status string        `json:"status"`
	Data   RefreshStatus `json:"data"`
}

// CreateSource creates a new source
func (c *Client) CreateSource(ctx context.Context, req *CreateSourceRequest) (*Source, error) {
	return c.CreateSourceWithToken(ctx, req, "")
}

// CreateSourceWithToken creates a new source using a specific workspace token
func (c *Client) CreateSourceWithToken(ctx context.Context, req *CreateSourceRequest, workspaceToken string) (*Source, error) {
	resp, err := c.makeRequestWithToken(ctx, http.MethodPost, "/sources", req, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make create source request: %w", err)
	}

	var result SourceResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	return result.Data, nil
}

// GetSource retrieves a source by ID
func (c *Client) GetSource(ctx context.Context, sourceID int) (*Source, error) {
	return c.GetSourceWithToken(ctx, sourceID, "")
}

// GetSourceWithToken retrieves a source by ID using a specific workspace token
func (c *Client) GetSourceWithToken(ctx context.Context, sourceID int, workspaceToken string) (*Source, error) {
	path := fmt.Sprintf("/sources/%d", sourceID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get source request: %w", err)
	}

	var result SourceResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	return result.Data, nil
}

// UpdateSource updates an existing source
func (c *Client) UpdateSource(ctx context.Context, sourceID int, req *UpdateSourceRequest) (*Source, error) {
	return c.UpdateSourceWithToken(ctx, sourceID, req, "")
}

// UpdateSourceWithToken updates an existing source using a specific workspace token
func (c *Client) UpdateSourceWithToken(ctx context.Context, sourceID int, req *UpdateSourceRequest, workspaceToken string) (*Source, error) {
	path := fmt.Sprintf("/sources/%d", sourceID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodPatch, path, req, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make update source request: %w", err)
	}

	var result SourceResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to update source: %w", err)
	}

	return result.Data, nil
}

// DeleteSource deletes a source by ID
func (c *Client) DeleteSource(ctx context.Context, sourceID int) error {
	return c.DeleteSourceWithToken(ctx, sourceID, "")
}

// DeleteSourceWithToken deletes a source by ID using a specific workspace token
func (c *Client) DeleteSourceWithToken(ctx context.Context, sourceID int, workspaceToken string) error {
	path := fmt.Sprintf("/sources/%d", sourceID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodDelete, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return fmt.Errorf("failed to make delete source request: %w", err)
	}

	if err := c.handleResponse(resp, nil); err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
	}

	return nil
}

// ListSources retrieves a list of sources
func (c *Client) ListSources(ctx context.Context, opts *ListOptions) ([]Source, *PaginationInfo, error) {
	return c.ListSourcesWithToken(ctx, opts, "")
}

// ListSourcesWithToken retrieves a list of sources using a specific workspace token
func (c *Client) ListSourcesWithToken(ctx context.Context, opts *ListOptions, workspaceToken string) ([]Source, *PaginationInfo, error) {
	params := make(map[string]string)
	if opts != nil {
		params = opts.ToParams()
	}

	fullURL := c.buildURL("/sources", params)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, fullURL, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make list sources request: %w", err)
	}

	var result SourceListResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to list sources: %w", err)
	}

	return result.Data, &result.Pagination, nil
}

// GetSourceObjects retrieves objects (tables, models, etc.) for a source
func (c *Client) GetSourceObjects(ctx context.Context, sourceID int) ([]SourceObject, error) {
	return c.GetSourceObjectsWithToken(ctx, sourceID, "")
}

// GetSourceObjectsWithToken retrieves objects (tables, models, etc.) for a source using a specific workspace token
func (c *Client) GetSourceObjectsWithToken(ctx context.Context, sourceID int, workspaceToken string) ([]SourceObject, error) {
	path := fmt.Sprintf("/sources/%d/objects", sourceID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get source objects request: %w", err)
	}

	var result SourceObjectsResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get source objects: %w", err)
	}

	return result.Data, nil
}

// CreateSourceConnectLink creates a connect link for source reauthorization
func (c *Client) CreateSourceConnectLink(ctx context.Context, sourceID int) (*ConnectLink, error) {
	return c.CreateSourceConnectLinkWithToken(ctx, sourceID, "")
}

// CreateSourceConnectLinkWithToken creates a connect link for source reauthorization using a specific workspace token
func (c *Client) CreateSourceConnectLinkWithToken(ctx context.Context, sourceID int, workspaceToken string) (*ConnectLink, error) {
	path := fmt.Sprintf("/sources/%d/connect_links", sourceID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodPost, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make create connect link request: %w", err)
	}

	var result ConnectLinkResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to create connect link: %w", err)
	}

	return &result.Data, nil
}

// RefreshSourceTables starts a table refresh for a source
func (c *Client) RefreshSourceTables(ctx context.Context, sourceID int) error {
	return c.RefreshSourceTablesWithToken(ctx, sourceID, "")
}

// RefreshSourceTablesWithToken starts a table refresh for a source using a specific workspace token
func (c *Client) RefreshSourceTablesWithToken(ctx context.Context, sourceID int, workspaceToken string) error {
	path := fmt.Sprintf("/sources/%d/refresh_tables", sourceID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodPost, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return fmt.Errorf("failed to make refresh tables request: %w", err)
	}

	if err := c.handleResponse(resp, nil); err != nil {
		return fmt.Errorf("failed to refresh tables: %w", err)
	}

	return nil
}

// GetSourceTableRefreshStatus checks the status of a table refresh
func (c *Client) GetSourceTableRefreshStatus(ctx context.Context, sourceID int) (*RefreshStatus, error) {
	return c.GetSourceTableRefreshStatusWithToken(ctx, sourceID, "")
}

// GetSourceTableRefreshStatusWithToken checks the status of a table refresh using a specific workspace token
func (c *Client) GetSourceTableRefreshStatusWithToken(ctx context.Context, sourceID int, workspaceToken string) (*RefreshStatus, error) {
	path := fmt.Sprintf("/sources/%d/refresh_tables_status", sourceID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get refresh status request: %w", err)
	}

	var result RefreshStatusResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get refresh status: %w", err)
	}

	return &result.Data, nil
}

// SourceTypeField represents a field configuration for a source type
type SourceTypeField struct {
	ID                    string      `json:"id"`
	Rules                 []string    `json:"rules"`
	Label                 string      `json:"label"`
	Type                  string      `json:"type"`
	Placeholder           interface{} `json:"placeholder"`
	IsPasswordTypeField   bool        `json:"is_password_type_field"`
	PossibleValues        []string    `json:"possible_values,omitempty"`
	Show                  interface{} `json:"show"`
	ConditionallyRequired interface{} `json:"conditionally_required"`
}

// SourceTypeConfiguration represents the configuration fields for a source type
type SourceTypeConfiguration struct {
	Fields []SourceTypeField `json:"fields"`
}

// SourceType represents a source type from the API
type SourceType struct {
	DocumentationSlug    string                  `json:"documentation_slug"`
	Label                string                  `json:"label"`
	ServiceName          string                  `json:"service_name"`
	SupportedSyncEngines []string                `json:"supported_sync_engines"`
	CreatableViaAPI      bool                    `json:"creatable_via_api"`
	EditableViaAPI       bool                    `json:"editable_via_api"`
	ConfigurationFields  SourceTypeConfiguration `json:"configuration_fields"`
}

// SourceTypesResponse represents the response from /source_types
type SourceTypesResponse struct {
	Status string       `json:"status"`
	Data   []SourceType `json:"data"`
}

// GetSourceTypes retrieves all available source types and their field requirements
func (c *Client) GetSourceTypes(ctx context.Context, workspaceToken string) ([]SourceType, error) {
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, "/source_types", nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get source types request: %w", err)
	}

	var result SourceTypesResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get source types: %w", err)
	}

	return result.Data, nil
}

// ValidateSourceCredentials validates source credentials against the source type requirements
func (c *Client) ValidateSourceCredentials(ctx context.Context, sourceType string, credentials map[string]interface{}, workspaceToken string) error {
	sourceTypes, err := c.GetSourceTypes(ctx, workspaceToken)
	if err != nil {
		return fmt.Errorf("failed to get source types for validation: %w", err)
	}

	// Find the source type
	var targetSourceType *SourceType
	for _, st := range sourceTypes {
		if st.ServiceName == sourceType {
			targetSourceType = &st
			break
		}
	}

	if targetSourceType == nil {
		return fmt.Errorf("unknown source type: %s", sourceType)
	}

	// Validate required fields and set defaults for optional ones
	for _, field := range targetSourceType.ConfigurationFields.Fields {
		// Check if field has required rule
		isRequired := false
		if len(field.Rules) > 0 {
			for _, rule := range field.Rules {
				if rule == "required" || rule == "required:notForEditing" {
					isRequired = true
					break
				}
			}
		}

		// Check for conditional requirements based on 'show' field
		if isRequired && field.Show != nil {
			// Handle conditional requirements like: "show": {"if": "ssh_tunnel_enabled"} or "show": {"if": {"ssh_tunnel_enabled": true}}
			if showMap, ok := field.Show.(map[string]interface{}); ok {
				if ifField, exists := showMap["if"]; exists {
					if ifFieldStr, ok := ifField.(string); ok {
						// Simple condition: "if": "ssh_tunnel_enabled"
						if conditionValue, conditionExists := credentials[ifFieldStr]; conditionExists {
							// Convert to boolean - only require if the condition is explicitly true
							if conditionStr, ok := conditionValue.(string); ok {
								if conditionStr != "true" && conditionStr != "1" {
									isRequired = false
								}
							} else if conditionBool, ok := conditionValue.(bool); ok {
								if !conditionBool {
									isRequired = false
								}
							} else {
								// If condition field is not set or falsy, don't require this field
								isRequired = false
							}
						} else {
							// If condition field doesn't exist, don't require this field
							isRequired = false
						}
					} else if ifFieldMap, ok := ifField.(map[string]interface{}); ok {
						// Complex condition: "if": {"ssh_tunnel_enabled": true}
						conditionMet := false
						for conditionField, expectedValue := range ifFieldMap {
							if conditionValue, conditionExists := credentials[conditionField]; conditionExists {
								// Check if the condition value matches the expected value
								if expectedBool, ok := expectedValue.(bool); ok {
									if conditionStr, ok := conditionValue.(string); ok {
										if expectedBool && (conditionStr == "true" || conditionStr == "1") {
											conditionMet = true
										} else if !expectedBool && (conditionStr == "false" || conditionStr == "0" || conditionStr == "") {
											conditionMet = true
										}
									} else if conditionBool, ok := conditionValue.(bool); ok {
										if conditionBool == expectedBool {
											conditionMet = true
										}
									}
								}
							}
							break // For now, just handle the first condition
						}
						if !conditionMet {
							isRequired = false
						}
					}
				} else if unlessField, exists := showMap["unless"]; exists {
					// Handle "unless" conditions: field is required UNLESS condition is true
					// Example: "show": {"unless": {"use_keypair": {"eq": true}}}
					conditionMet := false

					if unlessFieldMap, ok := unlessField.(map[string]interface{}); ok {
						for conditionField, expectedValue := range unlessFieldMap {
							if conditionValue, conditionExists := credentials[conditionField]; conditionExists {
								// Check if condition matches expected value
								if expectedMap, ok := expectedValue.(map[string]interface{}); ok {
									if eqValue, hasEq := expectedMap["eq"]; hasEq {
										// Handle {"eq": value} comparisons
										if expectedBool, ok := eqValue.(bool); ok {
											// Expected value is boolean
											if condBool, ok := conditionValue.(bool); ok {
												if condBool == expectedBool {
													conditionMet = true
												}
											} else if condStr, ok := conditionValue.(string); ok {
												// Convert string to boolean for comparison
												condBool := condStr == "true" || condStr == "1"
												if condBool == expectedBool {
													conditionMet = true
												}
											}
										} else if expectedStr, ok := eqValue.(string); ok {
											// Expected value is string
											if condStr, ok := conditionValue.(string); ok {
												if condStr == expectedStr {
													conditionMet = true
												}
											}
										}
									}
								} else if expectedBool, ok := expectedValue.(bool); ok {
									// Simple boolean comparison without "eq" wrapper
									if condBool, ok := conditionValue.(bool); ok {
										if condBool == expectedBool {
											conditionMet = true
										}
									} else if condStr, ok := conditionValue.(string); ok {
										condBool := condStr == "true" || condStr == "1"
										if condBool == expectedBool {
											conditionMet = true
										}
									}
								}
							}
							break // Handle first condition for now
						}
					}

					// If "unless" condition is met, field is NOT required
					if conditionMet {
						isRequired = false
					}
				}
			}
		}

		if isRequired {
			if _, exists := credentials[field.ID]; !exists {
				return fmt.Errorf("required field '%s' (%s) is missing", field.ID, field.Label)
			}
		} else {
			// Set empty string default for missing optional fields
			if _, exists := credentials[field.ID]; !exists {
				credentials[field.ID] = ""
			}
		}
	}

	return nil
}
