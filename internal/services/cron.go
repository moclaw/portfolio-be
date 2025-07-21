package services

import (
	"context"
	"log"
	"time"
)

// CronService handles scheduled tasks
type CronService struct {
	ctx             context.Context
	cancel          context.CancelFunc
	resourceService *ResourceService
	uploadService   *UploadService
	ticker          *time.Ticker
}

// NewCronService creates a new cron service
func NewCronService(resourceService *ResourceService, uploadService *UploadService) *CronService {
	ctx, cancel := context.WithCancel(context.Background())
	return &CronService{
		ctx:             ctx,
		cancel:          cancel,
		resourceService: resourceService,
		uploadService:   uploadService,
	}
}

// Start starts the cron service with all scheduled jobs
func (cs *CronService) Start() error {
	// Run jobs every hour for testing, change to 24 hours for production
	cs.ticker = time.NewTicker(1 * time.Hour) // Change to 24 * time.Hour for daily execution

	go cs.runScheduledJobs()

	log.Println("Cron service started successfully")
	return nil
}

// Stop stops the cron service
func (cs *CronService) Stop() {
	if cs.ticker != nil {
		cs.ticker.Stop()
	}
	cs.cancel()
	log.Println("Cron service stopped")
}

// runScheduledJobs runs the scheduled jobs
func (cs *CronService) runScheduledJobs() {
	for {
		select {
		case <-cs.ctx.Done():
			return
		case <-cs.ticker.C:
			cs.RefreshExpiredURLsJob()
			cs.CleanupExpiredUploadsJob()
		}
	}
}

// RefreshExpiredURLsJob is the job that refreshes expired URLs
func (cs *CronService) RefreshExpiredURLsJob() {
	log.Println("Starting URL refresh job...")

	start := time.Now()
	err := cs.resourceService.RefreshExpiredURLs()
	if err != nil {
		log.Printf("Error during URL refresh job: %v", err)
		return
	}

	duration := time.Since(start)
	log.Printf("URL refresh job completed successfully in %v", duration)
}

// CleanupExpiredUploadsJob is the job that cleans up truly expired uploads
func (cs *CronService) CleanupExpiredUploadsJob() {
	log.Println("Starting cleanup expired uploads job...")

	start := time.Now()
	err := cs.cleanupExpiredUploads()
	if err != nil {
		log.Printf("Error during cleanup job: %v", err)
		return
	}

	duration := time.Since(start)
	log.Printf("Cleanup job completed successfully in %v", duration)
}

// cleanupExpiredUploads removes uploads that have been expired for more than 30 days
func (cs *CronService) cleanupExpiredUploads() error {
	// This is a placeholder for cleanup logic
	// You might want to:
	// 1. Find uploads expired for more than 30 days
	// 2. Remove from S3
	// 3. Mark as inactive in database
	// 4. Optionally remove from database completely

	log.Println("Cleanup logic would run here...")
	// Implementation would depend on your business requirements

	return nil
}

// RunURLRefreshNow runs the URL refresh job immediately (for testing/manual trigger)
func (cs *CronService) RunURLRefreshNow() error {
	log.Println("Running URL refresh job manually...")
	cs.RefreshExpiredURLsJob()
	return nil
}

// GetStatus returns the current status of the cron service
func (cs *CronService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"running":   cs.ctx.Err() == nil,
		"next_tick": time.Now().Add(1 * time.Hour), // Approximate next run time
		"last_run":  time.Now(),                    // This should be tracked properly in production
	}
}
