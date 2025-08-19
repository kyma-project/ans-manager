## Usage

### Get tokens
You need separate tokens for each Service Binding. The tokens are used to authenticate the requests to the ANS Manager API.
- Run the script `refresh-tokens.sh` to get the tokens for both Service Bindings. The script will use the Service Binding files you downloaded in the previous step.
- The script will create two files:  `current-e.token` and `current-n.token` with the tokens for both Service Bindings.

### Create a notification type
There are the following scripts creating notification types:
- in-app-type.sh creates "POC_WebOnlyType"
- in-app-type2.sh creates "POC_WebOnlyType2"
- in-app-type3.sh creates "POC_WebOnlyType3"
- in-app-type4.sh creates "POC_WebOnlyType4"

The notification type used in end-to-end tests is "POC_WebOnlyType2" is defined as follows:
```json
{
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
}
```
In the template, there are placeholders for the resource data of the event. The placeholders are replaced with the actual values when the notification is sent. If the notification is sent without an event,
the placeholders are replaced Properties provided in the notification itself. 
If you send an event with notifications, placeholders are replaced with the values from the event as follows:
- `{resource.<name>}` is replaced with the `resource.<name>` value from the event e.g. 
  - `{resource.subAccount}` is replaced with the subaccount ID of the event (`ResourceEvent.Resource.SubAccount`)
- `{resource.tags.<name>}` is replaced with the `<name>` tag value from the resource (`ResourceEvent.Resource.Tags`)
- `{tags.<name>}` is replaced with the `<name>` tag value from the event (`ResourceEvent.Tags`)
- `{body}` is replaced with the body of the event (`ResourceEvent.Body`)
- `{subject}` is replaced with the subject of the event (`ResourceEvent.Subject`)
- `{eventType}` is replaced with the event type of the event (`ResourceEvent.EventType`)

### Patch a notification type
Run the script `patch-in-app.sh` to patch the template of notification type identified by given ID.

### Delete a notification type
Run the script `del-type.sh` to delete a notification type by its ID.


### Get all notification types
Run the script `get-types.sh` to get all available notification types.
The script will use the `current-n.token` file to authenticate the request. The output will be a JSON file with all available notification types.

### Get a notification type by ID
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
