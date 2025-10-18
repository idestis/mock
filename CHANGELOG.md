# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- GitHub Actions workflow for automated testing
- GitHub Actions workflow for Docker build and push to GHCR
- Multi-architecture Docker image support (linux/amd64, linux/arm64)
- Automated semantic versioning and tagging
- Makefile targets for GHCR operations
- Comprehensive GHCR setup documentation
- CI/CD documentation in README

### Changed
- Fixed invalid SHA tag format in Docker build workflow (changed from `{{branch}}-` to `sha-` prefix)
- Updated Helm values.yaml with GHCR repository comments

## [0.0.1] - 2025-10-18

### Added
- Initial project setup with Golang + Echo framework
- Three hardcoded jobs: user_annonymization, zendesk_import, blueshift_export
- Job status tracking: PENDING, RUNNING, SUCCESS, FAILED
- In-memory job storage with UUID tracking
- Background job execution with configurable duration and failure scenarios
- REST API endpoints:
  - `GET /health` - Health check
  - `POST /api/v1/system/scheduler/trigger` - Trigger job execution
  - `GET /api/v1/system/scheduler/status` - Get job status
- Environment variable configuration for job behavior
- Dockerfile for containerization
- Helm chart for Kubernetes deployment
- Complete documentation (README, setup guides)
- Makefile with build, test, and deployment targets
