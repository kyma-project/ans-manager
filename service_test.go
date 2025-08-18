package ans_manager

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/kyma-project/ans-manager/events"
	"github.com/kyma-project/ans-manager/notifications"
	"github.com/stretchr/testify/require"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

func Test_PostEvent(t *testing.T) {
	t.Skip()
	config := EndpointConfig{
		ClientID:               os.Getenv("E_CLIENT_ID"),
		ClientSecret:           os.Getenv("E_CLIENT_SECRET"),
		AuthURL:                "https://jp-notifications-lxe3vgwv.authentication.stagingaws.hanavlab.ondemand.com/oauth/token",
		ServiceURL:             "https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com",
		RateLimitingInterval:   2 * time.Second,
		MaxRequestsPerInterval: 5,
	}
	client := NewEventsClient(context.Background(), config, logger.With("component", "ANS-eventsClient"))
	require.NotNil(t, client)
	recipient := events.NewXsuaaRecipient(events.LevelSubaccount, "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad", []events.RoleName{"Subaccount admin"})
	require.NoError(t, recipient.Validate())
	recipients := events.NewRecipients([]events.XsuaaRecipient{*recipient}, nil)
	require.NoError(t, recipients.Validate())
	notificationMapping := events.NewNotificationMapping("POC_WebOnlyType2", *recipients)
	require.NoError(t, notificationMapping.Validate())
	resource := events.NewResource("broker",
		"keb",
		"2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
		"2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
		events.WithResourceGlobalAccount("8cd57dc2-edb2-45e0-af8b-7d881006e516"))
	require.NoError(t, resource.Validate())
	event, err := events.NewResourceEvent(
		"eventType",
		"body",
		"subject",
		resource,
		events.SeverityInfo,
		events.CategoryNotification,
		events.VisibilityOwnerSubAccount,
		*notificationMapping,
	)
	require.NotNil(t, event)
	require.NoError(t, event.Validate())
	eventAsJSON, err := json.Marshal(event)
	require.JSONEq(t, `{
  "body": "Test body",
  "subject": "PoC for KEB",
  "eventType": "POC-test-KEB-event",
  "resource": {
    "resourceType": "broker",
    "resourceName": "keb",
    "subAccount": "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
    "globalAccount": "8cd57dc2-edb2-45e0-af8b-7d881006e516",
    "resourceGroup": "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad"
  },
  "severity": "INFO",
  "category": "NOTIFICATION",
  "visibility": "OWNER_SUBACCOUNT",
  "notificationMapping": {
    "notificationTypeKey": "POC_WebOnlyType2",
    "recipients": {
      "xsuaa": [
        {
          "level": "SUBACCOUNT",
          "tenantId": "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
          "roleNames": ["Subaccount admin"]
        }
      ]
    }
  }
}`, string(eventAsJSON))
	require.JSONEq(t, `{
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
      "xsuaa": [
        {
          "tenantId": "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
          "level": "SUBACCOUNT",
          "roleNames": ["Subaccount admin"]
        }
      ]
    }
  }
}`, string(eventAsJSON))
	err = client.postEvent(*event)
	require.NoError(t, err)
}

func Test_PostNotifications(t *testing.T) {
	t.Skip()
	config := EndpointConfig{
		ClientID:               os.Getenv("N_CLIENT_ID"),
		ClientSecret:           os.Getenv("N_CLIENT_SECRET"),
		AuthURL:                "https://jp-notifications-lxe3vgwv.authentication.stagingaws.hanavlab.ondemand.com/oauth/token",
		ServiceURL:             "https://clm-sl-ans-canary-ans-service-api.cfapps.eu12.hana.ondemand.com",
		RateLimitingInterval:   2 * time.Second,
		MaxRequestsPerInterval: 5,
	}
	client := NewNotificationsClient(context.Background(), config, logger.With("component", "ANS-notificationsClient"))
	require.NotNil(t, client)
	recipient := notifications.NewRecipient("jaroslaw.pieszka@sap.com", "accounts.sap.com")
	require.NoError(t, recipient.Validate())
	require.NotNil(t, recipient)
	property := notifications.NewProperty("shoot", "c0123456")
	require.NoError(t, property.Validate())
	require.NotNil(t, property)
	notification := notifications.NewNotification("POC_WebOnlyType",
		[]notifications.Recipient{*recipient},
		notifications.WithProperties([]notifications.Property{*property}))
	require.NoError(t, notification.Validate())
	require.NotNil(t, notification)
	notificationAsJSON, err := json.Marshal(notification)
	require.JSONEq(t, `{
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
}`, string(notificationAsJSON))
	require.JSONEq(t, `{
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
}`, string(notificationAsJSON))
	require.NoError(t, err)
	err = client.postNotification(*notification)
	require.NoError(t, err)
}
