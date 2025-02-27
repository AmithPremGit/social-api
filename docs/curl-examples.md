# Social API cURL Examples

The following examples demonstrate how to interact with the Social API using cURL commands.

## Health Check

```bash
curl -X GET http://localhost:8080/api/v1/health | jq
```

## User Management

### Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }' | jq
```

### Login (Get Authentication Token)

```bash
curl -X POST http://localhost:8080/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' | jq
```

### Get Current User Profile

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
```

## Post Management

### Create a Post

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post on the platform."
  }' | jq
```

### Get a Post by ID

```bash
curl -X GET http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
```

### Update a Post

```bash
curl -X PUT http://localhost:8080/api/v1/posts/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "Updated Post Title",
    "content": "This is the updated content of my post."
  }' | jq
```

### List Posts (with Pagination and Filtering)

```bash
curl -X GET "http://localhost:8080/api/v1/posts?page=1&page_size=10&sort=desc&sort_by=created_at" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
```

### Filter Posts by User ID

```bash
curl -X GET "http://localhost:8080/api/v1/posts?user_id=1" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
```

### Search Posts by Title

```bash
curl -X GET "http://localhost:8080/api/v1/posts?title=first" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq
```

### Delete a Post

```bash
curl -X DELETE http://localhost:8080/api/v1/posts/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Notes

- Replace `YOUR_TOKEN_HERE` with the actual JWT token received during login or registration
- All commands use `localhost:8080` - update this to match your actual API server address
- The `jq` command formats JSON responses for better readability. If you don't have it installed, you can remove `| jq` from the commands