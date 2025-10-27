# Quick Start Guide

This guide will get you up and running with go-pdf in less than 5 minutes.

## 1. Prerequisites

- Docker Desktop installed and running

## 2. Start the Service

```bash
docker-compose up -d
```

This will:
- Build the Docker image (~400MB)
- Start the service on port 8080
- Run in the background

## 3. Test the Service

Run the included test script:

```bash
./test-api.sh
```

Or test manually with curl:

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: text/html" \
  -d '<html><body><h1>Hello PDF!</h1></body></html>' \
  -o output.pdf
```

## 4. View Logs

```bash
docker-compose logs -f
```

Press `Ctrl+C` to stop viewing logs.

## 5. Stop the Service

```bash
docker-compose down
```

## API Quick Reference

**Endpoint**: `POST /api/v1/generate`

**Request**:
- **Content-Type**: text/html
- **Body**: Raw HTML content (string)

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: text/html" \
  -d '<html>...</html>' \
  -o output.pdf
```

**Response**: PDF file (binary)

## Common Use Cases

### 1. Generate Invoice

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: text/html" \
  -d '<!DOCTYPE html><html><head><style>body{font-family:Arial;}</style></head><body><h1>Invoice #001</h1><p>Amount: $100</p></body></html>' \
  -o invoice.pdf
```

### 2. Generate Report with Styling

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: text/html" \
  -d '<!DOCTYPE html><html><head><style>.header{background:#4CAF50;color:white;padding:20px;}</style></head><body><div class="header"><h1>Monthly Report</h1></div></body></html>' \
  -o report.pdf
```

## Troubleshooting

### Service won't start
```bash
docker-compose logs
```

### Port 8080 is in use
Edit `docker-compose.yml` and change the port:
```yaml
ports:
  - "9090:8080"  # Use port 9090 instead
```

### PDF not generating
- Check that HTML is valid
- Ensure HTML is properly escaped in JSON
- View service logs: `docker logs go-pdf`

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check the [config.yml](config.yml) for configuration options
- Explore the source code in `apps/pdf/`

## Need Help?

1. Check the logs: `docker logs go-pdf`
2. Verify service is running: `docker ps`
3. Test with simple HTML first
4. Ensure JSON is properly formatted
