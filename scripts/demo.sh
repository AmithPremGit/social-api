#!/bin/bash
set -e

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

echo "Social API Demo"
echo "==============="
echo "Ensure API is running with 'docker-compose up -d' before running this script."

echo -e "\n1. Health Check"
curl -s "$BASE_URL/health" | jq

echo -e "\n2. Register a New User"
# Using random username and email to avoid duplicate registration errors
TIMESTAMP=$(date +%s)
RANDOM_USERNAME="demo$TIMESTAMP"
RANDOM_EMAIL="demo$TIMESTAMP@example.com"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$RANDOM_USERNAME\",\"email\":\"$RANDOM_EMAIL\",\"password\":\"password123\"}")
echo $REGISTER_RESPONSE | jq
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.data.token')
echo $REGISTER_RESPONSE | jq
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.data.token')

# Basic error handling
if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo "Error: Failed to get authentication token."
  exit 1
fi

echo -e "\n3. Get Current User Profile"
curl -s -X GET "$BASE_URL/users/me" \
  -H "Authorization: Bearer $TOKEN" | jq

echo -e "\n4. Create a New Post"
POST_RESPONSE=$(curl -s -X POST "$BASE_URL/posts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"My First Post","content":"This is the content of my first post on the platform. It needs to be at least 10 characters long."}')
echo $POST_RESPONSE | jq
POST_ID=$(echo $POST_RESPONSE | jq -r '.data.id')

echo -e "\n5. Get Post by ID"
curl -s -X GET "$BASE_URL/posts/$POST_ID" \
  -H "Authorization: Bearer $TOKEN" | jq

echo -e "\n6. Update Post"
curl -s -X PUT "$BASE_URL/posts/$POST_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Updated Post Title","content":"This post has been updated with new content that is longer than 10 characters."}' | jq

echo -e "\n7. List Posts with Pagination"
curl -s -X GET "$BASE_URL/posts?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq

echo -e "\n8. Delete Post"
curl -s -X DELETE "$BASE_URL/posts/$POST_ID" \
  -H "Authorization: Bearer $TOKEN" -v