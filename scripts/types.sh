#!/bin/bash
token=$(cat ./current-n.token|xargs echo)
curl --request GET \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/odatav2/NotificationType.svc/NotificationTypes?inlinecount=allpages' \
  --header 'Content-type: application/json' \
  --header 'Accept: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token"|jq
