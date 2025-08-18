token=$(cat ./current-n.token|xargs echo)
curl -v  --request POST \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/odatav2/Notification.svc/Notifications' \
  --header 'Accept: application/json' \
  --header 'Content-type: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token"\
  --data '{
  "NotificationTypeKey": "POC_WebOnlyType",
  "Recipients": [
    {
      "RecipientId": "jaroslaw.pieszka@sap.com",
      "IasHost": "accounts.sap.com"
    }
  ],
  "Properties": [
    {
      "Key": "shoot",
      "Value": "c0123456"
    }
  ]
}'
