token=$(cat ./current-n.token|xargs echo)
curl --request POST \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/odatav2/NotificationType.svc/NotificationTypes' \
  --header 'Accept: application/json' \
  --header 'Content-type: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token"\
  --data '{
  "NotificationTypeId": "8beb8936-eae4-4503-b01c-79d5850f3203",
  "NotificationTypeKey": "POC_WebOnlyType2",
  "NotificationTypeVersion": "1",
  "IsGroupable": true,
  "UserPreferencesVisibility": "ON_NOTIFICATION",
  "Templates": [
    {
      "Language": "EN",
      "TemplatePublic": "Test event for subaccount {resource.subAccount}",
      "TemplateSensitive": "Test sensitive event for subaccount {resource.subAccount}",
      "TemplateGrouped": "Provisioning for {resource.resourceName} in subaccount {resource.subAccount} started",
      "Description": "Kyma Environment Broker started provisioning",
      "TemplateLanguage": "DEFAULT",
      "Subtitle": "Kyma Environment Broker Notification"
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
