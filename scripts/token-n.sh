#!/usr/bin/env bash
SB_FILE=$1

# standard bash error handling
set -o nounset  # treat unset variables as an error and exit immediately.
set -o errexit  # exit immediately when a command fails.
set -E          # needs to be set if we want the ERR trap
set -o pipefail # prevents errors in a pipeline from being masked
CLIENT_ID=$(jq <$SB_FILE -r '.uaa.clientid')
CLIENT_SECRET=$(jq <$SB_FILE -r '.uaa.clientsecret')
ENCODED_CREDENTIALS=$(echo -n "${CLIENT_ID}:${CLIENT_SECRET}" | base64)
TOKEN_URL=$(jq <$SB_FILE -r '.uaa.url')
TOKEN=$(curl -X POST ${TOKEN_URL}/oauth/token -H "Content-Type: application/x-www-form-urlencoded" -H "Authorization: Basic $ENCODED_CREDENTIALS" --data-urlencode "grant_type=client_credentials")
echo ${TOKEN}| jq '.access_token' >current-n.token
