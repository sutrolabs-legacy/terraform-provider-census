package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Dataset represents a Census dataset (SQL-based data transformation)
type Dataset struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	Type               string          `json:"type"` // "sql"
	Description        *string         `json:"description,omitempty"`
	Query              string          `json:"query,omitempty"`
	SourceID           int             `json:"source_id,omitempty"`
	ResourceIdentifier string          `json:"resource_identifier,omitempty"`
	CachedRecordCount  *int            `json:"cached_record_count,omitempty"`
	Columns            []DatasetColumn `json:"columns,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	WorkspaceID        string          `json:"-"` // Not returned by API, set by client
}

// DatasetColumn represents a column in a dataset
type DatasetColumn struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
}

// CreateDatasetRequest represents the request to create a new dataset
type CreateDatasetRequest struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // "sql"
	Description *string `json:"description,omitempty"`
	Query       string  `json:"query"`
	SourceID    int     `json:"source_id"`
}

// UpdateDatasetRequest represents the request to update a dataset
type UpdateDatasetRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Query       *string `json:"query,omitempty"`
}

// DatasetResponse represents a single dataset response from the API
type DatasetResponse struct {
	Status string   `json:"status"`
	Data   *Dataset `json:"data"`
}

// DatasetListResponse represents a paginated dataset list response
type DatasetListResponse struct {
	Status     string         `json:"status"`
	Pagination PaginationInfo `json:"pagination"`
	Data       []Dataset      `json:"data"`
}

// CreateDataset creates a new dataset
func (c *Client) CreateDataset(ctx context.Context, req *CreateDatasetRequest) (*Dataset, error) {
	return c.CreateDatasetWithToken(ctx, req, "")
}

// CreateDatasetWithToken creates a new dataset using a workspace token
func (c *Client) CreateDatasetWithToken(ctx context.Context, req *CreateDatasetRequest, workspaceToken string) (*Dataset, error) {
	resp, err := c.makeRequestWithToken(ctx, http.MethodPost, "/datasets", req, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make create dataset request: %w", err)
	}

	var result DatasetResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to create dataset: %w", err)
	}

	return result.Data, nil
}

// GetDataset retrieves a dataset by ID
func (c *Client) GetDataset(ctx context.Context, id int) (*Dataset, error) {
	return c.GetDatasetWithToken(ctx, id, "")
}

// GetDatasetWithToken retrieves a dataset by ID using a workspace token
func (c *Client) GetDatasetWithToken(ctx context.Context, id int, workspaceToken string) (*Dataset, error) {
	path := fmt.Sprintf("/datasets/%d", id)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get dataset request: %w", err)
	}

	var result DatasetResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get dataset: %w", err)
	}

	return result.Data, nil
}

// UpdateDataset updates a dataset
func (c *Client) UpdateDataset(ctx context.Context, id int, req *UpdateDatasetRequest) (*Dataset, error) {
	return c.UpdateDatasetWithToken(ctx, id, req, "")
}

// UpdateDatasetWithToken updates a dataset using a workspace token
func (c *Client) UpdateDatasetWithToken(ctx context.Context, id int, req *UpdateDatasetRequest, workspaceToken string) (*Dataset, error) {
	path := fmt.Sprintf("/datasets/%d", id)
	resp, err := c.makeRequestWithToken(ctx, http.MethodPatch, path, req, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make update dataset request: %w", err)
	}

	var result DatasetResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to update dataset: %w", err)
	}

	return result.Data, nil
}

// DeleteDataset deletes a dataset
func (c *Client) DeleteDataset(ctx context.Context, id int) error {
	return c.DeleteDatasetWithToken(ctx, id, "")
}

// DeleteDatasetWithToken deletes a dataset using a workspace token
func (c *Client) DeleteDatasetWithToken(ctx context.Context, id int, workspaceToken string) error {
	path := fmt.Sprintf("/datasets/%d", id)
	resp, err := c.makeRequestWithToken(ctx, http.MethodDelete, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return fmt.Errorf("failed to make delete dataset request: %w", err)
	}

	// For DELETE, we expect either 200 or 204
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var apiErr APIError
		if err := c.handleResponse(resp, &apiErr); err != nil {
			return fmt.Errorf("failed to delete dataset: %w", err)
		}
		return &apiErr
	}

	return nil
}

// ListDatasets lists all datasets in a workspace
func (c *Client) ListDatasets(ctx context.Context) ([]Dataset, error) {
	return c.ListDatasetsWithToken(ctx, "")
}

// ListDatasetsWithToken lists all datasets in a workspace using a workspace token
func (c *Client) ListDatasetsWithToken(ctx context.Context, workspaceToken string) ([]Dataset, error) {
	// Filter by SQL type as per OpenAPI spec
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, "/datasets?type=sql", nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make list datasets request: %w", err)
	}

	var result DatasetListResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to list datasets: %w", err)
	}

	return result.Data, nil
}


type RefreshColumnsResponse struct {
	RefreshKey int `json:"refresh_key"`
}

type RefreshColumnsStatusResponse struct {
	Status  string  `json:"status"`           
	Message *string `json:"message,omitempty"` 
}

func (c *Client) RefreshDatasetColumns(ctx context.Context, datasetID int) (int, error) {
	return c.RefreshDatasetColumnsWithToken(ctx, datasetID, "")
}

func (c *Client) RefreshDatasetColumnsWithToken(ctx context.Context, datasetID int, workspaceToken string) (int, error) {
	path := fmt.Sprintf("/datasets/%d/refresh_columns", datasetID)
	resp, err := c.makeRequestWithToken(ctx, http.MethodPost, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return 0, fmt.Errorf("failed to make refresh columns request: %w", err)
	}

	var result RefreshColumnsResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to refresh columns: %w", err)
	}

	return result.RefreshKey, nil
}

func (c *Client) GetDatasetRefreshStatus(ctx context.Context, datasetID int, refreshKey int) (*RefreshColumnsStatusResponse, error) {
	return c.GetDatasetRefreshStatusWithToken(ctx, datasetID, refreshKey, "")
}


func (c *Client) GetDatasetRefreshStatusWithToken(ctx context.Context, datasetID int, refreshKey int, workspaceToken string) (*RefreshColumnsStatusResponse, error) {
	path := fmt.Sprintf("/datasets/%d/refresh_columns_status?refresh_key=%d", datasetID, refreshKey)
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, path, nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get refresh status request: %w", err)
	}

	var result RefreshColumnsStatusResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get refresh status: %w", err)
	}

	return &result, nil
}

