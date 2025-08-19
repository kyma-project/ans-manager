token=$(cat ./current-n.token|xargs echo)
curl -v --cookie cookie-jar.txt --request PATCH \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/odatav2/NotificationType.svc/NotificationTypes(guid'\''0edb0986-719a-42d6-a613-0cb4513bcd4b'\'')' \
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
