#!/bin/bash

# Script to add an Achievement event to the Redis stream.
# Example:
# scripts/add_achievement_event.sh Nandgopal-R https://github.com/IAmRiteshKoushik/devpool/issues/1 "IMPACT"

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <username> <url> <type>"
    echo "Type can be one of: IMPACT, DOC, BUG, TEST, HELP"
    exit 1
fi

USERNAME=$1
URL=$2
TYPE=$3

# Validate type value
case "$TYPE" in
    "IMPACT"|"DOC"|"BUG"|"TEST"|"HELP")
        # Valid type
        ;;
    *)
        echo "Invalid type: $TYPE"
        echo "Type can be one of: IMPACT, DOC, BUG, TEST, HELP"
        exit 1
        ;;
esac

JSON_PAYLOAD=$(printf '{"github_username":"%s","url":"%s","type":"%s"}' "$USERNAME" "$URL" "$TYPE")

STREAM_NAME="automatic-events-stream"

echo "Adding the following payload to stream '$STREAM_NAME':"
echo "$JSON_PAYLOAD"

redis-cli XADD "$STREAM_NAME" '*' data "$JSON_PAYLOAD"
