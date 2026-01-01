#!/bin/bash

# Script to add a BountyAction event to the Redis stream.
# scripts/add_bounty_event.sh Nandgopal-R 100 https://github.com/IAmRiteshKoushik/devpool/issues/1 "BOUNTY"

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <username> <amount> <url> <action>"
    exit 1
fi

USERNAME=$1
AMOUNT=$2
URL=$3
ACTION=$4

JSON_PAYLOAD=$(printf '{"github_username":"%s","amount":%d,"url":"%s","action":"%s"}' "$USERNAME" "$AMOUNT" "$URL" "$ACTION")

STREAM_NAME="bounty-stream"

echo "Adding the following payload to stream '$STREAM_NAME':"
echo "$JSON_PAYLOAD"

redis-cli XADD "$STREAM_NAME" '*' data "$JSON_PAYLOAD"
