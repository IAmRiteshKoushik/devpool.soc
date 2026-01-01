#!/bin/bash

# Script to add a Solution event to the Redis stream.
#
# Usage: ./add_solution_event.sh <username> <pr_url> <true|false>
#
# Example:
# ./add_solution_event.sh IAmRiteshKoushik https://github.com/IAmRiteshKoushik/devpool/pull/1 true

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <username> <pr_url> <true|false>"
    exit 1
fi

USERNAME=$1
PR_URL=$2
MERGED=$3

# Validate merged value
if [ "$MERGED" != "true" ] && [ "$MERGED" != "false" ]; then
    echo "Merged must be 'true' or 'false'"
    exit 1
fi

JSON_PAYLOAD=$(printf '{"github_username":"%s","pull_request_url":"%s","merged":%s}' "$USERNAME" "$PR_URL" "$MERGED")

STREAM_NAME="solution-merged-stream"

echo "Adding the following payload to stream '$STREAM_NAME':"
echo "$JSON_PAYLOAD"

redis-cli XADD "$STREAM_NAME" '*' data "$JSON_PAYLOAD"
