package factory

import (
	"fmt"

	"github.com/kyma-project/ans-manager/events"
)

const (
	notificationTypeKey    = "KEB_ANS_PoC_Always_Visible"
	iasHost                = "accounts.sap.com"
	eventTypeUser          = "ans-cli-user-notification"
	eventTypeSubaccount    = "ans-cli-subaccount-notification"
	eventTypeGlobalaccount = "ans-cli-globalaccount-notification"
)

func NewResourceEventByRecipients(recipientType string, recipients []string) (*events.ResourceEvent, error) {
	switch recipientType {
	case "user":
		return newUserResourceEvent(recipients), nil
	case "subaccount":
		return newSubaccountResourceEvent(recipients), nil
	case "globalaccount":
		return newGlobalAccountResourceEvent(recipients), nil
	default:
		return nil, fmt.Errorf("unsupported recipient type: %s", recipientType)
	}
}

func newUserResourceEvent(recipients []string) *events.ResourceEvent {
	var userRecipients []events.UserRecipient
	for _, email := range recipients {
		userRecipient := events.NewUserRecipient(email, iasHost)
		userRecipients = append(userRecipients, *userRecipient)
	}
	eventRecipients := events.NewRecipients(nil, userRecipients)
	notificationMapping := events.NewNotificationMapping(notificationTypeKey, *eventRecipients)

	return &events.ResourceEvent{
		EventType:           eventTypeUser,
		Visibility:          events.VisibilityOwnerGlobalAccount,
		NotificationMapping: *notificationMapping,
	}
}

func newSubaccountResourceEvent(recipients []string) *events.ResourceEvent {
	var xsuaaRecipients []events.XsuaaRecipient
	for _, uuid := range recipients {
		xsuaaRecipient := events.NewXsuaaRecipient(events.LevelSubaccount, uuid, []events.RoleName{events.RoleSubaccountAdministrator})
		xsuaaRecipients = append(xsuaaRecipients, *xsuaaRecipient)
	}
	eventRecipients := events.NewRecipients(xsuaaRecipients, nil)
	notificationMapping := events.NewNotificationMapping(notificationTypeKey, *eventRecipients)

	return &events.ResourceEvent{
		EventType:           eventTypeSubaccount,
		Visibility:          events.VisibilityOwnerSubaccount,
		NotificationMapping: *notificationMapping,
	}
}

func newGlobalAccountResourceEvent(recipients []string) *events.ResourceEvent {
	var xsuaaRecipients []events.XsuaaRecipient
	for _, uuid := range recipients {
		xsuaaRecipient := events.NewXsuaaRecipient(events.LevelGlobalAccount, uuid, []events.RoleName{events.RoleGlobalAccountAdministrator})
		xsuaaRecipients = append(xsuaaRecipients, *xsuaaRecipient)
	}
	eventRecipients := events.NewRecipients(xsuaaRecipients, nil)
	notificationMapping := events.NewNotificationMapping(notificationTypeKey, *eventRecipients)

	return &events.ResourceEvent{
		EventType:           eventTypeGlobalaccount,
		Visibility:          events.VisibilityOwnerGlobalAccount,
		NotificationMapping: *notificationMapping,
	}
}
