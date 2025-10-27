// Package main is the entry point for the go-pdf application
package main

import (
	"github.com/getevo/evo/v2"
	"github.com/getevo/evo/v2/lib/application"
	"go-pdf/apps/pdf"
)

// main is the entry point function that initializes and runs the application
func main() {
	// Initialize the Evo framework
	evo.Setup()

	// Get the application instance
	var apps = application.GetInstance()

	// Register the PDF app
	apps.Register(
		pdf.App{}, // PDF generation service
	)

	// Start the application
	evo.Run()
}