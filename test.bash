#!/bin/bash

IP="49.43.41.152"
URL="http://localhost:8080/api"
REQUESTS=200

for i in $(seq 1 $REQUESTS); do
  response=$(curl -s -o /dev/stderr -w "%{http_code}" -H "X-Forwarded-For: $IP" $URL)
  http_code=$(echo $response | tail -n 1) #get the http status code from the curl output.

  if [ "$http_code" -eq 429 ]; then #429 is http.StatusTooManyRequests
    echo "Request $i failed: Too Many Requests (429)"
  elif [ "$http_code" -eq 200 ] || [ "$http_code" -eq 201 ]; then #add other successful status codes as needed.
    echo "Request $i successful (HTTP $http_code)"
  else
    echo "Request $i failed: HTTP $http_code"
  fi
done

echo "All $REQUESTS requests completed."