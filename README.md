# go-pdf Service

[![Build and Publish Docker Image](https://github.com/getevo/go-pdf/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/getevo/go-pdf/actions/workflows/docker-publish.yml)
[![Docker Image](https://img.shields.io/badge/docker-ghcr.io%2Fgetevo%2Fgo--pdf-blue)](https://ghcr.io/getevo/go-pdf)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance, concurrent HTML to PDF conversion service built with Go and the [getevo/evo/v2](https://github.com/getevo/evo) framework. This service provides a simple REST API to convert HTML (with CSS and JavaScript) into PDF documents using wkhtmltopdf.

**Docker Image**: `ghcr.io/getevo/go-pdf:latest`

## Why This Service Matters

Imagine you need to generate an invoice or receipt in your application. What do you do?

**Option 1**: Create messy, unreadable code embedded in your application that requires recompiling the entire app for each design change?

**Option 2**: Use an internal service that takes an HTML template and generates PDFs for you?

**This is why go-pdf matters!** It lets you painlessly generate any PDF with your custom style. Just design your document as HTML/CSS (like you would for a webpage), send it to this service, and get back a professional PDF. No more wrestling with complex PDF libraries or maintaining rigid PDF generation code in your main application.

### The Developer Experience

- 🎨 **Design PDFs like webpages**: Use familiar HTML/CSS
- 🔄 **Iterate quickly**: Change designs without recompiling your app
- 📦 **Keep it separate**: PDF generation logic stays out of your main codebase
- 🚀 **Scale independently**: Deploy PDF service separately from your application
- 🎯 **Reuse templates**: Same HTML template for web display and PDF generation

## Features

- **Fast & Efficient**: Built with Go for high performance and low memory footprint
- **Concurrent Processing**: Handles multiple concurrent PDF generation requests
- **REST API**: Simple POST endpoint for PDF generation with configurable options
- **Automatic Cleanup**: Automatically removes cached files older than 1 hour
- **Docker Support**: Multi-stage Dockerfile for minimal image size
- **No Database Required**: Lightweight service with no database dependencies
- **Framework-based**: Uses getevo/evo/v2 framework for robust HTTP handling
- **Flexible Configuration**: Control page size, quality, JavaScript execution, and more via query parameters

## What is go-pdf Service?

go-pdf is a microservice that converts HTML documents (including CSS and JavaScript) into PDF files. It uses wkhtmltopdf under the hood, which renders HTML using the WebKit engine and generates high-quality PDF output.

The service temporarily stores HTML files and generated PDFs in a cache directory, then automatically cleans up files older than 1 hour to prevent disk space issues.

### Use Cases

- Generate invoices and receipts from HTML templates
- Create PDF reports from web content
- Convert HTML emails to PDF
- Generate printable documents from web applications
- Batch PDF generation from dynamic HTML content

## How to Use go-pdf

### API Endpoints

#### **POST /api/v1/generate**

Converts HTML content to PDF and returns the generated PDF file.

#### Request

- **Method**: POST
- **Content-Type**: text/html
- **Body**: Raw HTML content (string)

#### Response

- **Content-Type**: application/pdf
- **Headers**:
  - `Content-Disposition`: attachment; filename=`{hash}`.pdf
  - `X-Generated-At`: Timestamp of generation
- **Body**: PDF file (binary)

#### Query Parameters

All query parameters are optional and allow you to control the PDF generation behavior:

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `js_delay` | integer | `0` | `5000` | JavaScript delay in milliseconds before rendering. Use this if your HTML contains JavaScript that needs time to execute |
| `image_dpi` | integer | - | - | Set the DPI for images (e.g., `300` for high quality) |
| `image_quality` | integer | - | `100` | JPEG image quality (0-100, where 100 is best quality) |
| `lowquality` | boolean | `false` | - | Use `true` to generate lower quality PDFs (smaller file size) |
| `page_height` | string | - | - | Page height (e.g., `297mm`, `11.69in`) |
| `page_width` | string | - | - | Page width (e.g., `210mm`, `8.27in`) |
| `page_size` | string | - | - | Page size preset (e.g., `A4`, `Letter`, `Legal`) |
| `enable_forms` | boolean | `false` | - | Use `true` to enable HTML forms in the PDF |
| `enable_smart_shrinking` | boolean | `false` | - | Use `true` to enable smart content shrinking to fit page |
| `margin_top` | string | - | - | Top margin (e.g., `10mm`, `0.5in`) |
| `margin_bottom` | string | - | - | Bottom margin (e.g., `10mm`, `0.5in`) |
| `margin_left` | string | - | - | Left margin (e.g., `10mm`, `0.5in`) |
| `margin_right` | string | - | - | Right margin (e.g., `10mm`, `0.5in`) |
| `orientation` | string | - | - | Page orientation: `portrait` or `landscape` |

**Example with query parameters:**

```bash
# Generate PDF with JavaScript delay and custom page size
curl -X POST "http://localhost:8080/api/v1/generate?js_delay=2000&page_size=A4&image_quality=95" \
  -H "Content-Type: text/html" \
  -d '<html><body><h1>High Quality PDF</h1></body></html>' \
  -o output.pdf
```

```bash
# Generate PDF with custom page dimensions
curl -X POST "http://localhost:8080/api/v1/generate?page_width=210mm&page_height=297mm" \
  -H "Content-Type: text/html" \
  -d '<html><body><h1>Custom Size PDF</h1></body></html>' \
  -o custom.pdf
```

```bash
# Generate low quality PDF (smaller file size)
curl -X POST "http://localhost:8080/api/v1/generate?lowquality=true&image_quality=50" \
  -H "Content-Type: text/html" \
  -d '<html><body><h1>Compressed PDF</h1></body></html>' \
  -o compressed.pdf
```

```bash
# Generate landscape PDF with custom margins
curl -X POST "http://localhost:8080/api/v1/generate?orientation=landscape&margin_top=20mm&margin_bottom=20mm&margin_left=15mm&margin_right=15mm" \
  -H "Content-Type: text/html" \
  -d '<html><body><h1>Landscape Report</h1><p>This is a wide format document.</p></body></html>' \
  -o landscape.pdf
```

#### Example Usage with curl

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: text/html" \
  -d '<html><body><h1>Hello PDF!</h1></body></html>' \
  -o output.pdf
```

**With styled HTML:**

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: text/html" \
  -d '<!DOCTYPE html>
<html>
<head>
  <title>Test</title>
  <style>
    body { font-family: Arial; margin: 40px; }
    h1 { color: #333; }
  </style>
</head>
<body>
  <h1>Hello PDF!</h1>
  <p>This is a test document.</p>
</body>
</html>' \
  -o output.pdf
```

#### Example Usage with JavaScript/Fetch

```javascript
const html = `
<!DOCTYPE html>
<html>
<head>
    <title>Invoice</title>
    <style>
        body { font-family: Arial, sans-serif; }
        .header { background-color: #4CAF50; color: white; padding: 20px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Invoice #12345</h1>
    </div>
    <p>Thank you for your purchase!</p>
</body>
</html>
`;

fetch('http://localhost:8080/api/v1/generate', {
    method: 'POST',
    headers: {
        'Content-Type': 'text/html'
    },
    body: html
})
.then(response => response.blob())
.then(blob => {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'invoice.pdf';
    a.click();
});
```

#### Example Usage with Python

```python
import requests

html_content = """
<!DOCTYPE html>
<html>
<head>
    <title>Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
    </style>
</head>
<body>
    <h1>Monthly Report</h1>
    <p>This is the monthly report generated on demand.</p>
</body>
</html>
"""

response = requests.post(
    'http://localhost:8080/api/v1/generate',
    headers={'Content-Type': 'text/html'},
    data=html_content
)

if response.status_code == 200:
    with open('report.pdf', 'wb') as f:
        f.write(response.content)
    print('PDF generated successfully!')
else:
    print(f'Error: {response.status_code}')
```

#### **GET /health**

Health check endpoint for monitoring service status and readiness.

**Response (HTTP 200 - Healthy):**

```json
{
  "status": "healthy",
  "service": "go-pdf",
  "wkhtmltopdf": "wkhtmltopdf 0.12.6 (with patched qt)\n",
  "cached_files": 4,
  "timestamp": "2025-10-28T00:00:00Z"
}
```

**Response (HTTP 503 - Unhealthy):**

```json
{
  "status": "unhealthy",
  "message": "wkhtmltopdf is not available",
  "error": "exec: \"wkhtmltopdf\": executable file not found"
}
```

**Example Usage:**

```bash
# Check service health
curl http://localhost:8080/health
```

```python
# Python health check
import requests

response = requests.get('http://localhost:8080/health')
if response.status_code == 200:
    print('Service is healthy')
    print(f"Cached files: {response.json()['cached_files']}")
else:
    print('Service is unhealthy')
```

## How to Install Using Docker

### Prerequisites

- Docker Desktop installed on your PC
- Docker daemon running

### Method 1: Using Pre-built Image (Recommended)

Pull and run the pre-built image from GitHub Container Registry:

```bash
docker run -d -p 8080:8080 --name go-pdf ghcr.io/getevo/go-pdf:latest
```

That's it! The service will be available at `http://localhost:8080`

**With docker-compose:**

```yaml
version: '3.8'
services:
  go-pdf:
    image: ghcr.io/getevo/go-pdf:latest
    container_name: go-pdf
    ports:
      - "8080:8080"
    restart: unless-stopped
```

Then run:
```bash
docker-compose up -d
```

### Method 2: Build from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/getevo/go-pdf.git
   cd go-pdf
   ```

2. **Build the Docker image**:
   ```bash
   docker build -t go-pdf:latest .
   ```

3. **Run the container**:
   ```bash
   docker run -d -p 8080:8080 --name go-pdf go-pdf:latest
   ```

4. **Verify the service is running**:
   ```bash
   curl -X POST http://localhost:8080/api/v1/generate \
     -H "Content-Type: text/html" \
     -d '<html><body><h1>Test</h1></body></html>' \
     -o test.pdf
   ```

5. **View logs**:
   ```bash
   docker logs go-pdf
   ```

6. **Stop the container**:
   ```bash
   docker stop go-pdf
   ```

### Method 3: Using Docker Compose (Build from Source)

1. **Clone and start**:
   ```bash
   git clone https://github.com/getevo/go-pdf.git
   cd go-pdf
   docker-compose up -d
   ```

2. **View logs**:
   ```bash
   docker-compose logs -f
   ```

3. **Stop the service**:
   ```bash
   docker-compose down
   ```

### Configuration

The service is configured via `config.yml`. Key settings:

- **HTTP.Port**: 8080 (default) - The port the service listens on
- **HTTP.BodyLimit**: 50mb - Maximum request body size
- **HTTP.ReadTimeout**: 30s - Request read timeout
- **HTTP.WriteTimeout**: 60s - Response write timeout
- **Database.Enabled**: false - Database is disabled

### Environment Variables

You can override configuration using environment variables or by modifying the `config.yml` file before building the Docker image.

## Architecture

```
go-pdf/
├── main.go              # Application entry point
├── config.yml           # Configuration file
├── apps/
│   └── pdf/
│       ├── app.go       # App registration and routing
│       └── controller.go # PDF generation logic
├── cache/               # Temporary file storage (auto-created)
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Docker Compose configuration
└── README.md            # This file
```

### How It Works

1. **Request**: Client sends HTML content via POST to `/api/v1/generate`
2. **Validation**: Service validates the HTML content
3. **File Creation**: HTML is written to a temporary file in the cache directory
4. **Conversion**: wkhtmltopdf converts HTML to PDF
5. **Response**: PDF is read and sent back to the client
6. **Cleanup**: Files older than 1 hour are automatically deleted every 10 minutes

### Performance Considerations

- **Concurrent Requests**: The service is designed to handle concurrent requests efficiently
- **File Caching**: Temporary files are cached for 1 hour to balance between reusability and disk space
- **Memory Usage**: The service uses minimal memory thanks to Go's efficiency and streaming responses
- **Image Size**: Multi-stage Docker build results in a minimal image size (~400MB including wkhtmltopdf dependencies)

## Troubleshooting

### Container Won't Start

Check logs:
```bash
docker logs go-pdf
```

### PDF Generation Fails

1. Verify HTML is valid
2. Check if wkhtmltopdf is installed in the container:
   ```bash
   docker exec go-pdf wkhtmltopdf --version
   ```

### Port Already in Use

Change the port mapping:
```bash
docker run -d -p 9090:8080 --name go-pdf go-pdf:latest
```

## Development

### Building Locally

```bash
go mod download
go build -o go-pdf .
./go-pdf
```

### Testing

```bash
# Create a test HTML file
echo '<!DOCTYPE html><html><body><h1>Test</h1></body></html>' > test.html

# Generate PDF
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d "{\"html\":\"$(cat test.html)\"}" \
  -o test.pdf
```

## License

MIT License - See [LICENSE](LICENSE) file for details.

This project is released under the MIT License, which is one of the most permissive open-source licenses. You are free to:
- Use commercially
- Modify
- Distribute
- Use privately

The only requirement is to include the original copyright and license notice in any copy of the software/source.

## Credits

- Built with [getevo/evo/v2](https://github.com/getevo/evo)
- Uses [wkhtmltopdf](https://wkhtmltopdf.org/) for HTML to PDF conversion
- HTTP framework: [Fiber](https://gofiber.io/)
