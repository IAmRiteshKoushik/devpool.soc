#!/bin/bash

# Example to add an IssueAction event to the Redis stream.
# Usage: ./add_issue_event.sh Nandgopal-R https://github.com/IAmRiteshKoushik/jsonite/issues/3 true 

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <username> <url> <true|false>"
    exit 1
fi

USERNAME=$1
URL=$2
CLAIM=$3

# Validate claim value
if [ "$CLAIM" != "true" ] && [ "$CLAIM" != "false" ]; then
    echo "Claim must be 'true' or 'false'"
    exit 1
fi

JSON_PAYLOAD=$(printf '{"github_username":"%s","url":"%s","claimed":%s}' "$USERNAME" "$URL" "$CLAIM")

STREAM_NAME="issue-stream"

echo "Adding the following payload to stream '$STREAM_NAME':"
echo "$JSON_PAYLOAD"

redis-cli XADD "$STREAM_NAME" '*' data "$JSON_PAYLOAD"
