# Using the Social API Postman Collection

## Quick Start

1. Import both files in Postman:
   - `social-api.postman_collection.json`
   - `social-api.postman_environment.json`

2. Select "Social API Environment" from the environment dropdown

3. The collection is organized into folders:
   - Health - API health check
   - Users - Registration, login, profile
   - Posts - Create, read, update, delete posts

4. Authentication:
   - Use the "Login" request first
   - The token should be automatically stored in the environment
   - If authentication fails, manually copy the token from the login response and add it to the environment variables

## Base URL

The default base URL is `http://localhost:8080`. Change this in the environment variables if your API is hosted elsewhere.