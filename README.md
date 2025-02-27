# Social API

A RESTful API for social media platforms built with Go, featuring JWT authentication, PostgreSQL database, and Redis caching.

This application provides a complete backend solution for social networking services where users can register accounts, create and manage posts, and interact with content. Built with Go and the Chi router, the API implements best practices for authentication, data validation, rate limiting, and performance optimization.

The codebase follows a clean architecture pattern that separates concerns between different layers, making it modular and extensible for future feature additions.

## Features

- **User Management** - Registration, authentication and profiles
- **Content Management** - CRUD operations for posts
- **JWT Authentication** - Secure token-based auth (72hr default expiry)
- **Redis Caching** - Optional performance enhancement
- **Rate Limiting** - Protection against API abuse
- **Input Validation** - Request validation and sanitization
- **Structured Responses** - Consistent JSON formatting
- **Pagination & Filtering** - Support for large datasets
- **Error Handling** - Clear, structured error responses

## Technology Stack

- **Backend**: Go (Golang)
- **Router**: Chi
- **Database**: PostgreSQL
- **Caching**: Redis
- **Authentication**: JWT with bcrypt password hashing
- **Containerization**: Docker & Docker Compose
- **Testing**: Postman collection included

## Architecture

The application follows a clean architecture with:

- **Handler Layer**: HTTP request handling
- **Service Layer**: Business logic
- **Store Layer**: Data access interfaces
- **Repository Layer**: Database implementations

The system is containerized with Docker and includes:
- Go API server
- PostgreSQL database
- Redis cache

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for development)
- Postman (for testing)

### Running the API

1. Clone the repository:
   ```bash
   git clone https://github.com/AmithPremGit/social-api.git
   cd social-api
   ```

2. Start the Docker containers:
   ```bash
   docker-compose up -d
   ```

3. The API will be available at http://localhost:8080

### Running the Demo

To quickly explore the API functionality:

```bash
chmod +x scripts/demo.sh
./scripts/demo.sh
```

## API Documentation

### Authentication Endpoints

| Method | Endpoint         | Description     | Auth Required |
|--------|------------------|-----------------|---------------|
| POST   | /api/v1/users    | Register user   | No            |
| POST   | /api/v1/auth/token | Login         | No            |
| GET    | /api/v1/users/me | Get current user| Yes           |

### Post Endpoints

| Method | Endpoint         | Description     | Auth Required |
|--------|------------------|-----------------|---------------|
| GET    | /api/v1/posts    | List posts      | Yes           |
| POST   | /api/v1/posts    | Create post     | Yes           |
| GET    | /api/v1/posts/{id} | Get post      | Yes           |
| PUT    | /api/v1/posts/{id} | Update post   | Yes           |
| DELETE | /api/v1/posts/{id} | Delete post   | Yes           |

### System Endpoints

| Method | Endpoint         | Description     | Auth Required |
|--------|------------------|-----------------|---------------|
| GET    | /api/v1/health   | Health check    | No            |

## Testing the API

### Using Postman

A Postman collection is included in the `/docs` directory:

1. Import both the collection and environment files
2. Select "Social API Environment" from the dropdown
3. The collection includes all endpoints with pre-configured requests

### Using the Demo Script

```bash
chmod +x scripts/demo.sh
./scripts/demo.sh
```

This script demonstrates the API functionality with cURL requests.

### Using cURL Commands

For individual endpoint testing, see `docs/curl-examples.md` for comprehensive cURL examples.

## Project Structure

```
social-api/
├── cmd/                  # Application entry points
├── internal/             # Private application code
│   ├── auth/             # Authentication logic
│   ├── cache/            # Cache implementation
│   ├── config/           # Configuration
│   ├── db/               # Database connection
│   ├── handler/          # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── model/            # API models
│   └── store/            # Data store interfaces and implementations
├── migrations/           # Database migrations
├── docs/                 # Documentation and Postman collection
├── scripts/              # Utility scripts

├── docker-compose.yaml   # Docker Compose configuration
└── Dockerfile            # Dockerfile for the API
```

## Security Features

- Secure password hashing with bcrypt
- JWT tokens with configurable expiration
- Rate limiting for API protection
- Input validation and sanitization
- Request context timeouts

## Performance Optimizations

- Database connection pooling
- Redis caching for frequently accessed data
- SQL query optimization with proper indexing
- Pagination for large result sets
- Structured logging

## Future Improvements

- Image upload support
- Comment functionality
- Email verification
- Social interactions (likes, follows)
- Full-text search
- Extended test coverage
- WebSocket support for real-time notifications
