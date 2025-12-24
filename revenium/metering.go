package revenium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Package-level HTTP client with connection pooling for metering requests.
// This prevents creating a new client for each metering call, avoiding
// file descriptor exhaustion and TCP handshake overhead under high load.
var meteringHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true, // JSON is already small
	},
}

// MeteringClient handles communication with the Revenium metering API
type MeteringClient struct {
	config *Config
}

// NewMeteringClient creates a new metering client
func NewMeteringClient(config *Config) *MeteringClient {
	return &MeteringClient{
		config: config,
	}
}

// SendVideoMetering sends video generation metering data to Revenium
func (m *MeteringClient) SendVideoMetering(ctx context.Context, result *VideoGenerationResult, metadata *UsageMetadata) error {
	payload := m.buildMeteringPayload(result, metadata)

	// Send with retry logic
	return m.sendWithRetry(ctx, payload)
}

// buildMeteringPayload constructs the metering payload for video generation
func (m *MeteringClient) buildMeteringPayload(result *VideoGenerationResult, metadata *UsageMetadata) map[string]interface{} {
	now := time.Now()
	requestTime := now.Add(-result.Duration)

	// Determine stop reason
	stopReason := "END"
	if result.Status == TaskStatusFailed {
		stopReason = "ERROR"
	} else if result.Status == TaskStatusCanceled {
		stopReason = "CANCELLED"
	}

	// Extract video duration from metadata if available (default to 5 seconds for gen3a_turbo)
	var videoDurationSeconds float64 = 5.0 // Runway default
	if result.Metadata != nil {
		if dur, ok := result.Metadata["duration"].(int); ok {
			videoDurationSeconds = float64(dur)
		} else if dur, ok := result.Metadata["duration"].(float64); ok {
			videoDurationSeconds = dur
		} else if dur, ok := result.Metadata["durationSeconds"].(float64); ok {
			videoDurationSeconds = dur
		}
	}

	// Build base payload with durationSeconds at TOP LEVEL for billing (per API contract)
	payload := map[string]interface{}{
		"operationType":    "VIDEO",
		"provider":         "runway",
		"modelSource":      "RUNWAY",
		"model":            result.Model,
		"transactionId":    result.ID,
		"requestTime":      requestTime.Format(time.RFC3339),
		"responseTime":     now.Format(time.RFC3339),
		"requestDuration":  result.Duration.Milliseconds(),
		"durationSeconds":  videoDurationSeconds, // CRITICAL: video duration for billing
		"stopReason":       stopReason,
		"costType":         "AI",
		"isStreamed":       false,
		"middlewareSource": "revenium-middleware-runway-go",
	}

	// Add error information if failed
	if result.Error != nil {
		payload["errorReason"] = *result.Error
		payload["stopReason"] = "ERROR"
	}
	if result.FailureCode != nil {
		payload["failureCode"] = *result.FailureCode
	}

	// Add metadata from result
	if result.Metadata != nil {
		for k, v := range result.Metadata {
			// Only add if not already in payload
			if _, exists := payload[k]; !exists {
				payload[k] = v
			}
		}
	}

	// Add usage metadata if provided
	if metadata != nil {
		if metadata.OrganizationID != "" {
			payload["organizationId"] = metadata.OrganizationID
		}
		if metadata.ProductID != "" {
			payload["productId"] = metadata.ProductID
		}
		if metadata.TaskType != "" {
			payload["taskType"] = metadata.TaskType
		}
		if metadata.Agent != "" {
			payload["agent"] = metadata.Agent
		}
		if metadata.SubscriptionID != "" {
			payload["subscriptionId"] = metadata.SubscriptionID
		}
		if metadata.TraceID != "" {
			payload["traceId"] = metadata.TraceID
		}
		// Distributed tracing fields
		if metadata.ParentTransactionID != "" {
			payload["parentTransactionId"] = metadata.ParentTransactionID
		}
		if metadata.TraceType != "" {
			payload["traceType"] = metadata.TraceType
		}
		if metadata.TraceName != "" {
			payload["traceName"] = metadata.TraceName
		}
		if metadata.Environment != "" {
			payload["environment"] = metadata.Environment
		}
		if metadata.Region != "" {
			payload["region"] = metadata.Region
		}
		if metadata.RetryNumber != nil {
			payload["retryNumber"] = *metadata.RetryNumber
		}
		if metadata.CredentialAlias != "" {
			payload["credentialAlias"] = metadata.CredentialAlias
		}
		if metadata.Subscriber != nil {
			payload["subscriber"] = metadata.Subscriber
		}
		if metadata.TaskID != "" {
			payload["taskId"] = metadata.TaskID
		}
		if metadata.ResponseQualityScore != nil {
			payload["responseQualityScore"] = *metadata.ResponseQualityScore
		}
		if metadata.Custom != nil {
			for k, v := range metadata.Custom {
				// Only add if not already in payload
				if _, exists := payload[k]; !exists {
					payload[k] = v
				}
			}
		}
	}

	return payload
}

// sendWithRetry sends metering data with exponential backoff retry
func (m *MeteringClient) sendWithRetry(ctx context.Context, payload map[string]interface{}) error {
	const maxRetries = 3
	const initialBackoff = 100 * time.Millisecond

	var lastErr error
	backoff := initialBackoff

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
		}

		err := m.sendMeteringRequest(ctx, payload)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Don't retry on validation errors
		if IsValidationError(err) {
			return err
		}
	}

	return NewMeteringError("metering failed after retries", lastErr)
}

// sendMeteringRequest sends a single metering request to Revenium API
func (m *MeteringClient) sendMeteringRequest(ctx context.Context, payload map[string]interface{}) error {
	if m.config.ReveniumAPIKey == "" {
		return NewConfigError("Revenium API key not configured", nil)
	}

	// Build request URL - note: video endpoint is /meter/v2/ai/video
	baseURL := m.config.ReveniumBaseURL
	if baseURL == "" {
		baseURL = "https://api.revenium.ai"
	}
	url := baseURL + "/meter/v2/ai/video"

	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return NewMeteringError("failed to marshal metering payload", err)
	}

	Debug("[METERING] Sending video metering to %s: %s", url, string(jsonData))

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return NewMeteringError("failed to create metering request", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-api-key", m.config.ReveniumAPIKey)
	req.Header.Set("User-Agent", "revenium-middleware-runway-go/1.0")

	// Send request using pooled client (avoids creating new client per instance)
	resp, err := meteringHTTPClient.Do(req)
	if err != nil {
		return NewNetworkError("metering request failed", err)
	}
	defer resp.Body.Close()

	// Read response body for error details
	body, _ := io.ReadAll(resp.Body)

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// Validation error - don't retry
			return NewValidationError(
				fmt.Sprintf("metering API returned %d: %s", resp.StatusCode, string(body)),
				nil,
			)
		}
		return NewMeteringError("metering API error", fmt.Errorf("status %d: %s", resp.StatusCode, string(body)))
	}

	Debug("[METERING] Successfully sent metering data")
	return nil
}

// Close closes the metering client
func (m *MeteringClient) Close() error {
	// Nothing to clean up for HTTP client
	return nil
}
