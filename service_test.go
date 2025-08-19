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
	require.NotEmpty(t, os.Getenv("E_CLIENT_ID"))
	require.NotEmpty(t, os.Getenv("E_CLIENT_SECRET"))
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
	event, err := events.NewResourceEvent(
		"eventType",
		"body",
		"subject",
		events.NewResource("broker",
			"keb",
			"2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
			"2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
			events.WithResourceGlobalAccount("8cd57dc2-edb2-45e0-af8b-7d881006e516")),
		events.SeverityInfo,
		events.CategoryNotification,
		events.VisibilityOwnerSubAccount,
		*events.NewNotificationMapping("POC_WebOnlyType2",
			*events.NewRecipients(
				[]events.XsuaaRecipient{*events.NewXsuaaRecipient(
					events.LevelSubaccount, "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
					[]events.RoleName{"Subaccount admin"})},
				nil)),
	)
	require.NotNil(t, event)
	require.NoError(t, event.Validate())
	eventAsJSON, err := json.Marshal(event)
	require.JSONEq(t, `{
  "body" : "body",
  "subject" : "subject",
  "eventType" : "eventType",
  "resource" : {
    "resourceType" : "broker",
    "resourceName" : "keb",
    "subAccount" : "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
    "globalAccount" : "8cd57dc2-edb2-45e0-af8b-7d881006e516",
    "resourceGroup" : "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad"
  },
  "severity" : "INFO",
  "category" : "NOTIFICATION",
  "visibility" : "OWNER_SUBACCOUNT",
  "notificationMapping" : {
    "notificationTypeKey" : "POC_WebOnlyType2",
    "recipients" : {
      "xsuaa" : [ {
        "level" : "SUBACCOUNT",
        "tenantId" : "2fd47ed4-dd54-40b5-99d8-36c4dc3b8cad",
        "roleNames" : [ "Subaccount admin" ]
      } ]
    }
  }
}`, string(eventAsJSON))
	err = client.postEvent(*event)
	require.NoError(t, err)
}

func Test_PostNotifications(t *testing.T) {
	t.Skip()
	require.NotEmpty(t, os.Getenv("N_CLIENT_ID"))
	require.NotEmpty(t, os.Getenv("N_CLIENT_SECRET"))
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
	notification := notifications.NewNotification("POC_WebOnlyType",
		[]notifications.Recipient{*notifications.NewRecipient("jaroslaw.pieszka@sap.com", "accounts.sap.com")},
		notifications.WithProperties([]notifications.Property{*notifications.NewProperty("shoot", "c0123456")}))
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
	require.NoError(t, err)
	err = client.postNotification(*notification)
	require.NoError(t, err)
}
