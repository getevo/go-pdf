package pdf

import (
	"github.com/getevo/evo/v2"
	"github.com/getevo/evo/v2/lib/log"
)

type App struct{}

// Register initializes the PDF app
func (App) Register() error {
	log.Info("PDF app registered")
	return nil
}

// Router initializes the routes for the PDF generation service
func (App) Router() error {
	var controller Controller

	// POST endpoint for PDF generation
	evo.Post("/api/v1/generate", controller.Generate)

	log.Info("PDF routes registered: POST /api/v1/generate")
	return nil
}

// WhenReady is called when the application is ready
func (App) WhenReady() error {
	// Start the cache cleanup routine
	go cleanupCache()
	log.Info("PDF cache cleanup routine started")
	return nil
}

// Name returns the name of the app
func (App) Name() string {
	return "pdf"
}
