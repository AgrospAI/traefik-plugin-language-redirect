#!/bin/bash

# Function to test a URL with optional headers
test_request() {
  local url="$1"
  local expected_status="$2"
  local expected_location="$3"
  shift 3
  local headers=("$@")

  # Build curl header options
  local header_opts=()
  for h in "${headers[@]}"; do
    header_opts+=("-H" "$h")
  done

  # Make request
  response=$(curl -s -I -o - "${header_opts[@]}" "$url")
  status=$(echo "$response" | head -n 1 | awk '{print $2}')
  location=$(echo "$response" | grep -i '^Location:' | awk '{print $2}' | tr -d '\r\n')

  # Normalize trailing slash
  location="${location%/}"
  expected_location="${expected_location%/}"

  echo "Testing: $url ${headers[*]}"

  if [[ "$status" == "$expected_status" && "$location" == "$expected_location" ]]; then
    echo "✅ Test passed"
  else
    echo "Expected status: $expected_status, got: $status"
    echo "Expected location: '$expected_location', got: '$location'"
    echo "❌ Test failed"
  fi
  echo
}

# Tests
test_request "http://localhost/en" "200" ""
test_request "http://localhost/" "302" "http://localhost/en"
test_request "http://localhost/" "302" "http://localhost/es" "Accept-Language: es"
test_request "http://localhost/" "302" "http://localhost/de" "Cookie: lang=de"