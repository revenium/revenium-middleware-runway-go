// Package main demonstrates comprehensive metering field population for Runway Go middleware.
//
// This example populates ALL available metering fields with realistic enterprise values
// to verify:
// 1. All fields are transmitted correctly to Revenium
// 2. No values are hard-coded (compare with basic example)
// 3. Values represent sensible enterprise-customer scenarios
//
// Run with: go run main.go
// Requires: RUNWAY_API_KEY and REVENIUM_METERING_API_KEY environment variables
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

func main() {
	fmt.Println("=== Revenium Runway Middleware - Comprehensive Field Test ===")
	fmt.Println()
	fmt.Println("This example populates ALL available metering fields with realistic values.")
	fmt.Println("Compare the metering payload against the basic example to verify no hard-coding.")
	fmt.Println()

	// Enable verbose startup logging to see configuration details
	os.Setenv("REVENIUM_LOG_LEVEL", "DEBUG")
	os.Setenv("REVENIUM_VERBOSE_STARTUP", "true")

	// Initialize the middleware with options
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

	// Get the client
	client, err := revenium.GetClient()
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	// Create context with timeout (video generation can take several minutes)
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Minute)
	defer cancel()

	// Build COMPREHENSIVE metadata with ALL available fields
	// These values represent a realistic enterprise video production scenario
	metadata := buildComprehensiveMetadata()

	// Log the metadata we're sending
	fmt.Println("Comprehensive UsageMetadata being sent:")
	fmt.Println("----------------------------------------")
	printMetadata(metadata)
	fmt.Println()

	// Create video generation request with all options
	imageToVideoReq := &revenium.ImageToVideoRequest{
		PromptImage: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b6/Image_created_with_a_mobile_phone.png/1200px-Image_created_with_a_mobile_phone.png",
		PromptText:  "A sweeping cinematic drone shot revealing a futuristic cityscape at golden hour, volumetric lighting, 8K quality",
		Model:       "gen3a_turbo",
		Duration:    10, // Maximum duration for thorough testing
		Ratio:       "1280:768",
	}

	fmt.Println("Request Configuration:")
	fmt.Println("----------------------")
	fmt.Printf("  Model:    %s\n", imageToVideoReq.Model)
	fmt.Printf("  Duration: %d seconds\n", imageToVideoReq.Duration)
	fmt.Printf("  Ratio:    %s\n", imageToVideoReq.Ratio)
	fmt.Printf("  Prompt:   %s\n", imageToVideoReq.PromptText)
	fmt.Println()

	fmt.Println("Generating video... (this may take several minutes)")
	fmt.Println("===================================================")

	result, err := client.ImageToVideo(ctx, imageToVideoReq, metadata)
	if err != nil {
		log.Fatalf("Failed to generate video: %v", err)
	}

	// Display results
	fmt.Println()
	fmt.Println("Video Generation Complete!")
	fmt.Println("==========================")
	fmt.Printf("  Task ID:     %s\n", result.ID)
	fmt.Printf("  Status:      %s\n", result.Status)
	fmt.Printf("  Model:       %s\n", result.Model)
	fmt.Printf("  Duration:    %v\n", result.Duration)
	fmt.Printf("  Output URLs: %v\n", result.OutputURLs)
	fmt.Println()

	// Print summary of all metering fields that should have been transmitted
	fmt.Println("Expected Metering Payload Fields:")
	fmt.Println("==================================")
	printExpectedMeteringFields(metadata, result)

	fmt.Println()
	fmt.Println("Metering data sent asynchronously to Revenium.")
	fmt.Println("Check DEBUG logs above for actual payload, or verify in Revenium dashboard:")
	fmt.Println("  https://app.revenium.ai")
	fmt.Println()
	fmt.Println("Middleware version:", revenium.GetMiddlewareSource())
}

// buildComprehensiveMetadata creates UsageMetadata with ALL available fields populated
// with realistic enterprise values representing a video production company scenario.
func buildComprehensiveMetadata() *revenium.UsageMetadata {
	// RetryNumber must be a pointer
	retryNumber := 0

	// ResponseQualityScore must be a pointer (0.0-1.0 scale)
	qualityScore := 0.95

	return &revenium.UsageMetadata{
		// === ORGANIZATION & PRODUCT IDENTIFICATION ===
		// These identify the customer organization and which product/feature this usage belongs to
		OrganizationID: "org-videotech-studios-prod",
		ProductID:      "prod-ai-video-gen-platform-v2",

		// === TASK IDENTIFICATION ===
		// Describes what type of work is being performed
		TaskType: "marketing-campaign-hero-video",
		TaskID:   "task-2026-q1-brand-refresh-001",

		// === AGENT/WORKER IDENTIFICATION ===
		// Identifies which worker/service instance is making this request
		Agent: "video-rendering-worker-03-eu",

		// === SUBSCRIPTION & BILLING ===
		// Links this usage to a specific subscription for billing purposes
		SubscriptionID: "sub-enterprise-unlimited-annual-2026",

		// === DISTRIBUTED TRACING ===
		// For correlating requests across distributed systems
		TraceID:             "trace-vid-abc123-def456-789xyz",
		ParentTransactionID: "parent-tx-campaign-workflow-main",
		TraceType:           "distributed",
		TraceName:           "marketing-video-generation-pipeline",

		// === DEPLOYMENT CONTEXT ===
		// Where this request is running
		Environment: "production",
		Region:      "eu-central-1",

		// === CREDENTIAL MANAGEMENT ===
		// Which credential alias is being used (for multi-key setups)
		CredentialAlias: "runway-prod-key-primary-eu",

		// === RETRY TRACKING ===
		// Track retry attempts for reliability analysis
		RetryNumber: &retryNumber,

		// === QUALITY METRICS ===
		// Response quality scoring (0.0-1.0)
		ResponseQualityScore: &qualityScore,

		// === MULTIMODAL JOB IDENTIFIERS ===
		// For tracking related video/audio jobs in complex pipelines
		VideoJobID: "vjob-brand-refresh-hero-2026-001",
		AudioJobID: "ajob-brand-refresh-soundtrack-2026-001",

		// === SUBSCRIBER INFORMATION ===
		// Detailed information about the end-user making this request
		Subscriber: map[string]interface{}{
			"id":                  "user-enterprise-admin-12345",
			"email":               "video.producer@videotech-studios.com",
			"name":                "Alexandra Chen",
			"role":                "Senior Video Producer",
			"department":          "Marketing Content Production",
			"costCenter":          "CC-MKT-CONTENT-2026",
			"projectCode":         "PRJ-BRAND-REFRESH-Q1",
			"internalClientId":    "client-internal-marketing-team",
			"quotaGroup":          "enterprise-unlimited",
			"billingTier":         "premium",
			"accountCreatedAt":    "2023-06-15T10:30:00Z",
			"lastLoginAt":         "2026-01-22T08:45:00Z",
			"totalVideosGenerated": 1247,
			"monthlyBudgetUsed":   78.5,
			"monthlyBudgetLimit":  500.0,
		},

		// === CUSTOM FIELDS ===
		// Business-specific metadata that doesn't fit standard fields
		Custom: map[string]interface{}{
			// Campaign tracking
			"campaignId":        "camp-brand-refresh-2026-q1",
			"campaignName":      "Q1 2026 Brand Refresh Initiative",
			"campaignPhase":     "hero-content-creation",
			"deliverableType":   "hero-video-social-media",
			"targetPlatforms":   []string{"youtube", "linkedin", "instagram", "tiktok"},

			// Content classification
			"contentCategory":   "brand-marketing",
			"contentRating":     "general-audience",
			"brandGuidelines":   "v2.3-jan-2026",
			"stylePreset":       "corporate-cinematic",

			// Approval workflow
			"approvalRequired":  true,
			"approverEmail":     "creative.director@videotech-studios.com",
			"approvalDeadline":  "2026-01-25T18:00:00Z",
			"reviewCycleNumber": 1,

			// Asset management
			"assetLibraryRef":   "asset-lib-brand-2026",
			"sourceAssetId":     "img-hero-cityscape-001",
			"outputFormat":      "mp4-h265-4k",
			"outputDestination": "s3://videotech-assets/campaigns/brand-refresh-2026/",

			// Analytics tags
			"analyticsEnabled":  true,
			"abTestVariant":     "cinematic-style-a",
			"experimentId":      "exp-video-style-comparison-001",

			// Cost tracking
			"budgetCode":        "BUD-MKT-2026-Q1-VIDEO",
			"costAllocation":    "marketing-content-production",
			"invoiceReference":  "INV-2026-01-VIDEO-PROD",

			// Technical metadata
			"requestPriority":   "high",
			"queueName":         "enterprise-priority-eu",
			"processingNode":    "gpu-node-eu-central-03",
			"gpuType":           "nvidia-a100-80gb",

			// Compliance
			"dataResidency":     "eu",
			"gdprCompliant":     true,
			"retentionPolicy":   "90-days-then-archive",
		},
	}
}

// printMetadata pretty-prints the metadata structure
func printMetadata(m *revenium.UsageMetadata) {
	fmt.Printf("  OrganizationID:       %s\n", m.OrganizationID)
	fmt.Printf("  ProductID:            %s\n", m.ProductID)
	fmt.Printf("  TaskType:             %s\n", m.TaskType)
	fmt.Printf("  TaskID:               %s\n", m.TaskID)
	fmt.Printf("  Agent:                %s\n", m.Agent)
	fmt.Printf("  SubscriptionID:       %s\n", m.SubscriptionID)
	fmt.Printf("  TraceID:              %s\n", m.TraceID)
	fmt.Printf("  ParentTransactionID:  %s\n", m.ParentTransactionID)
	fmt.Printf("  TraceType:            %s\n", m.TraceType)
	fmt.Printf("  TraceName:            %s\n", m.TraceName)
	fmt.Printf("  Environment:          %s\n", m.Environment)
	fmt.Printf("  Region:               %s\n", m.Region)
	fmt.Printf("  CredentialAlias:      %s\n", m.CredentialAlias)
	if m.RetryNumber != nil {
		fmt.Printf("  RetryNumber:          %d\n", *m.RetryNumber)
	}
	if m.ResponseQualityScore != nil {
		fmt.Printf("  ResponseQualityScore: %.2f\n", *m.ResponseQualityScore)
	}
	fmt.Printf("  VideoJobID:           %s\n", m.VideoJobID)
	fmt.Printf("  AudioJobID:           %s\n", m.AudioJobID)

	if m.Subscriber != nil {
		fmt.Printf("  Subscriber:           (%d fields)\n", len(m.Subscriber))
		for k, v := range m.Subscriber {
			fmt.Printf("    - %s: %v\n", k, v)
		}
	}

	if m.Custom != nil {
		fmt.Printf("  Custom:               (%d fields)\n", len(m.Custom))
		for k, v := range m.Custom {
			fmt.Printf("    - %s: %v\n", k, v)
		}
	}
}

// printExpectedMeteringFields shows all fields that should appear in the metering payload
func printExpectedMeteringFields(m *revenium.UsageMetadata, result *revenium.VideoGenerationResult) {
	fmt.Println()
	fmt.Println("=== MIDDLEWARE-POPULATED FIELDS ===")
	fmt.Println("(Automatically set by the middleware)")
	fmt.Printf("  operationType:            VIDEO\n")
	fmt.Printf("  provider:                 runway\n")
	fmt.Printf("  modelSource:              RUNWAY\n")
	fmt.Printf("  model:                    %s\n", result.Model)
	fmt.Printf("  transactionId:            %s\n", result.ID)
	fmt.Printf("  requestTime:              (auto-calculated)\n")
	fmt.Printf("  responseTime:             (auto-calculated)\n")
	fmt.Printf("  requestDuration:          %d ms\n", result.Duration.Milliseconds())
	fmt.Printf("  durationSeconds:          (from video metadata)\n")
	fmt.Printf("  requestedDurationSeconds: 10 (from request)\n")
	fmt.Printf("  stopReason:               END\n")
	fmt.Printf("  costType:                 AI\n")
	fmt.Printf("  isStreamed:               false\n")
	fmt.Printf("  middlewareSource:         %s\n", revenium.GetMiddlewareSource())

	fmt.Println()
	fmt.Println("=== PROMPT CAPTURE FIELDS (opt-in) ===")
	fmt.Println("(Added when CapturePrompts is enabled)")
	fmt.Printf("  inputMessages:            [{\"role\":\"user\",\"content\":\"<prompt>\"}]\n")
	fmt.Printf("  outputResponse:           (JSON array of video URLs)\n")
	fmt.Printf("  promptsTruncated:         (true if prompt exceeded 50K chars)\n")

	fmt.Println()
	fmt.Println("=== USER-PROVIDED METADATA FIELDS ===")
	fmt.Println("(Set via UsageMetadata struct)")
	fmt.Printf("  organizationId:           %s\n", m.OrganizationID)
	fmt.Printf("  productId:                %s\n", m.ProductID)
	fmt.Printf("  taskType:                 %s\n", m.TaskType)
	fmt.Printf("  taskId:                   %s\n", m.TaskID)
	fmt.Printf("  agent:                    %s\n", m.Agent)
	fmt.Printf("  subscriptionId:           %s\n", m.SubscriptionID)
	fmt.Printf("  traceId:                  %s\n", m.TraceID)
	fmt.Printf("  parentTransactionId:      %s\n", m.ParentTransactionID)
	fmt.Printf("  traceType:                %s\n", m.TraceType)
	fmt.Printf("  traceName:                %s\n", m.TraceName)
	fmt.Printf("  environment:              %s\n", m.Environment)
	fmt.Printf("  region:                   %s\n", m.Region)
	fmt.Printf("  credentialAlias:          %s\n", m.CredentialAlias)
	if m.RetryNumber != nil {
		fmt.Printf("  retryNumber:              %d\n", *m.RetryNumber)
	}
	if m.ResponseQualityScore != nil {
		fmt.Printf("  responseQualityScore:     %.2f\n", *m.ResponseQualityScore)
	}
	fmt.Printf("  videoJobId:               %s\n", m.VideoJobID)
	fmt.Printf("  audioJobId:               %s\n", m.AudioJobID)
	fmt.Printf("  subscriber:               (%d nested fields)\n", len(m.Subscriber))
	fmt.Printf("  custom:                   (%d nested fields, merged at top level)\n", len(m.Custom))

	// Generate expected payload for verification
	fmt.Println()
	fmt.Println("=== EXPECTED METERING PAYLOAD (JSON preview) ===")
	payload := buildExpectedPayload(m, result)
	jsonBytes, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(jsonBytes))
}

// buildExpectedPayload constructs what we expect the metering payload to look like
func buildExpectedPayload(m *revenium.UsageMetadata, result *revenium.VideoGenerationResult) map[string]interface{} {
	payload := map[string]interface{}{
		// Middleware-populated fields
		"operationType":            "VIDEO",
		"provider":                 "runway",
		"modelSource":              "RUNWAY",
		"model":                    result.Model,
		"transactionId":            result.ID,
		"requestTime":              "(auto)",
		"responseTime":             "(auto)",
		"requestDuration":          result.Duration.Milliseconds(),
		"durationSeconds":          5.0, // Default
		"requestedDurationSeconds": 10,  // From request
		"stopReason":               "END",
		"costType":                 "AI",
		"isStreamed":               false,
		"middlewareSource":         revenium.GetMiddlewareSource(),

		// User metadata fields
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
		"retryNumber":          *m.RetryNumber,
		"responseQualityScore": *m.ResponseQualityScore,
		"videoJobId":           m.VideoJobID,
		"audioJobId":           m.AudioJobID,
		"subscriber":           m.Subscriber,
	}

	// Custom fields are merged at top level
	for k, v := range m.Custom {
		payload[k] = v
	}

	return payload
}
