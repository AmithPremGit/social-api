# Social API Architecture

This document explains the architecture and design patterns used in the Social API.

## Overview

The Social API follows clean architecture principles with separation of concerns between different layers of the application. The design makes the code modular, testable, and maintainable.

## System Components

### Docker Infrastructure

The application is containerized using Docker, with three main services:

1. **API Server**: A Go service that handles HTTP requests
2. **PostgreSQL**: Database for persistent storage
3. **Redis**: In-memory cache for improved performance

## Application Layers

The application code is organized into the following layers:

### 1. Handler Layer (internal/handler)

- Responsible for HTTP request handling
- Parses and validates incoming requests
- Calls the appropriate store methods
- Formats and returns responses
- Handles error responses

### 2. Middleware Layer (internal/middleware)

- Intercepts and processes HTTP requests before they reach handlers
- Implements cross-cutting concerns like:
  - Authentication
  - Logging
  - Rate limiting
  - Recovery from panics

### 3. Model Layer (internal/model)

- Defines data structures for API requests and responses
- Provides validation rules
- Implements JSON marshaling/unmarshaling
- Handles pagination and filtering

### 4. Store Layer (internal/store)

- Defines interfaces for data access
- Contains domain models
- Abstracts database operations
- Enables dependency injection and testing

### 5. Store Implementation Layer (internal/store/postgres)

- Implements the store interfaces using PostgreSQL
- Translates domain models to database queries
- Handles database-specific error handling

### 6. Authentication Layer (internal/auth)

- Implements JWT token generation and validation
- Provides middleware for securing routes
- Manages user sessions

### 7. Cache Layer (internal/cache)

- Abstracts caching operations
- Provides Redis implementation
- Improves performance for frequently accessed data

### 8. Configuration Layer (internal/config)

- Loads application configuration from environment variables
- Provides type-safe access to configuration values

### 9. Database Layer (internal/db)

- Manages database connections
- Handles connection pooling
- Provides transaction support

## Data Flow

1. HTTP request comes in
2. Middleware processes the request (logging, rate limiting, etc.)
3. If a protected route, authentication middleware verifies the JWT token
4. Handler processes the request and calls the appropriate store methods
5. Store implementation interacts with the database
6. Database returns data to the store
7. Store returns data to the handler
8. Handler formats the response and returns it to the client

## Design Patterns

### 1. Repository Pattern

The store interfaces and implementations follow the repository pattern, abstracting data access and allowing for different implementations.

### 2. Dependency Injection

Components receive their dependencies through constructors, making them testable and loosely coupled.

### 3. Middleware Pattern

HTTP middleware is used to implement cross-cutting concerns like authentication, logging, and rate limiting.

### 4. Factory Pattern

Factory functions like `NewUserStore` create and configure components with their dependencies.

### 5. Strategy Pattern

Interfaces like `Cache` allow for different implementation strategies (Redis, in-memory, etc.).

## Concurrency Model

- Go's goroutines and channels are used for concurrent processing
- Context is used for cancellation and timeouts
- Database connection pooling handles concurrent database access

## Error Handling

- Structured error responses with consistent format
- Custom error types for common scenarios
- Error handling at appropriate layers
- Consistent error logging

## Security Considerations

- Password hashing with bcrypt (default cost: 10)
- JWT token authentication with configurable expiration
- Input validation to prevent injection attacks
- Rate limiting (default: 20 requests per 5-second window)
- Parameterized queries to prevent SQL injection

## Performance Optimizations

- Redis caching for frequently accessed data
- Connection pooling (max 25 connections by default)
- Request timeouts (60s global timeout)
- Database indexing for query performance
- Pagination for large datasets

## Implementation Notes

- Redis cache invalidation occurs on data updates
- Error handling patterns vary slightly between packages
- JWT validation follows industry best practices
- Rate limiting uses a fixed-window algorithm
- Database connections use default timeout handling