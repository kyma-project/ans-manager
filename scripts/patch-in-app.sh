token=$(cat ./current-n.token|xargs echo)
curl --request PATCH \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/odatav2/NotificationType.svc/NotificationTypes(guid'\'"$1"\'')' \
  --header 'Accept: application/json' \
  --header 'Content-type: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token"\
  --data '{
  "Templates": [
    {
      "Language": "EN",
      "TemplatePublic": "Provisioning started {shoot}",
      "TemplateSensitive": "Provisioning started",
      "TemplateGrouped": "Provisioning",
      "Description": "Kyma Environment Broker started provisioning",
      "Subtitle": "Kyma Environment Broker Notification"
    }
  ]
}'|jq
