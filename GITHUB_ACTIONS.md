# GitHub Actions Setup

This document explains the automated Docker image building and publishing workflow for go-pdf.

## Overview

The GitHub Actions workflow automatically builds and publishes Docker images to GitHub Container Registry (ghcr.io) whenever changes are pushed to the repository.

## Workflow File

Location: `.github/workflows/docker-publish.yml`

## What Gets Triggered

The workflow runs automatically on:
- **Push to main branch**: Builds and publishes `latest` tag
- **Pull requests**: Builds the image (doesn't publish)
- **Tags (v*)**: Builds and publishes versioned tags (e.g., `v1.0.0`, `v1.0`, `v1`, `latest`)
- **Manual trigger**: Can be triggered manually from GitHub Actions tab

## What It Does

1. **Checkout Code**: Clones the repository
2. **Setup Docker Buildx**: Enables multi-platform builds
3. **Login to GHCR**: Authenticates with GitHub Container Registry using `GITHUB_TOKEN`
4. **Extract Metadata**: Generates Docker tags and labels
5. **Build & Push**:
   - Builds for multiple platforms: `linux/amd64` and `linux/arm64`
   - Uses layer caching for faster builds
   - Pushes to `ghcr.io/getevo/go-pdf`
6. **Attestation**: Generates build provenance for security

## Docker Image Tags

The workflow automatically creates the following tags:

### On Push to Main Branch
- `ghcr.io/getevo/go-pdf:latest`
- `ghcr.io/getevo/go-pdf:main`

### On Version Tag (e.g., v1.2.3)
- `ghcr.io/getevo/go-pdf:latest`
- `ghcr.io/getevo/go-pdf:1.2.3`
- `ghcr.io/getevo/go-pdf:1.2`
- `ghcr.io/getevo/go-pdf:1`

### On Pull Request
- `ghcr.io/getevo/go-pdf:pr-123` (build only, not pushed)

## Permissions

The workflow has the following permissions:
- `contents: read` - Read repository contents
- `packages: write` - Push to GitHub Container Registry
- `id-token: write` - Generate attestations

## Multi-Platform Support

Images are built for:
- **linux/amd64** - Standard x86_64 servers
- **linux/arm64** - ARM-based servers (Apple Silicon, AWS Graviton, etc.)

Docker automatically pulls the correct image for your platform.

## Using the Published Image

### Pull Latest Image
```bash
docker pull ghcr.io/getevo/go-pdf:latest
```

### Run Container
```bash
docker run -d -p 8080:8080 ghcr.io/getevo/go-pdf:latest
```

### Use Specific Version
```bash
docker pull ghcr.io/getevo/go-pdf:1.0.0
docker run -d -p 8080:8080 ghcr.io/getevo/go-pdf:1.0.0
```

### Docker Compose
```yaml
version: '3.8'
services:
  go-pdf:
    image: ghcr.io/getevo/go-pdf:latest
    ports:
      - "8080:8080"
    restart: unless-stopped
```

## Image Visibility

By default, images pushed to ghcr.io are **public** if the repository is public. Anyone can pull and use the image without authentication.

To make it public (if needed):
1. Go to https://github.com/orgs/getevo/packages/container/go-pdf/settings
2. Change visibility to "Public"
3. Save changes

## Creating a Release

To publish a versioned release:

1. **Create and push a tag**:
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically**:
   - Builds the Docker image
   - Tags it with `1.0.0`, `1.0`, `1`, and `latest`
   - Pushes to ghcr.io

3. **Users can then pull**:
   ```bash
   docker pull ghcr.io/getevo/go-pdf:1.0.0
   ```

## Build Cache

The workflow uses GitHub Actions cache to speed up builds:
- First build: ~2-3 minutes
- Subsequent builds: ~30-60 seconds (with cache)

Cache is shared across workflow runs and branches.

## Monitoring Builds

1. Go to: https://github.com/getevo/go-pdf/actions
2. Click on "Build and Publish Docker Image"
3. View logs and build status

## Troubleshooting

### Build Fails
- Check the Actions tab for error logs
- Verify Dockerfile syntax
- Ensure all dependencies are available

### Image Not Visible
- Check package visibility settings
- Ensure workflow completed successfully
- Wait a few minutes for registry to update

### Permission Denied
- The `GITHUB_TOKEN` is automatically provided
- No manual configuration needed
- Check repository permissions if issues persist

## Security

- **Build Provenance**: Each image includes attestation proving it was built by GitHub Actions
- **SBOM**: Software Bill of Materials can be generated
- **Scanning**: Can integrate with security scanners
- **No Secrets Required**: Uses built-in `GITHUB_TOKEN`

## Workflow Status Badge

Add to README:
```markdown
[![Build and Publish Docker Image](https://github.com/getevo/go-pdf/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/getevo/go-pdf/actions/workflows/docker-publish.yml)
```

Result:
[![Build and Publish Docker Image](https://github.com/getevo/go-pdf/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/getevo/go-pdf/actions/workflows/docker-publish.yml)

## Cost

- GitHub Actions: **Free** for public repositories
- GitHub Container Registry (ghcr.io): **Free** for public images
- Bandwidth: **Free** for pulling public images

## Summary

✅ Fully automated - no manual Docker builds needed
✅ Multi-platform support (amd64 + arm64)
✅ Fast builds with caching
✅ Automatic versioning
✅ Free for public repositories
✅ Secure with attestations

The Docker image is automatically built and published on every push to main or version tag!
