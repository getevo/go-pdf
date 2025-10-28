package pdf

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/getevo/evo/v2"
	"github.com/getevo/evo/v2/lib/log"
)

type Controller struct{}

var (
	// Cache directory for temporary files
	cacheDir = "./cache"
	// Mutex for thread-safe file operations
	cacheMutex sync.RWMutex
	// Map to track file creation times
	fileCreationTimes = make(map[string]time.Time)
)

// init initializes the cache directory
func init() {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		log.Error("Failed to create cache directory: %v", err)
	}
}

// Generate handles POST /api/v1/generate
// It converts HTML to PDF using wkhtmltopdf
// Accepts HTML directly in the request body
func (c Controller) Generate(request *evo.Request) any {
	// Get HTML content directly from request body
	htmlContent := string(request.Body())

	// Validate HTML content
	if htmlContent == "" {
		request.Status(400)
		return map[string]any{
			"error":   "HTML content is required",
			"message": "Request body cannot be empty",
		}
	}

	// Generate unique filename based on HTML content hash
	hash := md5.Sum([]byte(htmlContent + time.Now().String()))
	filename := hex.EncodeToString(hash[:])

	htmlPath := filepath.Join(cacheDir, filename+".html")
	pdfPath := filepath.Join(cacheDir, filename+".pdf")

	// Write HTML to temporary file
	if err := os.WriteFile(htmlPath, []byte(htmlContent), 0644); err != nil {
		log.Error("Failed to write HTML file: %v", err)
		request.Status(500)
		return map[string]any{
			"error":   "Failed to write HTML file",
			"message": err.Error(),
		}
	}

	// Track file creation time
	cacheMutex.Lock()
	fileCreationTimes[htmlPath] = time.Now()
	fileCreationTimes[pdfPath] = time.Now()
	cacheMutex.Unlock()

	// Build wkhtmltopdf command with query parameters
	args := []string{
		"--enable-local-file-access",
		"--print-media-type",
		"--no-stop-slow-scripts",
	}

	// Parse query parameters
	// 1. JavaScript delay (default: 0, max: 5000ms)
	jsDelayStr := request.Query("js_delay").String()
	if jsDelayStr == "" {
		jsDelayStr = "0"
	} else {
		// Validate and cap at 5000ms
		if delay, err := strconv.Atoi(jsDelayStr); err == nil {
			if delay > 5000 {
				jsDelayStr = "5000"
			} else if delay < 0 {
				jsDelayStr = "0"
			}
		} else {
			jsDelayStr = "0"
		}
	}
	args = append(args, "--javascript-delay", jsDelayStr)

	// 2. Image DPI
	if imageDPI := request.Query("image_dpi").String(); imageDPI != "" {
		args = append(args, "--image-dpi", imageDPI)
	}

	// 3. Image Quality
	if imageQuality := request.Query("image_quality").String(); imageQuality != "" {
		args = append(args, "--image-quality", imageQuality)
	}

	// 4. Low Quality
	if lowQuality := request.Query("lowquality").String(); lowQuality == "true" {
		args = append(args, "--lowquality")
	}

	// 5. Page Height
	if pageHeight := request.Query("page_height").String(); pageHeight != "" {
		args = append(args, "--page-height", pageHeight)
	}

	// 6. Page Width
	if pageWidth := request.Query("page_width").String(); pageWidth != "" {
		args = append(args, "--page-width", pageWidth)
	}

	// 7. Page Size
	if pageSize := request.Query("page_size").String(); pageSize != "" {
		args = append(args, "--page-size", pageSize)
	}

	// 8. Enable Forms
	if enableForms := request.Query("enable_forms").String(); enableForms == "true" {
		args = append(args, "--enable-forms")
	}

	// 9. Enable Smart Shrinking (disabled by default, can be enabled)
	if enableShrinking := request.Query("enable_smart_shrinking").String(); enableShrinking == "true" {
		args = append(args, "--enable-smart-shrinking")
	} else {
		args = append(args, "--disable-smart-shrinking")
	}

	// 10. Margin Top
	if marginTop := request.Query("margin_top").String(); marginTop != "" {
		args = append(args, "--margin-top", marginTop)
	}

	// 11. Margin Bottom
	if marginBottom := request.Query("margin_bottom").String(); marginBottom != "" {
		args = append(args, "--margin-bottom", marginBottom)
	}

	// 12. Margin Left
	if marginLeft := request.Query("margin_left").String(); marginLeft != "" {
		args = append(args, "--margin-left", marginLeft)
	}

	// 13. Margin Right
	if marginRight := request.Query("margin_right").String(); marginRight != "" {
		args = append(args, "--margin-right", marginRight)
	}

	// 14. Orientation (Portrait or Landscape)
	if orientation := request.Query("orientation").String(); orientation != "" {
		// Validate orientation (case insensitive)
		switch orientation {
		case "portrait", "Portrait":
			args = append(args, "--orientation", "Portrait")
		case "landscape", "Landscape":
			args = append(args, "--orientation", "Landscape")
		}
	}

	// Add input and output paths
	args = append(args, htmlPath, pdfPath)

	// Convert HTML to PDF using wkhtmltopdf
	cmd := exec.Command("wkhtmltopdf", args...)

	// Capture stderr for debugging
	var stderr []byte
	var err error
	stderr, err = cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to generate PDF: %v, stderr: %s", err, string(stderr))
		// Clean up HTML file
		os.Remove(htmlPath)
		cacheMutex.Lock()
		delete(fileCreationTimes, htmlPath)
		delete(fileCreationTimes, pdfPath)
		cacheMutex.Unlock()

		request.Status(500)
		return map[string]any{
			"error":   "Failed to generate PDF",
			"message": fmt.Sprintf("%v: %s", err, string(stderr)),
		}
	}

	// Read the generated PDF
	pdfData, err := os.ReadFile(pdfPath)
	if err != nil {
		log.Error("Failed to read PDF file: %v", err)
		// Clean up files
		os.Remove(htmlPath)
		os.Remove(pdfPath)
		cacheMutex.Lock()
		delete(fileCreationTimes, htmlPath)
		delete(fileCreationTimes, pdfPath)
		cacheMutex.Unlock()

		request.Status(500)
		return map[string]any{
			"error":   "Failed to read PDF file",
			"message": err.Error(),
		}
	}

	// Set response headers
	request.Set("Content-Type", "application/pdf")
	request.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", filename))
	request.Set("X-Generated-At", time.Now().Format(time.RFC3339))

	log.Info("PDF generated successfully: %s", filename)

	// Return PDF file
	return pdfData
}

// Health handles GET /health
// Returns health status of the service
func (c Controller) Health(request *evo.Request) any {
	// Check if wkhtmltopdf is available
	cmd := exec.Command("wkhtmltopdf", "--version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		request.Status(503)
		return map[string]any{
			"status":  "unhealthy",
			"message": "wkhtmltopdf is not available",
			"error":   err.Error(),
		}
	}

	// Count cached files
	cacheMutex.RLock()
	cachedFiles := len(fileCreationTimes)
	cacheMutex.RUnlock()

	request.Status(200)
	return map[string]any{
		"status":        "healthy",
		"service":       "go-pdf",
		"wkhtmltopdf":   string(output),
		"cached_files":  cachedFiles,
		"timestamp":     time.Now().Format(time.RFC3339),
	}
}

// cleanupCache removes cached files older than 1 hour
func cleanupCache() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		cacheMutex.Lock()

		for filePath, creationTime := range fileCreationTimes {
			// Delete files older than 1 hour
			if now.Sub(creationTime) > 1*time.Hour {
				if err := os.Remove(filePath); err != nil {
					if !os.IsNotExist(err) {
						log.Error("Failed to delete cached file %s: %v", filePath, err)
					}
				} else {
					log.Debug("Deleted cached file: %s", filePath)
				}
				delete(fileCreationTimes, filePath)
			}
		}

		cacheMutex.Unlock()
	}
}
