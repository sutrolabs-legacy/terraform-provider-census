package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Workspace represents a Census workspace
type Workspace struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	OrganizationID     int       `json:"organization_id"`
	CreatedAt          time.Time `json:"created_at"`
	NotificationEmails []string  `json:"notification_emails"`
	APIKey             string    `json:"api_key,omitempty"`
}

// CreateWorkspaceRequest represents the request to create a workspace
type CreateWorkspaceRequest struct {
	Name                  string   `json:"name"`
	NotificationEmails    []string `json:"notification_emails,omitempty"`
	ReturnWorkspaceAPIKey bool     `json:"return_workspace_api_key,omitempty"`
}

// UpdateWorkspaceRequest represents the request to update a workspace
type UpdateWorkspaceRequest struct {
	Name               string   `json:"name"`
	NotificationEmails []string `json:"notification_emails,omitempty"`
}

// WorkspaceResponse represents a single workspace response
type WorkspaceResponse struct {
	Status string     `json:"status"`
	Data   *Workspace `json:"data"`
}

// WorkspaceListResponse represents a paginated workspace list response
type WorkspaceListResponse struct {
	Status     string         `json:"status"`
	Pagination PaginationInfo `json:"pagination"`
	Data       []Workspace    `json:"data"`
}

// CreateWorkspace creates a new workspace
func (c *Client) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest) (*Workspace, error) {
	resp, err := c.makeRequest(ctx, http.MethodPost, "/workspaces", req, TokenTypePersonal)
	if err != nil {
		return nil, fmt.Errorf("failed to make create workspace request: %w", err)
	}

	var result WorkspaceResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	return result.Data, nil
}

// GetWorkspace retrieves a workspace by ID
func (c *Client) GetWorkspace(ctx context.Context, workspaceID int) (*Workspace, error) {
	path := fmt.Sprintf("/workspaces/%d", workspaceID)
	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, TokenTypePersonal)
	if err != nil {
		return nil, fmt.Errorf("failed to make get workspace request: %w", err)
	}

	var result WorkspaceResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	return result.Data, nil
}

// UpdateWorkspace updates an existing workspace
func (c *Client) UpdateWorkspace(ctx context.Context, workspaceID int, req *UpdateWorkspaceRequest) (*Workspace, error) {
	path := fmt.Sprintf("/workspaces/%d", workspaceID)
	resp, err := c.makeRequest(ctx, http.MethodPatch, path, req, TokenTypePersonal)
	if err != nil {
		return nil, fmt.Errorf("failed to make update workspace request: %w", err)
	}

	var result WorkspaceResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to update workspace: %w", err)
	}

	return result.Data, nil
}

// DeleteWorkspace deletes a workspace by ID
func (c *Client) DeleteWorkspace(ctx context.Context, workspaceID int) error {
	path := fmt.Sprintf("/workspaces/%d", workspaceID)
	resp, err := c.makeRequest(ctx, http.MethodDelete, path, nil, TokenTypePersonal)
	if err != nil {
		return fmt.Errorf("failed to make delete workspace request: %w", err)
	}

	if err := c.handleResponse(resp, nil); err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	return nil
}

// ListWorkspaces retrieves a list of workspaces
func (c *Client) ListWorkspaces(ctx context.Context, opts *ListOptions) ([]Workspace, *PaginationInfo, error) {
	params := make(map[string]string)
	if opts != nil {
		params = opts.ToParams()
	}

	fullURL := c.buildURL("/workspaces", params)
	resp, err := c.makeRequest(ctx, http.MethodGet, fullURL, nil, TokenTypePersonal)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make list workspaces request: %w", err)
	}

	var result WorkspaceListResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to list workspaces: %w", err)
	}

	return result.Data, &result.Pagination, nil
}

// GetAuthenticatedWorkspace retrieves the workspace for the authenticated workspace token
func (c *Client) GetAuthenticatedWorkspace(ctx context.Context) (*Workspace, error) {
	return c.GetAuthenticatedWorkspaceWithToken(ctx, "")
}

// GetAuthenticatedWorkspaceWithToken retrieves the workspace for a specific workspace token
func (c *Client) GetAuthenticatedWorkspaceWithToken(ctx context.Context, workspaceToken string) (*Workspace, error) {
	resp, err := c.makeRequestWithToken(ctx, http.MethodGet, "/workspace", nil, TokenTypeWorkspace, workspaceToken)
	if err != nil {
		return nil, fmt.Errorf("failed to make get authenticated workspace request: %w", err)
	}

	var result WorkspaceResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to get authenticated workspace: %w", err)
	}

	return result.Data, nil
}

// WorkspaceAPIKeyResponse represents the response from getting a workspace API key
type WorkspaceAPIKeyResponse struct {
	APIKey string `json:"api_key"`
}

// GetWorkspaceAPIKey retrieves the API key for a specific workspace.
// Results are cached for the lifetime of the client to avoid redundant API calls
// during a single Terraform run.
// Requires organization-level permissions (personal access token).
func (c *Client) GetWorkspaceAPIKey(ctx context.Context, workspaceID int) (string, error) {
	// Check cache first (read lock)
	c.tokenMu.RLock()
	if token, ok := c.workspaceTokenCache[workspaceID]; ok {
		c.tokenMu.RUnlock()
		return token, nil
	}
	c.tokenMu.RUnlock()

	// Cache miss — fetch from API
	path := fmt.Sprintf("/workspaces/%d/api_key", workspaceID)
	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil, TokenTypePersonal)
	if err != nil {
		return "", fmt.Errorf("failed to make get workspace API key request: %w", err)
	}

	var result WorkspaceAPIKeyResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return "", fmt.Errorf("failed to get workspace API key: %w", err)
	}

	// Store in cache (write lock)
	c.tokenMu.Lock()
	c.workspaceTokenCache[workspaceID] = result.APIKey
	c.tokenMu.Unlock()

	return result.APIKey, nil
}
