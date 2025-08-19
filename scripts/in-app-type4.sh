token=$(cat ./current-n.token|xargs echo)
curl --request POST \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/odatav2/NotificationType.svc/NotificationTypes' \
  --header 'Accept: application/json' \
  --header 'Content-type: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token"\
  --data '{
  "NotificationTypeId": "a8195025-937e-42cf-b191-00dbc4012a04",
  "NotificationTypeKey": "POC_WebOnlyType4",
  "NotificationTypeVersion": "1",
  "IsGroupable": true,
  "UserPreferencesVisibility": "ON_NOTIFICATION",
  "Templates": [
    {
      "Language": "EN",
      "TemplatePublic": "{subject}: {body}",
      "TemplateSensitive": "{subject}",
      "TemplateGrouped": "KEB",
      "Description": "Kyma Environment Broker started step",
      "TemplateLanguage": "DEFAULT",
      "Subtitle": "KEB started step"
    }
  ],
  "DeliveryChannels": [
    {
      "Type": "WEB",
      "Enabled": true,
      "DefaultPreference": true,
      "EditablePreference": true
    }
  ]
}'|jq
