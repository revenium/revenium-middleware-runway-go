package main

import (
	"context"
	"fmt"
	"log"

	"github.com/revenium/revenium-middleware-runway-go/revenium"
)

func main() {
	fmt.Println("=== Revenium Runway Middleware - Basic Example ===")
	fmt.Println()

	// Initialize the middleware
	if err := revenium.Initialize(); err != nil {
		log.Fatalf("Failed to initialize middleware: %v", err)
	}

	// Get the client
	client, err := revenium.GetClient()
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	// Create context
	ctx := context.Background()

	// Define usage metadata
	metadata := &revenium.UsageMetadata{
		OrganizationID: "org-example-123",
		ProductID:      "product-runway-demo",
		TaskType:       "image-to-video-generation",
		Subscriber: map[string]interface{}{
			"id":    "user-123",
			"email": "demo@example.com",
		},
	}

	// Example 1: Image to Video
	fmt.Println("Example 1: Generating video from image...")
	fmt.Println("─────────────────────────────────────────")

	imageToVideoReq := &revenium.ImageToVideoRequest{
		PromptImage: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b6/Image_created_with_a_mobile_phone.png/1200px-Image_created_with_a_mobile_phone.png", // Valid public image
		PromptText:  "A cinematic shot of mountains at sunset",
		Model:       "gen3a_turbo",
		Duration:    5,
		Ratio:       "1280:768", // Resolution ratio (landscape)
	}

	result, err := client.ImageToVideo(ctx, imageToVideoReq, metadata)
	if err != nil {
		log.Fatalf("Failed to generate video: %v", err)
	}

	fmt.Printf("Task ID: %s\n", result.ID)
	fmt.Printf("Status: %s\n", result.Status)
	fmt.Printf("Model: %s\n", result.Model)
	fmt.Printf("Duration: %v\n", result.Duration)
	fmt.Printf("Output URLs: %v\n", result.OutputURLs)
	fmt.Println()

	// Example 2: Video Upscale (commented out - requires video URL)
	/*
	fmt.Println("Example 2: Upscaling video...")
	fmt.Println("─────────────────────────────────────────")

	upscaleReq := &revenium.VideoUpscaleRequest{
		PromptVideo: "https://example.com/sample-video.mp4", // Replace with actual video URL
		Model:       "upscale",
	}

	upscaleResult, err := client.UpscaleVideo(ctx, upscaleReq, metadata)
	if err != nil {
		log.Fatalf("Failed to upscale video: %v", err)
	}

	fmt.Printf("Task ID: %s\n", upscaleResult.ID)
	fmt.Printf("Status: %s\n", upscaleResult.Status)
	fmt.Printf("Output URLs: %v\n", upscaleResult.OutputURLs)
	fmt.Println()
	*/

	fmt.Println("Example completed successfully!")
	fmt.Println()
	fmt.Println("Note: Metering data has been sent to Revenium asynchronously.")
	fmt.Println("Check your Revenium dashboard at https://app.revenium.ai")
}
