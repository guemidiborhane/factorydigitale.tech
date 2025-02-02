#!/bin/sh

METHOD="${2:-GET}"
BASE_URL="${3:-localhost:3000}"
URL="$1"

# Function to extract CSRF token from cookies
get_csrf_token() {
	cookies="$1"
	echo "$cookies" | grep csrf_ | awk -F'=' '{print $2}' | cut -d ';' -f1
}

# Perform an initial request to obtain cookies
initial_response="$(curl -X GET -si "$BASE_URL/api")"
cookies="$(echo "$initial_response" | grep Set-Cookie | sed 's/Set-Cookie: //g')"

# Extract the CSRF token
csrf_token="$(get_csrf_token "$cookies")"

# Create the authentication payload
username="admin"
password="password"
payload='{"username": "'$username'", "password": "'$password'"}'

# Perform authentication and obtain authentication cookies
auth_response="$(curl -X POST -si -H 'Content-Type: application/json' -H "X-CSRF-Token: $csrf_token" -H "Cookie: $(echo "$cookies" | tr '\n' ';')" -d "$payload" "$BASE_URL/api/auth")"
auth_cookies="$(echo "$auth_response" | grep Set-Cookie | sed 's/Set-Cookie: //g')"

# Extract the CSRF token from authentication cookies
csrf_token="$(get_csrf_token "$auth_cookies")"

# Perform the main request with the provided method
curl -X "$METHOD" -H "Cookie: $(echo "$auth_cookies" | tr '\n' ';')" -s -D /dev/stderr "$BASE_URL$URL"

# Output the method, base URL, and URL to stderr
echo "$METHOD $BASE_URL$URL" >/dev/stderr

# Perform a cleanup request to delete the authentication session
curl -X DELETE -H "Cookie: $(echo "$auth_cookies" | tr '\n' ';')" -H "X-CSRF-Token: $csrf_token" -s "$BASE_URL/api/auth" >/dev/null
