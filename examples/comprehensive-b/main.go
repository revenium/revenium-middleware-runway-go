// Package main demonstrates Scenario B: Different values for hard-coding detection.
//
// This example uses COMPLETELY DIFFERENT values from comprehensive/main.go
// to verify no values are accidentally hard-coded in the middleware.
//
// Compare the metering payloads from both scenarios - they should differ
// in every user-settable field.
//
// Run with: RUNWAY_API_KEY=... REVENIUM_METERING_API_KEY=hak_... go run main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/revenium/revenium-middleware-runway-go/revenium"
)

// GetScenarioBMetadata - INTENTIONALLY DIFFERENT from Scenario A
// Every user-settable field has a DIFFERENT value to detect hard-coding
func GetScenarioBMetadata() *revenium.UsageMetadata {
	retryNumber := 3 // Different: was 0
	qualityScore := 0.68 // Different: was 0.95

	return &revenium.UsageMetadata{
		// Organization & Product - ALL DIFFERENT FROM SCENARIO A
		OrganizationID: "org-indie-studio-dev",               // Different: was "org-videotech-studios-prod"
		ProductID:      "prod-experimental-video-prototype",  // Different

		// Task Identification - ALL DIFFERENT
		TaskType: "prototype-short-clip-test",                // Different
		TaskID:   "task-prototype-iteration-42",              // Different

		// Agent/Worker - DIFFERENT
		Agent: "local-dev-macbook-m3-01",                     // Different

		// Subscription - DIFFERENT
		SubscriptionID: "sub-indie-monthly-jan2026",          // Different

		// Distributed Tracing - ALL DIFFERENT
		TraceID:             "trace-local-9999-aaaa-bbbb-cccc", // Different
		ParentTransactionID: "parent-local-debug-main",         // Different
		TraceType:           "local-debug",                     // Different: was "distributed"
		TraceName:           "prototype-video-experiment",      // Different

		// Environment - ALL DIFFERENT
		Environment: "development",                           // Different: was "production"
		Region:      "us-west-1",                             // Different: was "eu-central-1"

		// Credentials - DIFFERENT
		CredentialAlias: "runway-dev-key-backup",             // Different

		// Retry - DIFFERENT
		RetryNumber: &retryNumber,

		// Quality - DIFFERENT
		ResponseQualityScore: &qualityScore,

		// Multimodal Jobs - ALL DIFFERENT
		VideoJobID: "vjob-prototype-test-999",                // Different
		AudioJobID: "ajob-prototype-sfx-999",                 // Different

		// Subscriber - COMPLETELY DIFFERENT
		Subscriber: map[string]interface{}{
			"id":               "user-indie-dev-alice",       // Different
			"email":            "alice@indie-studio.dev",     // Different
			"name":             "Alice Indie Developer",      // Different
			"role":             "Solo Developer",             // Different
			"department":       "Engineering (Solo)",         // Different
			"costCenter":       "CC-INDIE-PERSONAL",          // Different
			"projectCode":      "PRJ-WEEKEND-EXPERIMENT",     // Different
			"quotaGroup":       "indie-limited",              // Different
			"billingTier":      "basic",                      // Different: was "premium"
		},

		// Custom Fields - ALL DIFFERENT
		Custom: map[string]interface{}{
			// Different campaign
			"projectName":     "Weekend Experiment Project",
			"experimentType":  "performance-benchmark",
			"debugMode":       true,                          // Different: enterprise had false
			"testIteration":   42,

			// Different content
			"contentType":     "test-clip",
			"targetPlatform":  "internal-review",

			// Different cost tracking
			"budgetCategory":  "r-and-d-experiments",
			"freeCreditsUsed": true,
		},
	}
}

func main() {
	fmt.Println("============================================================")
	fmt.Println("  Runway Middleware - SCENARIO B: Different Values Test")
	fmt.Println("  Purpose: Detect hard-coded values by using different data")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Compare these values against Scenario A (comprehensive/main.go)")
	fmt.Println("Every field should be DIFFERENT in the metering payload.")
	fmt.Println()

	// Enable verbose logging
	os.Setenv("REVENIUM_LOG_LEVEL", "DEBUG")
	os.Setenv("REVENIUM_VERBOSE_STARTUP", "true")

	// Initialize with prompt capture enabled
	// Enable prompt capture for analytics (opt-in, default: false)
	// When enabled, generation prompts are captured and sent with metering data
	// Fields added to metering payload:
	//   - inputMessages: JSON array with role/content format
	//   - outputResponse: Generated video URLs
	//   - promptsTruncated: true if prompt exceeded 50K chars
	if err := revenium.Initialize(
		revenium.WithCapturePrompts(true), // Enable prompt capture for analytics
	); err != nil {
		log.Fatalf("Failed to initialize middleware: %v", err)
	}

	client, err := revenium.GetClient()
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Minute)
	defer cancel()

	// Get Scenario B metadata (DIFFERENT from Scenario A)
	metadata := GetScenarioBMetadata()

	// Print the metadata being sent
	fmt.Println("=== SCENARIO B METADATA (should differ from A) ===")
	fmt.Println()
	printMetadata(metadata)
	fmt.Println()

	// Create request - shorter duration for faster testing
	imageToVideoReq := &revenium.ImageToVideoRequest{
		PromptImage: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b6/Image_created_with_a_mobile_phone.png/1200px-Image_created_with_a_mobile_phone.png",
		PromptText:  "A static test pattern slowly shifting colors, minimal motion, debug overlay visible",
		Model:       "gen3a_turbo",
		Duration:    5, // Different: was 10 (shorter for faster test)
		Ratio:       "768:1280", // Different: was "1280:768" (portrait vs landscape)
	}

	fmt.Println("Request Configuration (Scenario B):")
	fmt.Printf("  Model:    %s\n", imageToVideoReq.Model)
	fmt.Printf("  Duration: %d seconds (shorter than A)\n", imageToVideoReq.Duration)
	fmt.Printf("  Ratio:    %s (portrait, different from A)\n", imageToVideoReq.Ratio)
	fmt.Println()

	fmt.Println("Generating video... (this may take several minutes)")
	fmt.Println()

	result, err := client.ImageToVideo(ctx, imageToVideoReq, metadata)
	if err != nil {
		log.Fatalf("Failed to generate video: %v", err)
	}

	// Display results
	fmt.Println()
	fmt.Println("=== SCENARIO B RESULTS ===")
	fmt.Printf("  Task ID:     %s\n", result.ID)
	fmt.Printf("  Status:      %s\n", result.Status)
	fmt.Printf("  Model:       %s\n", result.Model)
	fmt.Printf("  Duration:    %v\n", result.Duration)
	fmt.Printf("  Output URLs: %v\n", result.OutputURLs)
	fmt.Println()

	// Generate expected payload for comparison
	fmt.Println("=== EXPECTED PAYLOAD DIFFERENCES FROM SCENARIO A ===")
	fmt.Println()
	fmt.Println("Field                    | Scenario A                      | Scenario B (this)")
	fmt.Println("-------------------------|--------------------------------|----------------------------------")
	fmt.Println("organizationId           | org-videotech-studios-prod     | org-indie-studio-dev")
	fmt.Println("productId                | prod-ai-video-gen-platform-v2  | prod-experimental-video-prototype")
	fmt.Println("environment              | production                      | development")
	fmt.Println("region                   | eu-central-1                    | us-west-1")
	fmt.Println("traceType                | distributed                     | local-debug")
	fmt.Println("retryNumber              | 0                               | 3")
	fmt.Println("responseQualityScore     | 0.95                            | 0.68")
	fmt.Println("requestedDurationSeconds | 10                              | 5")
	fmt.Println()

	fmt.Println("============================================================")
	fmt.Println("  SCENARIO B Complete")
	fmt.Println()
	fmt.Println("  VALIDATION STEPS:")
	fmt.Println("  1. Run comprehensive/main.go (Scenario A)")
	fmt.Println("  2. Run comprehensive-b/main.go (this, Scenario B)")
	fmt.Println("  3. Compare DEBUG logs for metering payloads")
	fmt.Println("  4. ALL user-settable fields should be DIFFERENT")
	fmt.Println("============================================================")
}

func printMetadata(m *revenium.UsageMetadata) {
	// Pretty print as JSON for comparison
	jsonBytes, _ := json.MarshalIndent(map[string]interface{}{
		"organizationId":       m.OrganizationID,
		"productId":            m.ProductID,
		"taskType":             m.TaskType,
		"taskId":               m.TaskID,
		"agent":                m.Agent,
		"subscriptionId":       m.SubscriptionID,
		"traceId":              m.TraceID,
		"parentTransactionId":  m.ParentTransactionID,
		"traceType":            m.TraceType,
		"traceName":            m.TraceName,
		"environment":          m.Environment,
		"region":               m.Region,
		"credentialAlias":      m.CredentialAlias,
		"retryNumber":          m.RetryNumber,
		"responseQualityScore": m.ResponseQualityScore,
		"videoJobId":           m.VideoJobID,
		"audioJobId":           m.AudioJobID,
		"subscriber":           m.Subscriber,
		"custom":               m.Custom,
	}, "", "  ")
	fmt.Println(string(jsonBytes))
}
