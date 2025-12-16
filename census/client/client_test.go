package client

import (
	"net/http"
	"testing"
	"time"
)

func TestParseRetryAfter_DelaySeconds(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected time.Duration
		wantErr  bool
	}{
		{
			name:     "valid delay seconds",
			header:   "120",
			expected: 120 * time.Second,
			wantErr:  false,
		},
		{
			name:     "zero delay",
			header:   "0",
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "large delay",
			header:   "3600",
			expected: 3600 * time.Second,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRetryAfter(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRetryAfter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("parseRetryAfter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseRetryAfter_HTTPDate(t *testing.T) {
	// Test RFC 1123 format
	futureTime := time.Now().Add(30 * time.Second)
	header := futureTime.Format(time.RFC1123)

	got, err := parseRetryAfter(header)
	if err != nil {
		t.Errorf("parseRetryAfter() unexpected error: %v", err)
	}

	// Allow 2 second tolerance for test execution time
	if got < 28*time.Second || got > 32*time.Second {
		t.Errorf("parseRetryAfter() = %v, want approximately 30s", got)
	}
}

func TestParseRetryAfter_PastDate(t *testing.T) {
	// Date in the past should return 0 (retry immediately)
	pastTime := time.Now().Add(-30 * time.Second)
	header := pastTime.Format(time.RFC1123)

	got, err := parseRetryAfter(header)
	if err != nil {
		t.Errorf("parseRetryAfter() unexpected error: %v", err)
	}
	if got != 0 {
		t.Errorf("parseRetryAfter() = %v, want 0 for past date", got)
	}
}

func TestParseRetryAfter_Invalid(t *testing.T) {
	tests := []string{
		"invalid",
		"abc123",
		"not a date",
		"",
	}

	for _, header := range tests {
		t.Run(header, func(t *testing.T) {
			_, err := parseRetryAfter(header)
			if err == nil {
				t.Errorf("parseRetryAfter(%q) expected error, got nil", header)
			}
		})
	}
}

func TestCalculateRetryDelay_ExponentialBackoff(t *testing.T) {
	deadline := time.Now().Add(10 * time.Minute)

	tests := []struct {
		name        string
		attempt     int
		minExpected time.Duration
		maxExpected time.Duration
	}{
		{name: "attempt_1", attempt: 1, minExpected: 800 * time.Millisecond, maxExpected: 1200 * time.Millisecond},   // 1s ± 20%
		{name: "attempt_2", attempt: 2, minExpected: 1600 * time.Millisecond, maxExpected: 2400 * time.Millisecond},  // 2s ± 20%
		{name: "attempt_3", attempt: 3, minExpected: 3200 * time.Millisecond, maxExpected: 4800 * time.Millisecond},  // 4s ± 20%
		{name: "attempt_4", attempt: 4, minExpected: 6400 * time.Millisecond, maxExpected: 9600 * time.Millisecond},  // 8s ± 20%
		{name: "attempt_5", attempt: 5, minExpected: 12800 * time.Millisecond, maxExpected: 19200 * time.Millisecond}, // 16s ± 20%
		{name: "attempt_10", attempt: 10, minExpected: 72 * time.Second, maxExpected: 108 * time.Second},              // Capped at 90s ± 20%
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock response without Retry-After header
			resp := &http.Response{Header: http.Header{}}

			got := calculateRetryDelay(resp, tt.attempt, deadline)

			if got < tt.minExpected || got > tt.maxExpected {
				t.Errorf("calculateRetryDelay(attempt=%d) = %v, want between %v and %v",
					tt.attempt, got, tt.minExpected, tt.maxExpected)
			}
		})
	}
}

func TestCalculateRetryDelay_RetryAfterHeader(t *testing.T) {
	deadline := time.Now().Add(10 * time.Minute)

	// Test that Retry-After header is respected
	resp := &http.Response{
		Header: http.Header{
			"Retry-After": []string{"5"},
		},
	}

	got := calculateRetryDelay(resp, 10, deadline) // Attempt 10 would normally be 90s

	// Should use Retry-After value (5s) instead of exponential backoff (90s)
	if got != 5*time.Second {
		t.Errorf("calculateRetryDelay() with Retry-After = %v, want 5s", got)
	}
}

func TestCalculateRetryDelay_DeadlineRespect(t *testing.T) {
	// Deadline in 5 seconds
	deadline := time.Now().Add(5 * time.Second)

	// Attempt that would normally wait 90+ seconds
	resp := &http.Response{Header: http.Header{}}
	got := calculateRetryDelay(resp, 10, deadline)

	// Should be capped at remaining time (~5 seconds)
	if got > 5*time.Second {
		t.Errorf("calculateRetryDelay() = %v, should not exceed remaining deadline time of ~5s", got)
	}
}

func TestCalculateRetryDelay_RetryAfterExceedsDeadline(t *testing.T) {
	// Deadline in 2 seconds
	deadline := time.Now().Add(2 * time.Second)

	// Retry-After says wait 10 seconds
	resp := &http.Response{
		Header: http.Header{
			"Retry-After": []string{"10"},
		},
	}

	got := calculateRetryDelay(resp, 1, deadline)

	// Should be capped at remaining time (~2 seconds)
	if got > 2*time.Second {
		t.Errorf("calculateRetryDelay() = %v, should not exceed remaining deadline time of ~2s", got)
	}
}
