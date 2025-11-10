# Stream Microservices

A video streaming platform built with microservices architecture using Go, PostgreSQL, MongoDB, and RabbitMQ.

## Services

- **History Service** (`:8081`) - Tracks user viewing history using PostgreSQL
- **Recommendations Service** (`:8082`) - Provides video recommendations using MongoDB
- **Video Storage Service** (`:8080`) - Manages video metadata using MongoDB
- **Video Streaming Service** (`:8083`) - Handles video streaming using MongoDB

## Tech Stack

- **Backend**: Go
- **Databases**: PostgreSQL, MongoDB
- **Message Queue**: RabbitMQ
- **Containerization**: Docker

## Quick Start

1. Navigate to each service directory
2. Copy `.env.example` to `.env` and configure
3. Run with Docker Compose:
   ```bash
   docker-compose up -d
   ```

## Health Checks

All services expose `/health` endpoints for monitoring.
