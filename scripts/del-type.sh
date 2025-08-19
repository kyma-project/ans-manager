#!/bin/bash
token=$(cat ./current-n.token|xargs echo)
curl -v --request DELETE \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com:443/odatav2/NotificationType.svc/NotificationTypes(guid'\'"$1"\'')' \
  --header 'Content-type: application/json' \
  --header 'Accept: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token"|jq
