token=$(cat ./current-e.token|xargs echo)
curl --request POST \
  --url 'https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com/cf/producer/service/v1/resource-events' \
  --header 'Accept: application/json' \
  --header 'Content-type: application/json' \
  --header 'DataServiceVersion: 2.0' \
  --header "Authorization: Bearer $token" \
  --data '{
  "eventType": "POC-test-KEB-event",
  "body": "Test body",
  "subject": "PoC for KEB",
  "severity": "INFO",
  "visibility": "OWNER_SUBACCOUNT",
  "category": "NOTIFICATION",
  "resource": {
    "globalAccount": "8cd57dc2-edb2-45e0-af8b-7d881006e516",
    "subAccount": "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
    "resourceGroup": "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
    "resourceType": "broker",
    "resourceName": "keb"
  },
  "notificationMapping": {
    "notificationTypeKey": "POC_WebOnlyType2",
    "recipients": {
      "xsuaa":[
	{ 
          "tenantId":"2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
          "level":"SUBACCOUNT",
          "roleNames": ["Subaccount admin"]}
      ]
    }
  }
}'|jq
