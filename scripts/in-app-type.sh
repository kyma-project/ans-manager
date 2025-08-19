token=$(cat ./current-n.token|xargs echo)
curl --request POST \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/odatav2/NotificationType.svc/NotificationTypes' \
  --header 'Accept: application/json' \
  --header 'Content-type: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header 'X-CSRF-Token: B624D9994A3299400DB5DA07CE933459' \
  --header "Authorization: Bearer $token"\
  --data '{
  "NotificationTypeId": "0edb0986-719a-42d6-a613-0cb4513bcd4b",
  "NotificationTypeKey": "POC_WebOnlyType",
  "NotificationTypeVersion": "1",
  "Templates": [
    {
      "Language": "EN",
      "Description": "This is a description for sample notification",
      "Subtitle": "Sample Notification - subtitle",
      "TemplatePublic": "PoC Template Public",
      "TemplateSensitive": "PoC Template Public",
      "TemplateGrouped": "PoC Template Grouped",
      "EmailSubject": "PoC Web Notification",
      "EmailHtml": "<html>This is a web notification. It does not need your attention</html>",
      "EmailText": "This is a web notification. It does not need your attention",
      "EmailSenderName": "Gophers"
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
