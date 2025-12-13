package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Config holds the configuration for the Census API client
type Config struct {
	PersonalAccessToken  string
	WorkspaceAccessToken string
	BaseURL              string
	Region               string
	HTTPClient           *http.Client
}

// Client represents a Census API client
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new Census API client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 8 * time.Minute,
		}
	}

	return &Client{
		config:     config,
		httpClient: httpClient,
	}, nil
}

// APIError represents an error response from the Census API
type APIError struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message,omitempty"`
	Status     string `json:"status_text,omitempty"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Census API error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("Census API error (status %d)", e.StatusCode)
}

// PaginationInfo holds pagination information from API responses
type PaginationInfo struct {
	TotalRecords int  `json:"total_records"`
	PerPage      int  `json:"per_page"`
	PrevPage     *int `json:"prev_page"`
	Page         int  `json:"page"`
	NextPage     *int `json:"next_page"`
	LastPage     int  `json:"last_page"`
}

// PaginatedResponse represents a paginated response from the Census API
type PaginatedResponse struct {
	Status     string         `json:"status"`
	Pagination PaginationInfo `json:"pagination"`
	Data       interface{}    `json:"data"`
}

// makeRequest performs an HTTP request to the Census API
func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}, tokenType TokenType) (*http.Response, error) {
	return c.makeRequestWithToken(ctx, method, path, body, tokenType, "")
}

// makeRequestWithToken performs an HTTP request to the Census API with a specific token
func (c *Client) makeRequestWithToken(ctx context.Context, method, path string, body interface{}, tokenType TokenType, specificToken string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	fullURL := c.config.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "terraform-provider-census")

	// Set authentication based on token type and availability
	token := ""
	if specificToken != "" {
		// Use the specific token provided (e.g., dynamically retrieved workspace token)
		token = specificToken
	} else {
		// Fall back to configured tokens
		switch tokenType {
		case TokenTypePersonal:
			token = c.config.PersonalAccessToken
		case TokenTypeWorkspace:
			token = c.config.WorkspaceAccessToken
		default:
			return nil, fmt.Errorf("invalid token type: %v", tokenType)
		}
	}

	if token == "" {
		return nil, fmt.Errorf("required token not provided for token type: %v", tokenType)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	return c.httpClient.Do(req)
}

// TokenType represents the type of authentication token to use
type TokenType int

const (
	TokenTypePersonal TokenType = iota
	TokenTypeWorkspace
)

// handleResponse processes an HTTP response and handles errors
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		apiErr.StatusCode = resp.StatusCode

		// Try to decode JSON error response
		if json.Unmarshal(body, &apiErr) != nil {
			// If JSON decode fails, use raw body as message
			apiErr.Message = string(body)
		}

		return &apiErr
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to decode response JSON: %w", err)
		}
	}

	return nil
}

// buildURL constructs a URL with query parameters
func (c *Client) buildURL(path string, params map[string]string) string {
	fullURL := c.config.BaseURL + path

	if len(params) == 0 {
		return fullURL
	}

	u, _ := url.Parse(fullURL)
	q := u.Query()

	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

// ListOptions represents options for list operations
type ListOptions struct {
	Page    int
	PerPage int
	Order   string
}

// ToParams converts ListOptions to URL parameters
func (opts *ListOptions) ToParams() map[string]string {
	params := make(map[string]string)

	if opts.Page > 0 {
		params["page"] = strconv.Itoa(opts.Page)
	}
	if opts.PerPage > 0 {
		params["per_page"] = strconv.Itoa(opts.PerPage)
	}
	if opts.Order != "" {
		params["order"] = opts.Order
	}

	return params
}
