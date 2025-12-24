package revenium

import (
	"context"
	"sync"
	"time"
)

// ReveniumRunway is the main middleware client that wraps Runway API
// and adds metering capabilities
type ReveniumRunway struct {
	runwayClient   *RunwayClient
	meteringClient *MeteringClient
	config         *Config
	mu             sync.RWMutex
}

var (
	globalClient *ReveniumRunway
	globalMu     sync.RWMutex
	initialized  bool
)

// Initialize sets up the global Revenium middleware with configuration
func Initialize(opts ...Option) error {
	globalMu.Lock()
	defer globalMu.Unlock()

	if initialized {
		return nil
	}

	// Initialize logger first
	InitializeLogger()
	Info("Initializing Revenium Runway middleware...")

	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Load from environment if not provided
	if err := cfg.LoadFromEnv(); err != nil {
		Warn("Failed to load configuration from environment: %v", err)
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Create clients
	runwayClient := NewRunwayClient(cfg)
	meteringClient := NewMeteringClient(cfg)

	globalClient = &ReveniumRunway{
		runwayClient:   runwayClient,
		meteringClient: meteringClient,
		config:         cfg,
	}

	initialized = true
	Info("Revenium Runway middleware initialized successfully")
	return nil
}

// IsInitialized checks if the middleware is properly initialized
func IsInitialized() bool {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return initialized
}

// GetClient returns the global Revenium client
func GetClient() (*ReveniumRunway, error) {
	globalMu.RLock()
	defer globalMu.RUnlock()

	if !initialized {
		return nil, NewConfigError("middleware not initialized, call Initialize() first", nil)
	}

	return globalClient, nil
}

// NewReveniumRunway creates a new Revenium client with explicit configuration
func NewReveniumRunway(cfg *Config) (*ReveniumRunway, error) {
	if cfg == nil {
		return nil, NewConfigError("config cannot be nil", nil)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	runwayClient := NewRunwayClient(cfg)
	meteringClient := NewMeteringClient(cfg)

	return &ReveniumRunway{
		runwayClient:   runwayClient,
		meteringClient: meteringClient,
		config:         cfg,
	}, nil
}

// GetConfig returns the configuration
func (r *ReveniumRunway) GetConfig() *Config {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.config
}

// ImageToVideo generates a video from an image with automatic metering
func (r *ReveniumRunway) ImageToVideo(ctx context.Context, req *ImageToVideoRequest, metadata *UsageMetadata) (*VideoGenerationResult, error) {
	startTime := time.Now()

	// Set default model if not specified
	if req.Model == "" {
		req.Model = "gen3a_turbo"
	}

	// Create task
	Debug("Creating image-to-video task with model: %s", req.Model)
	taskResp, err := r.runwayClient.CreateImageToVideo(ctx, req)
	if err != nil {
		return nil, err
	}

	// Wait for task completion
	Info("Waiting for task %s to complete...", taskResp.ID)
	statusResp, err := r.runwayClient.WaitForTaskCompletion(ctx, taskResp.ID, DefaultPollingConfig())
	if err != nil {
		return nil, err
	}

	// Build result
	duration := time.Since(startTime)
	result := &VideoGenerationResult{
		ID:         taskResp.ID,
		Status:     statusResp.Status,
		OutputURLs: statusResp.Output,
		Duration:   duration,
		Model:      req.Model,
	}

	// Copy error information if failed
	if statusResp.Error != nil {
		result.Error = statusResp.Error
	}
	if statusResp.FailureCode != nil {
		result.FailureCode = statusResp.FailureCode
	}

	// Send metering asynchronously (fire-and-forget)
	go r.sendMetering(context.Background(), result, metadata)

	return result, nil
}

// VideoToVideo transforms a video with automatic metering
func (r *ReveniumRunway) VideoToVideo(ctx context.Context, req *VideoToVideoRequest, metadata *UsageMetadata) (*VideoGenerationResult, error) {
	startTime := time.Now()

	// Set default model if not specified
	if req.Model == "" {
		req.Model = "gen3a_turbo"
	}

	// Create task
	Debug("Creating video-to-video task with model: %s", req.Model)
	taskResp, err := r.runwayClient.CreateVideoToVideo(ctx, req)
	if err != nil {
		return nil, err
	}

	// Wait for task completion
	Info("Waiting for task %s to complete...", taskResp.ID)
	statusResp, err := r.runwayClient.WaitForTaskCompletion(ctx, taskResp.ID, DefaultPollingConfig())
	if err != nil {
		return nil, err
	}

	// Build result
	duration := time.Since(startTime)
	result := &VideoGenerationResult{
		ID:         taskResp.ID,
		Status:     statusResp.Status,
		OutputURLs: statusResp.Output,
		Duration:   duration,
		Model:      req.Model,
	}

	// Copy error information if failed
	if statusResp.Error != nil {
		result.Error = statusResp.Error
	}
	if statusResp.FailureCode != nil {
		result.FailureCode = statusResp.FailureCode
	}

	// Send metering asynchronously (fire-and-forget)
	go r.sendMetering(context.Background(), result, metadata)

	return result, nil
}

// UpscaleVideo upscales a video with automatic metering
func (r *ReveniumRunway) UpscaleVideo(ctx context.Context, req *VideoUpscaleRequest, metadata *UsageMetadata) (*VideoGenerationResult, error) {
	startTime := time.Now()

	// Set default model if not specified
	if req.Model == "" {
		req.Model = "upscale"
	}

	// Create task
	Debug("Creating video upscale task with model: %s", req.Model)
	taskResp, err := r.runwayClient.CreateVideoUpscale(ctx, req)
	if err != nil {
		return nil, err
	}

	// Wait for task completion
	Info("Waiting for task %s to complete...", taskResp.ID)
	statusResp, err := r.runwayClient.WaitForTaskCompletion(ctx, taskResp.ID, DefaultPollingConfig())
	if err != nil {
		return nil, err
	}

	// Build result
	duration := time.Since(startTime)
	result := &VideoGenerationResult{
		ID:         taskResp.ID,
		Status:     statusResp.Status,
		OutputURLs: statusResp.Output,
		Duration:   duration,
		Model:      req.Model,
	}

	// Copy error information if failed
	if statusResp.Error != nil {
		result.Error = statusResp.Error
	}
	if statusResp.FailureCode != nil {
		result.FailureCode = statusResp.FailureCode
	}

	// Send metering asynchronously (fire-and-forget)
	go r.sendMetering(context.Background(), result, metadata)

	return result, nil
}

// sendMetering sends metering data asynchronously
func (r *ReveniumRunway) sendMetering(ctx context.Context, result *VideoGenerationResult, metadata *UsageMetadata) {
	defer func() {
		if rec := recover(); rec != nil {
			Error("Metering goroutine panic: %v", rec)
		}
	}()

	if err := r.meteringClient.SendVideoMetering(ctx, result, metadata); err != nil {
		Error("Failed to send metering data: %v", err)
	}
}

// Close closes the client and cleans up resources
func (r *ReveniumRunway) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.runwayClient.Close(); err != nil {
		return err
	}
	if err := r.meteringClient.Close(); err != nil {
		return err
	}

	return nil
}

// Reset resets the global middleware state for testing
func Reset() {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalClient != nil {
		globalClient.Close()
		globalClient = nil
	}

	initialized = false
}
