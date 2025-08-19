token=$(cat ./current-n.token|xargs echo)
curl --request POST \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/odatav2/NotificationType.svc/NotificationTypes' \
  --header 'Accept: application/json' \
  --header 'Content-type: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token"\
  --data '{
  "NotificationTypeId": "fe915a5a-9abe-4a70-ad01-bdc5af3b1756",
  "NotificationTypeKey": "POC_WebOnlyType3",
  "NotificationTypeVersion": "1",
  "IsGroupable": true,
  "UserPreferencesVisibility": "ON_NOTIFICATION",
  "Templates": [
    {
      "Language": "EN",
      "TemplatePublic": "{operation}: started {step} for shoot {shoot} in subaccount {resource.subAccount}",
      "TemplateSensitive": "started step in subaccount {resource.subAccount}",
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
