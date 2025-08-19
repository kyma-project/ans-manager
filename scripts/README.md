## Usage

### Get tokens
You need separate tokens for each Service Binding. The tokens are used to authenticate the requests to the ANS Manager API.
- Run the script `refresh-tokens.sh` to get the tokens for both Service Bindings. The script will use the Service Binding files you downloaded in the previous step.
- The script will create two files:  `current-e.token` and `current-n.token` with the tokens for both Service Bindings.

### Get all notification types
Run the script `get-types.sh` to get all available notification types.
The script will use the `current-n.token` file to authenticate the request. The output will be a JSON file with all available notification types.

### Get notification types by ID
Run the script `get-type.sh` to get a specific notification type by its ID.
The script will use the `current-n.token` file to authenticate the request. The output will be a JSON file with the notification type.
Example:
```shell
 ./get-type.sh 0edb0986-719a-42d6-a613-0cb4513bcd4b
```

### Post notification
`post-in-app.sh` script posts a sample notification using the `POC_WebOnlyType` type. The script uses the `current-n.token` file to authenticate the request.
It assumes that this type is already created in the ANS Manager. The script will create a notification with the following content:

Sample
```json
{
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
}
```
This particular notification will be delivered to BTP Cockpit to the user "jaroslaw.pieszka@sap.com" in the subaccount where the ANS Manager Service Instance is created. The notification will be visible in the Notifications section of the BTP Cockpit.

### Post event

Run the script `event.sh` to post a sample event using the notification type `POC_WebOnlyType2` (and assumes this type is created). The script uses the `current-e.token` file to authenticate the request.

```json
{
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
}
```
This script generates event and maps it to the notification type `POC_WebOnlyType2`. The event will be delivered to the BTP Cockpit of the subaccount admins of subaccount `2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad`.
In the Notifications section of the BTP Cockpit you will see the event with the subject "PoC for KEB" and body "Test body". The event will be visible to all subaccount admins of the subaccount.
