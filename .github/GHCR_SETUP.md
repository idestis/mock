# GitHub Container Registry (GHCR) Setup Guide

This guide explains how to set up and use GitHub Container Registry for this project.

## Prerequisites

- GitHub account with access to this repository
- Docker installed locally
- GitHub Personal Access Token (PAT) with `write:packages` scope

## Creating a GitHub Personal Access Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Give it a descriptive name (e.g., "GHCR Access")
4. Select scopes:
   - `write:packages` (to upload packages)
   - `read:packages` (to download packages)
   - `delete:packages` (optional, to delete packages)
5. Click "Generate token"
6. **Copy the token** - you won't be able to see it again!

## Local Setup

### 1. Export Your GitHub Token

```bash
export GITHUB_TOKEN=ghp_your_token_here
export GITHUB_USER=your-github-username
```

Add these to your `~/.bashrc` or `~/.zshrc` for persistence:

```bash
echo 'export GITHUB_TOKEN=ghp_your_token_here' >> ~/.zshrc
echo 'export GITHUB_USER=your-github-username' >> ~/.zshrc
```

### 2. Login to GHCR

```bash
make ghcr-login
```

Or manually:

```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_USER --password-stdin
```

### 3. Build and Push Images

Using Makefile:

```bash
# Build multi-arch image and push
make ghcr-build-push VERSION=v1.0.0

# Or build first, then push
make docker-build
make ghcr-push VERSION=v1.0.0
```

Manual commands:

```bash
# Build multi-architecture image
docker buildx build --platform linux/amd64,linux/arm64 \
  -t ghcr.io/your-username/mock:v1.0.0 \
  -t ghcr.io/your-username/mock:latest \
  --push .
```

## GitHub Actions Setup

The workflows are already configured and will run automatically when you:

1. **Push to main/develop branches** - Builds and pushes with branch name tag
2. **Create a release tag** (e.g., `v1.0.0`) - Builds and pushes with version tags

### Enabling GitHub Actions

1. Go to your repository on GitHub
2. Click on the "Actions" tab
3. If prompted, enable GitHub Actions
4. Workflows will run automatically on push/tag events

### Making Packages Public

By default, packages are private. To make them public:

1. Go to your GitHub profile
2. Click on "Packages"
3. Find the `mock` package
4. Click "Package settings"
5. Scroll down to "Danger Zone"
6. Click "Change visibility" → "Public"

## Using Images from GHCR

### Pull and Run

```bash
# Pull the image
docker pull ghcr.io/your-username/mock:latest

# Run the container
docker run -p 8080:8080 ghcr.io/your-username/mock:latest
```

### Deploy to Kubernetes

Create image pull secret (if package is private):

```bash
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=$GITHUB_USER \
  --docker-password=$GITHUB_TOKEN \
  --docker-email=your-email@example.com
```

Deploy with Helm:

```bash
helm install mock ./helm/mock \
  --set image.repository=ghcr.io/your-username/mock \
  --set image.tag=v1.0.0 \
  --set imagePullSecrets[0].name=ghcr-secret
```

## Versioning Strategy

### Semantic Versioning

Follow semantic versioning (semver) for releases:

- **Major version** (`v1.0.0` → `v2.0.0`): Breaking changes
- **Minor version** (`v1.0.0` → `v1.1.0`): New features, backwards compatible
- **Patch version** (`v1.0.0` → `v1.0.1`): Bug fixes, backwards compatible

### Creating a Release

```bash
# Create a new version tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# Push the tag
git push origin v1.0.0
```

GitHub Actions will automatically create these tags:
- `ghcr.io/your-username/mock:v1.0.0`
- `ghcr.io/your-username/mock:v1.0`
- `ghcr.io/your-username/mock:v1`
- `ghcr.io/your-username/mock:latest`

### Tag-Only Builds

This project only builds Docker images on version tags. No automatic builds on branch pushes.
This keeps your container registry clean and ensures only versioned releases are published.

## Troubleshooting

### Authentication Failed

```
Error: denied: permission_denied
```

**Solution**: Check that your token has `write:packages` scope and you're logged in.

### Image Not Found

```
Error: manifest unknown
```

**Solution**: Verify the image name and tag. Check your GitHub packages page.

### Multi-arch Build Issues

```
Error: multiple platforms feature is currently not supported
```

**Solution**: Create a buildx builder:

```bash
docker buildx create --use
docker buildx inspect --bootstrap
```

### Rate Limiting

GitHub Container Registry has rate limits for unauthenticated pulls. Always authenticate for better limits.

## Best Practices

1. **Tag releases properly** - Use semantic versioning
2. **Keep latest tag updated** - Push to main to update `latest`
3. **Document breaking changes** - In release notes when bumping major version
4. **Clean up old images** - Delete unused tags to save space
5. **Use specific versions in production** - Don't use `latest` in prod deployments

## Useful Commands

```bash
# List all make targets
make help

# View workflow runs
gh run list

# View workflow logs
gh run view

# Delete a package version
gh api -X DELETE /user/packages/container/mock/versions/VERSION_ID
```

## Resources

- [GitHub Container Registry Documentation](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Buildx Documentation](https://docs.docker.com/buildx/working-with-buildx/)
