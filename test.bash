#!/bin/bash

IP="49.43.41.152"
URL="http://localhost:8080/api"
REQUESTS=200

for i in $(seq 1 $REQUESTS); do
  response=$(curl -s -w "%{http_code}" -H "X-Forwarded-For: $IP" $URL)
  http_code=$(echo "$response" | tail -n 1)
  http_body=$(echo "$response" | head -n -1)

  if [[ "$http_code" =~ ^(2..) ]]; then # Check for 2xx success codes
    echo "Request $i successful (HTTP $http_code) - Response: $http_body"
  else
    echo "Request $i failed (HTTP $http_code) - Response: $http_body"
  fi
done

echo "All $REQUESTS requests completed."