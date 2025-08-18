#!/bin/bash
token=$(cat ./current-n.token|xargs echo)
curl --request GET \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com:443/odatav2/NotificationType.svc/NotificationTypes(guid'\'"$1"\'')?$expand=Templates,DeliveryChannels'\
  --header 'Content-type: application/json' \
  --header 'Accept: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token"|jq
