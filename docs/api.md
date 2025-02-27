# Social API Documentation

This document provides information about the API endpoints, request formats, and responses.

## Base URL

All endpoints are prefixed with: `/api/v1` 

## Authentication

The API uses JWT tokens for authentication. To authenticate requests, include the token in the Authorization header:

```
Authorization: Bearer <your_token>
```

**Note**: Tokens expire after 72 hours by default (configurable in environment variables)

### Obtaining a Token

You can obtain a token through:
1. Registering a new user (POST /users)
2. Logging in with existing credentials (POST /auth/token)

## API Endpoints

### Health Check

**Endpoint:** `GET /health`

**Description:** Check if the API is running

**Authentication Required:** No

**Response Example:**
```json
{
  "data": {
    "status": "ok",
    "version": "1.0.0",
    "env": "development"
  }
}
```

### User Management

#### Register User

**Endpoint:** `POST /users`

**Description:** Register a new user

**Authentication Required:** No

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response Example:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "created_at": "2025-02-27T10:30:45Z"
    },
    "expires_at": "2025-03-01T10:30:45Z"
  }
}
```

**Validation:**
- `username`: Required, min 3 chars, max 100 chars
- `email`: Required, valid email format, max 255 chars
- `password`: Required, min 8 chars, max 72 chars

#### Login

**Endpoint:** `POST /auth/token`

**Description:** Login with email and password to get a token

**Authentication Required:** No

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response Example:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "created_at": "2025-02-27T10:30:45Z"
    },
    "expires_at": "2025-03-01T10:30:45Z"
  }
}
```

#### Get Current User

**Endpoint:** `GET /users/me`

**Description:** Get the profile of the currently authenticated user

**Authentication Required:** Yes

**Response Example:**
```json
{
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-02-27T10:30:45Z"
  }
}
```

### Post Management

#### Create Post

**Endpoint:** `POST /posts`

**Description:** Create a new post

**Authentication Required:** Yes

**Request Body:**
```json
{
  "title": "My First Post",
  "content": "This is the content of my first post on the platform."
}
```

**Response Example:**
```json
{
  "data": {
    "id": 1,
    "title": "My First Post",
    "content": "This is the content of my first post on the platform.",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "created_at": "2025-02-27T10:30:45Z"
    },
    "created_at": "2025-02-27T14:15:30Z",
    "updated_at": "2025-02-27T14:15:30Z"
  }
}
```

**Validation:**
- `title`: Required, min 3 chars, max 200 chars
- `content`: Required, min 10 chars

#### Get Post by ID

**Endpoint:** `GET /posts/{id}`

**Description:** Get a post by its ID

**Authentication Required:** Yes

**URL Parameters:**
- `id`: Post ID (integer)

**Response Example:**
```json
{
  "data": {
    "id": 1,
    "title": "My First Post",
    "content": "This is the content of my first post on the platform.",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "created_at": "2025-02-27T10:30:45Z"
    },
    "created_at": "2025-02-27T14:15:30Z",
    "updated_at": "2025-02-27T14:15:30Z"
  }
}
```

#### Update Post

**Endpoint:** `PUT /posts/{id}`

**Description:** Update a post (user must be the author)

**Authentication Required:** Yes

**URL Parameters:**
- `id`: Post ID (integer)

**Request Body:**
```json
{
  "title": "Updated Post Title",
  "content": "This is the updated content of my post."
}
```

**Response Example:**
```json
{
  "data": {
    "id": 1,
    "title": "Updated Post Title",
    "content": "This is the updated content of my post.",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "created_at": "2025-02-27T10:30:45Z"
    },
    "created_at": "2025-02-27T14:15:30Z",
    "updated_at": "2025-02-27T14:30:22Z"
  }
}
```

**Validation:**
- `title`: Optional, min 3 chars, max 200 chars if provided
- `content`: Optional, min 10 chars if provided

#### Delete Post

**Endpoint:** `DELETE /posts/{id}`

**Description:** Delete a post (user must be the author)

**Authentication Required:** Yes

**URL Parameters:**
- `id`: Post ID (integer)

**Response:** 
- Status: 204 No Content (No response body)

#### List Posts

**Endpoint:** `GET /posts`

**Description:** List posts with pagination, sorting, and filtering

**Authentication Required:** Yes

**Query Parameters:**
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 20, max: 100)
- `sort`: Sort direction - "asc" or "desc" (default: "desc")
- `sort_by`: Field to sort by (default: "created_at")
- `user_id`: Filter by user ID
- `title`: Filter by title (case-insensitive partial match)
- `content`: Filter by content (case-insensitive partial match)

**Response Example:**
```json
{
  "data": [
    {
      "id": 2,
      "title": "My Second Post",
      "content": "This is another post.",
      "user": {
        "id": 1,
        "username": "johndoe",
        "email": "john@example.com",
        "created_at": "2025-02-27T10:30:45Z"
      },
      "created_at": "2025-02-27T15:20:10Z",
      "updated_at": "2025-02-27T15:20:10Z"
    },
    {
      "id": 1,
      "title": "My First Post",
      "content": "This is the content of my first post.",
      "user": {
        "id": 1,
        "username": "johndoe",
        "email": "john@example.com",
        "created_at": "2025-02-27T10:30:45Z"
      },
      "created_at": "2025-02-27T14:15:30Z",
      "updated_at": "2025-02-27T14:15:30Z"
    }
  ],
  "meta": {
    "current_page": 1,
    "page_size": 20,
    "first_page": 1,
    "last_page": 1,
    "total_records": 2
  }
}
```

## Error Responses

The API returns structured error responses with appropriate HTTP status codes:

### 400 Bad Request

```json
{
  "error": "body must only contain a single JSON object"
}
```

### 401 Unauthorized

```json
{
  "error": "authorization header is required"
}
```

### 404 Not Found

```json
{
  "error": "The requested resource could not be found"
}
```

### 422 Unprocessable Entity (Validation Errors)

```json
{
  "errors": [
    {
      "field": "username",
      "error": "This field is required"
    },
    {
      "field": "password",
      "error": "Must be at least 8 characters long"
    }
  ]
}
```

### 429 Too Many Requests

```json
{
  "error": "Rate limit exceeded. Try again in 5 seconds"
}
```

### 500 Internal Server Error

```json
{
  "error": "The server encountered a problem and could not process your request"
}
```