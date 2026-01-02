package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Sync represents a Census data sync
type Sync struct {
	ID          int       `json:"id"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	Label       string    `json:"label"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Source and destination configuration
	SourceAttributes      map[string]interface{} `json:"source_attributes,omitempty"`
	DestinationAttributes map[string]interface{} `json:"destination_attributes,omitempty"`

	// Field mappings and configuration (API returns mappings, not field_mappings)
	Mappings      []MappingAttributes `json:"mappings,omitempty"`
	FieldMappings []FieldMapping      `json:"field_mappings,omitempty"` // Keep for backward compatibility
	SyncKey       []string            `json:"sync_key,omitempty"`
	Operation     string              `json:"operation,omitempty"` // mirror, upsert, append, etc.

	// Scheduling configuration from API response
	ScheduleFrequency string  `json:"schedule_frequency,omitempty"`
	ScheduleDay       *string `json:"schedule_day,omitempty"`
	ScheduleHour      *int    `json:"schedule_hour,omitempty"`
	ScheduleMinute    *int    `json:"schedule_minute,omitempty"`
	CronExpression    string  `json:"cron_expression,omitempty"`

	// For terraform config - keep existing schedule structure
	Schedule *SyncSchedule `json:"schedule,omitempty"`
	Paused   bool          `json:"paused,omitempty"`

	// Mode - live vs triggered with trigger configurations
	Mode *SyncMode `json:"mode,omitempty"`

	// Run information
	LastRunAt *time.Time `json:"last_run_at,omitempty"`
	NextRunAt *time.Time `json:"next_run_at,omitempty"`
	LastRunID *int       `json:"last_run_id,omitempty"`

	// Field configuration
	FieldBehavior      string `json:"field_behavior,omitempty"`      // sync_all_properties or specific_properties
	FieldNormalization string `json:"field_normalization,omitempty"` // snake_case, camel_case, etc.
	FieldOrder         string `json:"field_order,omitempty"`         // alphabetical_column_name or mapping_order

	// Sync behavior family
	SyncBehaviorFamily string `json:"sync_behavior_family,omitempty"` // activateEvents or mapRecords

	// Advanced configuration - destination-specific options
	AdvancedConfiguration map[string]interface{} `json:"advanced_configuration,omitempty"`

	// High water mark attribute - timestamp column for append syncs
	HighWaterMarkAttribute string `json:"high_water_mark_attribute,omitempty"`

	// Historical sync operation - how first sync handles existing records
	HistoricalSyncOperation string `json:"historical_sync_operation,omitempty"` // skip_current_records or backfill_all_records

	// Mirror strategy - how mirror syncs keep destination in sync (when operation=mirror)
	MirrorStrategy string `json:"mirror_strategy,omitempty"` // sync_updates_and_deletes, sync_updates_and_nulls, upload_and_swap

	// Alert configuration
	AlertAttributes []AlertAttribute `json:"alert_attributes,omitempty"`
}

// AlertAttribute represents an alert configuration for a sync
type AlertAttribute struct {
	ID                 int                    `json:"id,omitempty"`
	Type               string                 `json:"type"`
	SendFor            string                 `json:"send_for"`
	ShouldSendRecovery bool                   `json:"should_send_recovery"`
	Options            map[string]interface{} `json:"options"`
}

// FieldMapping represents a field mapping between source and destination (for Terraform config)
type FieldMapping struct {
	From                string      `json:"from"`
	To                  string      `json:"to"`
	Type                string      `json:"operation,omitempty"`             // direct, hash, constant, sync_metadata, segment_membership, liquid_template - JSON is still "operation" for API compatibility
	Constant            interface{} `json:"constant,omitempty"`              // For constant mappings
	SyncMetadataKey     string      `json:"sync_metadata_key,omitempty"`     // For sync_metadata mappings (e.g., "sync_run_id")
	SegmentIdentifyBy   string      `json:"segment_identify_by,omitempty"`   // For segment_membership mappings (e.g., "name")
	LiquidTemplate      string      `json:"liquid_template,omitempty"`       // For liquid_template mappings (template content)
	IsPrimaryIdentifier bool        `json:"is_primary_identifier,omitempty"` // Whether this field is the sync key
	LookupObject        string      `json:"lookup_object,omitempty"`         // Object to lookup for relationship mapping
	LookupField         string      `json:"lookup_field,omitempty"`          // Field to use for lookup
	PreserveValues      bool        `json:"preserve_values,omitempty"`       // Whether to preserve existing destination values
	GenerateField       bool        `json:"generate_field,omitempty"`        // Whether Census should generate this field
	SyncNullValues      *bool       `json:"sync_null_values,omitempty"`      // Whether to sync null values (pointer to distinguish from false vs not set)
	ArrayField          bool        `json:"array_field,omitempty"`           // Whether the destination field is an array type
	FieldType           string      `json:"field_type,omitempty"`            // The type of the destination field (for user-defined fields)
	FollowSourceType    bool        `json:"follow_source_type,omitempty"`    // Whether destination field type should follow source column type
}

// MappingAttributes represents Census API mapping format (OpenAPI compliant)
type MappingAttributes struct {
	From                MappingFrom `json:"from"`
	To                  string      `json:"to"`
	IsPrimaryIdentifier bool        `json:"is_primary_identifier"`
	LookupObject        string      `json:"lookup_object,omitempty"`
	LookupField         string      `json:"lookup_field,omitempty"`
	PreserveValues      bool        `json:"preserve_values,omitempty"`
	GenerateField       bool        `json:"generate_field,omitempty"`
	SyncNullValues      *bool       `json:"sync_null_values,omitempty"` // Pointer to distinguish from false vs not set (default is true)
	ArrayField          bool        `json:"array_field,omitempty"`
	FieldType           string      `json:"field_type,omitempty"`
	FollowSourceType    bool        `json:"follow_source_type,omitempty"`
}

// MappingFrom represents the source of a mapping
type MappingFrom struct {
	Type string      `json:"type"` // "column", "constant_value", "sync_metadata", "segment_membership", "liquid_template"
	Data interface{} `json:"data"` // Data format varies by type
}

// SyncSchedule represents sync scheduling configuration
type SyncSchedule struct {
	Frequency string `json:"frequency"`             // hourly, daily, weekly, etc.
	Minute    int    `json:"minute,omitempty"`      // minute of hour to run (0-59)
	Hour      int    `json:"hour,omitempty"`        // for daily/weekly
	DayOfWeek int    `json:"day_of_week,omitempty"` // for weekly (0=Sunday)
	Timezone  string `json:"timezone,omitempty"`
}

// SyncMode represents the mode configuration for a sync (live vs triggered)
type SyncMode struct {
	Type     string        `json:"type"`               // live or triggered
	Triggers *SyncTriggers `json:"triggers,omitempty"` // trigger configurations
}

// SyncTriggers represents all possible trigger types for a sync
type SyncTriggers struct {
	Schedule     *TriggerSchedule     `json:"schedule,omitempty"`
	DbtCloud     *DbtCloudTrigger     `json:"dbt_cloud,omitempty"`
	Fivetran     *FivetranTrigger     `json:"fivetran,omitempty"`
	SyncSequence *SyncSequenceTrigger `json:"sync_sequence,omitempty"`
}

// TriggerSchedule represents schedule-based trigger configuration
type TriggerSchedule struct {
	Frequency      string `json:"frequency"`                 // never, continuous, quarter_hourly, hourly, daily, weekly, expression
	Day            string `json:"day,omitempty"`             // Sunday-Saturday (for weekly)
	Hour           int    `json:"hour,omitempty"`            // 0-24
	Minute         int    `json:"minute,omitempty"`          // 0-59
	CronExpression string `json:"cron_expression,omitempty"` // for frequency="expression"
}

// DbtCloudTrigger represents dbt Cloud job trigger configuration
type DbtCloudTrigger struct {
	ProjectId string `json:"project_id"`
	JobId     string `json:"job_id"`
}

// FivetranTrigger represents Fivetran job trigger configuration
type FivetranTrigger struct {
	JobId   string `json:"job_id"`
	JobName string `json:"job_name"`
}

// SyncSequenceTrigger represents sync dependency trigger configuration
type SyncSequenceTrigger struct {
	SyncId int `json:"sync_id"`
}

// CreateSyncRequest represents the request to create a sync (OpenAPI compliant)
type CreateSyncRequest struct {
	// Required fields per BaseSyncAttributes
	Operation             string                 `json:"operation"`              // Required: How records are synced
	SourceAttributes      map[string]interface{} `json:"source_attributes"`      // Required: Source configuration
	DestinationAttributes map[string]interface{} `json:"destination_attributes"` // Required: Destination configuration
	Mappings              []MappingAttributes    `json:"mappings"`               // Required: Field mappings

	// Optional fields
	Label           string           `json:"label,omitempty"`
	AlertAttributes []AlertAttribute `json:"alert_attributes,omitempty"`

	// Schedule fields - Census Management API expects these as flat fields, not nested object
	ScheduleFrequency string `json:"schedule_frequency,omitempty"`
	ScheduleDay       *int   `json:"schedule_day,omitempty"`
	ScheduleHour      *int   `json:"schedule_hour,omitempty"`
	ScheduleMinute    *int   `json:"schedule_minute,omitempty"`
	CronExpression    string `json:"cron_expression,omitempty"`

	// Mode - live vs triggered with trigger configurations (replaces schedule fields)
	Mode *SyncMode `json:"mode,omitempty"`

	Paused bool `json:"paused"`

	// Field configuration
	FieldBehavior      string `json:"field_behavior,omitempty"`      // sync_all_properties or specific_properties
	FieldNormalization string `json:"field_normalization,omitempty"` // snake_case, camel_case, etc.
	FieldOrder         string `json:"field_order,omitempty"`         // alphabetical_column_name or mapping_order

	// Sync behavior family
	SyncBehaviorFamily string `json:"sync_behavior_family,omitempty"` // activateEvents or mapRecords

	// Advanced configuration - destination-specific options
	AdvancedConfiguration map[string]interface{} `json:"advanced_configuration,omitempty"`

	// High water mark attribute - timestamp column for append syncs
	HighWaterMarkAttribute string `json:"high_water_mark_attribute,omitempty"`

	// Historical sync operation - how first sync handles existing records
	HistoricalSyncOperation string `json:"historical_sync_operation,omitempty"` // skip_current_records or backfill_all_records

	// Mirror strategy - how mirror syncs keep destination in sync (when operation=mirror)
	MirrorStrategy string `json:"mirror_strategy,omitempty"` // sync_updates_and_deletes, sync_updates_and_nulls, upload_and_swap
}

// UpdateSyncRequest represents the request to update a sync
type UpdateSyncRequest struct {
	Label                 string                 `json:"label,omitempty"`
	SourceAttributes      map[string]interface{} `json:"source_attributes,omitempty"`
	DestinationAttributes map[string]interface{} `json:"destination_attributes,omitempty"`
	Mappings              []MappingAttributes    `json:"mappings,omitempty"` // Use OpenAPI-compliant format like CreateSyncRequest
	SyncKey               []string               `json:"sync_key,omitempty"`
	SyncMode              string                 `json:"sync_mode,omitempty"`

	// Schedule fields - Census Management API expects these as flat fields, not nested object
	ScheduleFrequency string `json:"schedule_frequency,omitempty"`
	ScheduleDay       *int   `json:"schedule_day,omitempty"`
	ScheduleHour      *int   `json:"schedule_hour,omitempty"`
	ScheduleMinute    *int   `json:"schedule_minute,omitempty"`
	CronExpression    string `json:"cron_expression,omitempty"`

	// Mode - live vs triggered with trigger configurations (replaces schedule fields)
	Mode *SyncMode `json:"mode,omitempty"`

	Paused bool `json:"paused"`

	// Field configuration
	FieldBehavior      string `json:"field_behavior,omitempty"`      // sync_all_properties or specific_properties
	FieldNormalization string `json:"field_normalization,omitempty"` // snake_case, camel_case, etc.
	FieldOrder         string `json:"field_order,omitempty"`         // alphabetical_column_name or mapping_order

	// Sync behavior family
	SyncBehaviorFamily string `json:"sync_behavior_family,omitempty"` // activateEvents or mapRecords

	// Advanced configuration - destination-specific options
	AdvancedConfiguration map[string]interface{} `json:"advanced_configuration,omitempty"`

	// High water mark attribute - timestamp column for append syncs
	HighWaterMarkAttribute string `json:"high_water_mark_attribute,omitempty"`

	// Historical sync operation - how first sync handles existing records
	HistoricalSyncOperation string `json:"historical_sync_operation,omitempty"` // skip_current_records or backfill_all_records

	// Mirror strategy - how mirror syncs keep destination in sync (when operation=mirror)
	MirrorStrategy string `json:"mirror_strategy,omitempty"` // sync_updates_and_deletes, sync_updates_and_nulls, upload_and_swap

	// Alert configuration
	// Note: No omitempty tag - we need to send empty array to delete all alerts
	AlertAttributes []AlertAttribute `json:"alert_attributes"`
}

// SyncResponse represents a single sync response
type SyncResponse struct {
	Status string `json:"status"`
	Data   *Sync  `json:"data"`
}

// CreateSyncResponse represents the response from creating a sync
type CreateSyncResponse struct {
	Status string `json:"status"`
	Data   struct {
		SyncID int `json:"sync_id"`
	} `json:"data"`
}

// SyncListResponse represents a paginated sync list response
type SyncListResponse struct {
	Status     string         `json:"status"`
	Pagination PaginationInfo `json:"pagination"`
	Data       []Sync         `json:"data"`
}

// SyncRun represents a sync execution
type SyncRun struct {
	ID               int        `json:"id"`
	SyncID           int        `json:"sync_id"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	StartedAt        *time.Time `json:"started_at,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	RecordsProcessed int        `json:"records_processed,omitempty"`
	RecordsSucceeded int        `json:"records_succeeded,omitempty"`
	RecordsFailed    int        `json:"records_failed,omitempty"`
	ErrorMessage     string     `json:"error_message,omitempty"`
}

// TriggerSyncRequest represents a request to trigger a sync
type TriggerSyncRequest struct {
	ForceFullSync bool `json:"force_full_sync,omitempty"`
}

// TriggerSyncResponse represents the response from triggering a sync
type TriggerSyncResponse struct {
	Status string `json:"status"`
	Data   struct {
		SyncRunID int `json:"sync_run_id"`
	} `json:"data"`
}

// SyncRunResponse represents a single sync run response
type SyncRunResponse struct {
	Status string   `json:"status"`
	Data   *SyncRun `json:"data"`
}

// CreateSync creates a new sync
func (c *Client) CreateSync(ctx context.Context, req *CreateSyncRequest) (*Sync, error) {
	return c.CreateSyncWithToken(ctx, req, "")
}

// CreateSyncWithToken creates a new sync using a specific workspace token
func (c *Client) CreateSyncWithToken(ctx context.Context, req *CreateSyncRequest, workspaceToken string) (*Sync, error) {
	// Log the request being sent
	fmt.Printf("[DEBUG] Create sync request: %+v\n", req)

	resp, err := c.makeRequestWithToken(ctx, http.MethodPost, "/syncs", req, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make create sync request: %w", err)
	}

	// Read the raw response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("[DEBUG] Create sync raw response: %s\n", string(bodyBytes))

	// Write request and response to debug file
	debugFile := "/tmp/census_sync_debug.log"
	debugContent := fmt.Sprintf("=== CREATE SYNC DEBUG ===\n\nREQUEST:\n%+v\n\nRAW RESPONSE:\n%s\n\n", req, string(bodyBytes))
	os.WriteFile(debugFile, []byte(debugContent), 0644)
	fmt.Printf("[DEBUG] Debug info written to %s\n", debugFile)

	// Reset the response body so handleResponse can read it
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var result CreateSyncResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to create sync: %w", err)
	}

	fmt.Printf("[DEBUG] Parsed create sync response: %+v\n", result)

	// Create a minimal Sync object from the response
	sync := &Sync{
		ID: result.Data.SyncID,
	}

	fmt.Printf("[DEBUG] Created sync object with ID: %d\n", sync.ID)
	return sync, nil
}

// GetSync retrieves a sync by ID
func (c *Client) GetSync(ctx context.Context, syncID int) (*Sync, error) {
	return c.GetSyncWithToken(ctx, syncID, "")
}

// GetSyncWithToken retrieves a sync by ID using a specific workspace token
func (c *Client) GetSyncWithToken(ctx context.Context, syncID int, workspaceToken string) (*Sync, error) {
	path := fmt.Sprintf("/syncs/%d", syncID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get sync request: %w", err)
	}

	var result SyncResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get sync: %w", err)
	}

	return result.Data, nil
}

// UpdateSync updates an existing sync
func (c *Client) UpdateSync(ctx context.Context, syncID int, req *UpdateSyncRequest) (*Sync, error) {
	return c.UpdateSyncWithToken(ctx, syncID, req, "")
}

// UpdateSyncWithToken updates an existing sync using a specific workspace token
func (c *Client) UpdateSyncWithToken(ctx context.Context, syncID int, req *UpdateSyncRequest, workspaceToken string) (*Sync, error) {
	// Log the request being sent
	fmt.Printf("[DEBUG] Update sync %d request: %+v\n", syncID, req)

	// Also log the JSON that will be sent
	if reqJSON, err := json.MarshalIndent(req, "", "  "); err == nil {
		fmt.Printf("[DEBUG] Update sync %d request JSON:\n%s\n", syncID, string(reqJSON))
	}

	path := fmt.Sprintf("/syncs/%d", syncID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodPatch, path, req, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make update sync request: %w", err)
	}

	// Read the raw response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("[DEBUG] Update sync %d raw response: %s\n", syncID, string(bodyBytes))

	// Write request and response to debug file
	debugFile := fmt.Sprintf("/tmp/census_sync_update_%d_debug.log", syncID)
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	debugContent := fmt.Sprintf("=== UPDATE SYNC %d DEBUG ===\n\nREQUEST:\n%+v\n\nREQUEST JSON:\n%s\n\nRAW RESPONSE:\n%s\n\n", syncID, req, string(reqJSON), string(bodyBytes))
	os.WriteFile(debugFile, []byte(debugContent), 0644)
	fmt.Printf("[DEBUG] Debug info written to %s\n", debugFile)

	// Reset the response body so handleResponse can read it
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var result SyncResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to update sync: %w", err)
	}

	return result.Data, nil
}

// DeleteSync deletes a sync by ID
func (c *Client) DeleteSync(ctx context.Context, syncID int) error {
	return c.DeleteSyncWithToken(ctx, syncID, "")
}

// DeleteSyncWithToken deletes a sync by ID using a specific workspace token
func (c *Client) DeleteSyncWithToken(ctx context.Context, syncID int, workspaceToken string) error {
	path := fmt.Sprintf("/syncs/%d", syncID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodDelete, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return fmt.Errorf("failed to make delete sync request: %w", err)
	}

	if err := c.handleResponse(resp, nil); err != nil {
		return fmt.Errorf("failed to delete sync: %w", err)
	}

	return nil
}

// ListSyncs retrieves a list of syncs
func (c *Client) ListSyncs(ctx context.Context, opts *ListOptions) ([]Sync, *PaginationInfo, error) {
	return c.ListSyncsWithToken(ctx, opts, "")
}

// ListSyncsWithToken retrieves a list of syncs using a specific workspace token
func (c *Client) ListSyncsWithToken(ctx context.Context, opts *ListOptions, workspaceToken string) ([]Sync, *PaginationInfo, error) {
	params := make(map[string]string)
	if opts != nil {
		params = opts.ToParams()
	}

	fullURL := c.buildURL("/syncs", params)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, fullURL, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make list syncs request: %w", err)
	}

	var result SyncListResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to list syncs: %w", err)
	}

	return result.Data, &result.Pagination, nil
}

// TriggerSync triggers a sync execution
func (c *Client) TriggerSync(ctx context.Context, syncID int, req *TriggerSyncRequest) (int, error) {
	return c.TriggerSyncWithToken(ctx, syncID, req, "")
}

// TriggerSyncWithToken triggers a sync execution using a specific workspace token
func (c *Client) TriggerSyncWithToken(ctx context.Context, syncID int, req *TriggerSyncRequest, workspaceToken string) (int, error) {
	path := fmt.Sprintf("/syncs/%d/trigger", syncID)

	var body interface{}
	if req != nil {
		body = req
	}

	resp, err := c.makeRequestWithToken(ctx, http.MethodPost, path, body, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return 0, fmt.Errorf("failed to make trigger sync request: %w", err)
	}

	var result TriggerSyncResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to trigger sync: %w", err)
	}

	return result.Data.SyncRunID, nil
}

// GetSyncRun retrieves details about a sync run
func (c *Client) GetSyncRun(ctx context.Context, syncRunID int) (*SyncRun, error) {
	return c.GetSyncRunWithToken(ctx, syncRunID, "")
}

// GetSyncRunWithToken retrieves details about a sync run using a specific workspace token
func (c *Client) GetSyncRunWithToken(ctx context.Context, syncRunID int, workspaceToken string) (*SyncRun, error) {
	path := fmt.Sprintf("/sync_runs/%d", syncRunID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get sync run request: %w", err)
	}

	var result SyncRunResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get sync run: %w", err)
	}

	return result.Data, nil
}

// CancelSyncRun cancels a running sync
func (c *Client) CancelSyncRun(ctx context.Context, syncRunID int) error {
	return c.CancelSyncRunWithToken(ctx, syncRunID, "")
}

// CancelSyncRunWithToken cancels a running sync using a specific workspace token
func (c *Client) CancelSyncRunWithToken(ctx context.Context, syncRunID int, workspaceToken string) error {
	path := fmt.Sprintf("/sync_runs/%d/cancel", syncRunID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodPost, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return fmt.Errorf("failed to make cancel sync run request: %w", err)
	}

	if err := c.handleResponse(resp, nil); err != nil {
		return fmt.Errorf("failed to cancel sync run: %w", err)
	}

	return nil
}
