// E2E Test for Runway Go Middleware
// Run with: RUNWAY_API_KEY=xxx REVENIUM_METERING_API_KEY=xxx go test -v -run TestE2E
package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/revenium/revenium-middleware-runway-go/revenium"
)

func TestE2EImageToVideo(t *testing.T) {
	// Check required env vars
	runwayKey := os.Getenv("RUNWAY_API_KEY")
	revKey := os.Getenv("REVENIUM_METERING_API_KEY")
	if runwayKey == "" || revKey == "" {
		t.Skip("Skipping E2E test: RUNWAY_API_KEY and REVENIUM_METERING_API_KEY required")
	}

	// Set up environment for DEV
	os.Setenv("REVENIUM_METERING_BASE_URL", "https://api.dev.hcapp.io")
	os.Setenv("REVENIUM_DEBUG", "true")

	// Initialize middleware
	if err := revenium.Initialize(); err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	client, err := revenium.GetClient()
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	// Create context
	ctx := context.Background()

	// Create comprehensive UsageMetadata with ALL fields
	metadata := &revenium.UsageMetadata{
		// Business context
		OrganizationID: "org-e2e-runway-test",
		ProductID:      "product-runway-video-test",
		SubscriptionID: "sub-runway-e2e-789",
		TaskType:       "e2e-image-to-video",
		Agent:          "e2e-test-agent-runway",
		// Tracing fields (ALL 10)
		TraceID:             "trace-runway-" + time.Now().Format("20060102-150405"),
		ParentTransactionID: "parent-runway-txn-456",
		TraceType:           "E2E_RUNWAY_TEST",
		TraceName:           "Runway Image-to-Video E2E Test",
		Environment:         "dev",
		Region:              "us-east-1",
		CredentialAlias:     "runway-e2e-test",
		TaskID:              "runway-task-" + time.Now().Format("150405"),
		// Multimodal job IDs (NEW)
		VideoJobID: "runway-vjob-" + time.Now().Format("150405"),
		AudioJobID: "", // Not applicable for image-to-video
		// Retry tracking
		RetryNumber: intPtr(0),
		// Quality score (for testing)
		ResponseQualityScore: float64Ptr(0.95),
		// Subscriber info
		Subscriber: map[string]interface{}{
			"id":           "user-runway-e2e",
			"email":        "runway-e2e@test.revenium.io",
			"tier":         "professional",
			"creditsBefore": 100,
		},
		// Custom fields
		Custom: map[string]interface{}{
			"testRun":     true,
			"testVersion": "1.0.0",
			"testDate":    time.Now().Format("2006-01-02"),
		},
	}

	// Create image-to-video request
	// Using a reliable public image from Unsplash CDN
	imageToVideoReq := &revenium.ImageToVideoRequest{
		PromptImage: "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=640",
		PromptText:  "Gentle zoom and subtle movement",
		Model:       "gen3a_turbo",
		Duration:    5, // Shortest duration for cost efficiency
		Ratio:       "1280:768",
	}

	t.Log("Generating video with Runway (gen3a_turbo)...")
	t.Log("NOTE: Video generation may take 1-3 minutes...")
	startTime := time.Now()

	// Generate video
	result, err := client.ImageToVideo(ctx, imageToVideoReq, metadata)
	if err != nil {
		t.Fatalf("Failed to generate video: %v", err)
	}

	duration := time.Since(startTime)
	t.Logf("Video generated in %v", duration)
	t.Logf("Task ID: %s", result.ID)
	t.Logf("Status: %s", result.Status)
	t.Logf("Model: %s", result.Model)
	t.Logf("Processing Duration: %v", result.Duration)

	if len(result.OutputURLs) > 0 {
		t.Logf("Output URL: %s", result.OutputURLs[0][:80]+"...")
	}

	// Wait for metering to complete
	t.Log("Waiting for metering data to be sent...")
	time.Sleep(3 * time.Second)

	t.Log("SUCCESS: Video generated and metering data sent to Revenium DEV")
	t.Log("Check dashboard: https://app.dev.hcapp.io -> AI Transaction Log")
	t.Log("")
	t.Log("VERIFY the following fields appear in the dashboard:")
	t.Log("  - videoJobId: " + metadata.VideoJobID)
	t.Log("  - traceId: " + metadata.TraceID)
	t.Log("  - environment: " + metadata.Environment)
	t.Log("  - region: " + metadata.Region)
	t.Log("  - responseQualityScore: 0.95")
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func main() {
	fmt.Println("Run with: go test -v -run TestE2E")
}
