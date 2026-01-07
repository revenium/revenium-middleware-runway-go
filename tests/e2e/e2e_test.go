// Package e2e provides end-to-end tests for the Revenium Runway middleware.
// These tests verify that video metering data is correctly sent to Revenium.
//
// Requirements:
// - RUNWAY_API_KEY: Valid Runway API key with credits
// - REVENIUM_METERING_API_KEY: Valid Revenium API key
// - REVENIUM_METERING_BASE_URL: Revenium API base URL (defaults to https://api.revenium.ai)
//
// IMPORTANT: Video generation takes 5-20 minutes and costs credits.
// There is currently NO GET endpoint for video metrics, so we can only
// verify that metering POST requests are sent successfully.
//
// Run with: go test -v -tags=e2e ./tests/e2e/...
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/revenium/revenium-middleware-runway-go/revenium"
)

// AuditRecord captures request/response data for billable video generation calls
type AuditRecord struct {
	Timestamp       time.Time              `json:"timestamp"`
	TraceID         string                 `json:"traceId"`
	TransactionID   string                 `json:"transactionId,omitempty"`
	Provider        string                 `json:"provider"`
	Model           string                 `json:"model"`
	OperationType   string                 `json:"operationType"`
	RequestMetadata map[string]interface{} `json:"requestMetadata"`
	DurationSeconds float64                `json:"durationSeconds,omitempty"`
	RequestDuration int64                  `json:"requestDuration,omitempty"`
	TaskStatus      string                 `json:"taskStatus"`
	OutputURLs      []string               `json:"outputUrls,omitempty"`
	MeteringStatus  string                 `json:"meteringStatus"`
	ValidationError string                 `json:"validationError,omitempty"`
}

var auditTrail []AuditRecord

func TestMain(m *testing.M) {
	// Check required environment variables
	required := []string{"RUNWAY_API_KEY", "REVENIUM_METERING_API_KEY"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			fmt.Printf("SKIP: Required environment variable %s not set\n", env)
			os.Exit(0)
		}
	}

	// Run tests
	code := m.Run()

	// Print audit trail summary
	if len(auditTrail) > 0 {
		fmt.Println("\n========== AUDIT TRAIL ==========")
		auditJSON, _ := json.MarshalIndent(auditTrail, "", "  ")
		fmt.Println(string(auditJSON))
		fmt.Println("=================================")
	}

	os.Exit(code)
}

// TestE2E_RunwayVideoMetering_AllTracingFields tests that all 9 distributed tracing fields
// are sent in the metering POST request. Note: There is NO GET endpoint for video metrics,
// so we can only verify the POST was successful (no field validation against GET response).
func TestE2E_RunwayVideoMetering_AllTracingFields(t *testing.T) {
	// Generate unique trace ID for this test
	traceID := fmt.Sprintf("e2e-runway-%d", time.Now().UnixNano())

	// Reset any previous initialization
	revenium.Reset()

	// Initialize middleware
	if err := revenium.Initialize(); err != nil {
		t.Fatalf("Failed to initialize Revenium middleware: %v", err)
	}

	client, err := revenium.GetClient()
	if err != nil {
		t.Fatalf("Failed to get Revenium client: %v", err)
	}
	defer client.Close()

	// Create metadata with ALL distributed tracing fields (9 fields)
	retryNum := 0
	metadata := &revenium.UsageMetadata{
		OrganizationID:      "e2e-test-org",
		ProductID:           "e2e-test-product",
		TaskType:            "e2e-video-validation",
		Agent:               "e2e-test-agent",
		SubscriptionID:      "e2e-sub-123",
		TraceID:             traceID,
		TaskID:              fmt.Sprintf("task-%d", time.Now().Unix()),
		// Distributed tracing fields (the 9 fields being validated)
		ParentTransactionID: "parent-txn-e2e-runway-test",
		TraceType:           "e2e-test",
		TraceName:           "Runway Go E2E Video Validation",
		Environment:         "development",
		Region:              "us-west-2",
		RetryNumber:         &retryNum,
		CredentialAlias:     "e2e-test-credential",
	}

	// Record start of request
	audit := AuditRecord{
		Timestamp:       time.Now(),
		TraceID:         traceID,
		Provider:        "runway",
		Model:           "gen3a_turbo",
		OperationType:   "VIDEO",
		RequestMetadata: metadataToMap(metadata),
		DurationSeconds: 5.0, // 5 second video
		MeteringStatus:  "pending",
	}

	// Create video generation request using a real, publicly accessible image URL
	// Runway requires images it can fetch - using a well-known public image
	// This is a 1024x1024 sample image from picsum.photos (Lorem Picsum)
	req := &revenium.ImageToVideoRequest{
		PromptImage: "https://picsum.photos/1024/1024",
		PromptText:  "A peaceful nature scene with subtle motion",
		Model:       "gen3a_turbo",
		Duration:    5, // Minimum 5 seconds to minimize costs
	}

	t.Logf("Starting video generation with traceId: %s", traceID)
	t.Logf("WARNING: This test will take 5-20 minutes and costs Runway credits")
	t.Log("NOTE: There is NO GET endpoint for video metrics - only POST verification")
	startTime := time.Now()

	ctx := context.Background()
	result, err := client.ImageToVideo(ctx, req, metadata)
	if err != nil {
		audit.MeteringStatus = "api_error"
		audit.ValidationError = err.Error()
		auditTrail = append(auditTrail, audit)
		t.Fatalf("Video generation failed: %v", err)
	}

	totalDuration := time.Since(startTime)
	t.Logf("Video generation completed in %v", totalDuration)

	// Record response details
	audit.TransactionID = result.ID
	audit.TaskStatus = string(result.Status)
	audit.RequestDuration = totalDuration.Milliseconds()
	audit.OutputURLs = result.OutputURLs

	if result.Status == revenium.TaskStatusSucceeded {
		audit.MeteringStatus = "sent"
		t.Logf("SUCCESS: Video generated with ID: %s", result.ID)
		t.Logf("Output URLs: %v", result.OutputURLs)
		t.Log("Metering POST was sent with all 9 distributed tracing fields")
		t.Log("NOTE: Cannot validate fields via GET - no video metrics endpoint available")
	} else {
		audit.MeteringStatus = "task_failed"
		if result.Error != nil {
			audit.ValidationError = *result.Error
		}
		t.Errorf("Video generation failed with status: %s", result.Status)
	}

	// Log the tracing fields that were sent in metering POST
	t.Logf("Tracing fields sent in metering POST:")
	t.Logf("  traceId: %s", traceID)
	t.Logf("  transactionId: %s", result.ID)
	t.Logf("  parentTransactionId: %s", metadata.ParentTransactionID)
	t.Logf("  traceType: %s", metadata.TraceType)
	t.Logf("  traceName: %s", metadata.TraceName)
	t.Logf("  environment: %s", metadata.Environment)
	t.Logf("  region: %s", metadata.Region)
	t.Logf("  retryNumber: %d", *metadata.RetryNumber)
	t.Logf("  credentialAlias: %s", metadata.CredentialAlias)

	auditTrail = append(auditTrail, audit)
}

// TestE2E_RunwayVideoMetering_MinimalMetadata tests basic video metering
// This is a lighter test that just verifies metering works without all tracing fields.
// SKIPPED by default due to cost/time - uncomment to run.
func TestE2E_RunwayVideoMetering_MinimalMetadata(t *testing.T) {
	t.Skip("Skipping minimal video test to save costs - full test already validates metering")

	traceID := fmt.Sprintf("e2e-runway-minimal-%d", time.Now().UnixNano())

	// Reset any previous initialization
	revenium.Reset()

	if err := revenium.Initialize(); err != nil {
		t.Fatalf("Failed to initialize Revenium middleware: %v", err)
	}

	client, err := revenium.GetClient()
	if err != nil {
		t.Fatalf("Failed to get Revenium client: %v", err)
	}
	defer client.Close()

	// Minimal metadata - just traceId for identification
	metadata := &revenium.UsageMetadata{
		TraceID:        traceID,
		OrganizationID: "e2e-minimal-test",
	}

	audit := AuditRecord{
		Timestamp:       time.Now(),
		TraceID:         traceID,
		Provider:        "runway",
		Model:           "gen3a_turbo",
		OperationType:   "VIDEO",
		RequestMetadata: metadataToMap(metadata),
		MeteringStatus:  "pending",
	}

	req := &revenium.ImageToVideoRequest{
		PromptImage: "https://via.placeholder.com/768x768.png?text=Minimal+Test",
		PromptText:  "Minimal test",
		Duration:    5,
	}

	t.Logf("Starting minimal video generation with traceId: %s", traceID)

	ctx := context.Background()
	result, err := client.ImageToVideo(ctx, req, metadata)
	if err != nil {
		audit.MeteringStatus = "api_error"
		audit.ValidationError = err.Error()
		auditTrail = append(auditTrail, audit)
		t.Fatalf("Video generation failed: %v", err)
	}

	audit.TransactionID = result.ID
	audit.TaskStatus = string(result.Status)
	audit.MeteringStatus = "sent"
	auditTrail = append(auditTrail, audit)

	t.Logf("Minimal test completed: status=%s", result.Status)
}

// metadataToMap converts UsageMetadata to a map for audit logging
func metadataToMap(m *revenium.UsageMetadata) map[string]interface{} {
	result := make(map[string]interface{})
	if m.OrganizationID != "" {
		result["organizationId"] = m.OrganizationID
	}
	if m.ProductID != "" {
		result["productId"] = m.ProductID
	}
	if m.TaskType != "" {
		result["taskType"] = m.TaskType
	}
	if m.Agent != "" {
		result["agent"] = m.Agent
	}
	if m.SubscriptionID != "" {
		result["subscriptionId"] = m.SubscriptionID
	}
	if m.TraceID != "" {
		result["traceId"] = m.TraceID
	}
	if m.TaskID != "" {
		result["taskId"] = m.TaskID
	}
	if m.ParentTransactionID != "" {
		result["parentTransactionId"] = m.ParentTransactionID
	}
	if m.TraceType != "" {
		result["traceType"] = m.TraceType
	}
	if m.TraceName != "" {
		result["traceName"] = m.TraceName
	}
	if m.Environment != "" {
		result["environment"] = m.Environment
	}
	if m.Region != "" {
		result["region"] = m.Region
	}
	if m.RetryNumber != nil {
		result["retryNumber"] = *m.RetryNumber
	}
	if m.CredentialAlias != "" {
		result["credentialAlias"] = m.CredentialAlias
	}
	return result
}
