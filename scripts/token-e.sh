#!/usr/bin/env bash 
SB_FILE=$1

# standard bash error handling
set -o nounset  # treat unset variables as an error and exit immediately.
set -o errexit  # exit immediately when a command fails.
set -E          # needs to be set if we want the ERR trap
set -o pipefail # prevents errors in a pipeline from being masked
CLIENT_ID=$(jq <$SB_FILE -r '.client_id')
CLIENT_SECRET=$(jq <$SB_FILE -r '.client_secret')
ENCODED_CREDENTIALS=$(echo -n "${CLIENT_ID}:${CLIENT_SECRET}" | base64)
TOKEN_URL=$(jq <$SB_FILE -r '.oauth_url')
TOKEN=$(curl -X POST ${TOKEN_URL} -H "Content-Type: application/x-www-form-urlencoded" -H "Authorization: Basic $ENCODED_CREDENTIALS")
echo ${TOKEN}| jq '.access_token' >current-e.token
