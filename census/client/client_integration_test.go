package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestRateLimitRetry_Success verifies that 429 responses are retried and eventually succeed
func TestRateLimitRetry_Success(t *testing.T) {
	attemptCount := int32(0)

	// Mock server that returns 429 twice, then 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attemptCount, 1)

		if count <= 2 {
			w.Header().Set("Retry-After", "1") // 1 second
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"message": "rate limited"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(&Config{
		BaseURL:             server.URL,
		PersonalAccessToken: "test-token",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Make request
	ctx := context.Background()
	resp, err := client.makeRequest(ctx, http.MethodGet, "/test", nil, TokenTypePersonal)

	if err != nil {
		t.Fatalf("Expected successful retry, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attemptCount) != 3 {
		t.Errorf("Expected 3 attempts (2 failures + 1 success), got %d", attemptCount)
	}
}

// TestRateLimitRetry_Timeout verifies that retries stop when context times out
func TestRateLimitRetry_Timeout(t *testing.T) {
	// Mock server that always returns 429
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "2") // 2 second delay
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"message": "rate limited"}`))
	}))
	defer server.Close()

	// Create client with short timeout
	client, err := NewClient(&Config{
		BaseURL:             server.URL,
		PersonalAccessToken: "test-token",
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Make request with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startTime := time.Now()
	_, err = client.makeRequest(ctx, http.MethodGet, "/test", nil, TokenTypePersonal)
	elapsed := time.Since(startTime)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	// Should timeout around 5 seconds, allow some tolerance
	if elapsed < 4*time.Second || elapsed > 7*time.Second {
		t.Errorf("Expected timeout around 5s, took %v", elapsed)
	}
}

// TestNon429Error_NoRetry verifies that non-429 errors are not retried
func TestNon429Error_NoRetry(t *testing.T) {
	attemptCount := int32(0)

	// Mock server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attemptCount, 1)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		BaseURL:             server.URL,
		PersonalAccessToken: "test-token",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.makeRequest(ctx, http.MethodGet, "/test", nil, TokenTypePersonal)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attemptCount) != 1 {
		t.Errorf("Expected exactly 1 attempt (no retry for 500), got %d", attemptCount)
	}
}

// TestRateLimitRetry_RetryAfterHeader verifies that Retry-After header is respected
func TestRateLimitRetry_RetryAfterHeader(t *testing.T) {
	attemptCount := int32(0)
	var lastAttemptTime time.Time

	// Mock server that returns 429 once with Retry-After, then 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attemptCount, 1)

		if count == 1 {
			lastAttemptTime = time.Now()
			w.Header().Set("Retry-After", "2") // 2 seconds
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"message": "rate limited"}`))
			return
		}

		// Check that at least 2 seconds elapsed since first attempt
		elapsed := time.Since(lastAttemptTime)
		if elapsed < 1800*time.Millisecond { // Allow some tolerance
			t.Errorf("Second attempt came too soon: %v (expected ~2s)", elapsed)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		BaseURL:             server.URL,
		PersonalAccessToken: "test-token",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.makeRequest(ctx, http.MethodGet, "/test", nil, TokenTypePersonal)

	if err != nil {
		t.Fatalf("Expected successful retry, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attemptCount) != 2 {
		t.Errorf("Expected 2 attempts, got %d", attemptCount)
	}
}

// TestRateLimitRetry_ContextCancellation verifies that context cancellation stops retries
func TestRateLimitRetry_ContextCancellation(t *testing.T) {
	attemptCount := int32(0)

	// Mock server that always returns 429
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attemptCount, 1)
		w.Header().Set("Retry-After", "10") // Long wait
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"message": "rate limited"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		BaseURL:             server.URL,
		PersonalAccessToken: "test-token",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	startTime := time.Now()
	_, err = client.makeRequest(ctx, http.MethodGet, "/test", nil, TokenTypePersonal)
	elapsed := time.Since(startTime)

	if err == nil {
		t.Fatal("Expected context cancellation error, got nil")
	}

	// Should have cancelled quickly (within ~1-2 seconds)
	if elapsed > 3*time.Second {
		t.Errorf("Context cancellation took too long: %v", elapsed)
	}

	// Should have made at least 1 attempt before cancellation
	if atomic.LoadInt32(&attemptCount) < 1 {
		t.Errorf("Expected at least 1 attempt, got %d", attemptCount)
	}
}

// TestRateLimitRetry_BodyReplay verifies that POST requests with bodies are properly replayed
func TestRateLimitRetry_BodyReplay(t *testing.T) {
	attemptCount := int32(0)
	var receivedBodies []string

	// Mock server that returns 429 once, then 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&attemptCount, 1)

		// Read and store the body
		body := make([]byte, 1024)
		n, _ := r.Body.Read(body)
		receivedBodies = append(receivedBodies, string(body[:n]))

		if count == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"message": "rate limited"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	client, err := NewClient(&Config{
		BaseURL:             server.URL,
		PersonalAccessToken: "test-token",
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	testBody := map[string]string{"key": "value"}
	resp, err := client.makeRequest(ctx, http.MethodPost, "/test", testBody, TokenTypePersonal)

	if err != nil {
		t.Fatalf("Expected successful retry, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if atomic.LoadInt32(&attemptCount) != 2 {
		t.Errorf("Expected 2 attempts, got %d", attemptCount)
	}

	// Verify that both attempts received the same body
	if len(receivedBodies) != 2 {
		t.Fatalf("Expected 2 bodies, got %d", len(receivedBodies))
	}

	if receivedBodies[0] != receivedBodies[1] {
		t.Errorf("Body mismatch between attempts:\nFirst:  %s\nSecond: %s", receivedBodies[0], receivedBodies[1])
	}

	expectedBody := `{"key":"value"}`
	if receivedBodies[0] != expectedBody {
		t.Errorf("Unexpected body content: got %s, want %s", receivedBodies[0], expectedBody)
	}
}
