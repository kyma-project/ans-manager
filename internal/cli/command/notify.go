package command

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/kyma-project/ans-manager/events"
	"github.com/kyma-project/ans-manager/internal/cli/client"
	"github.com/kyma-project/ans-manager/internal/cli/credentials"
)

const (
	notificationTypeKey    = "KEB_ANS_PoC_Always_Visible"
	iasHost                = "accounts.sap.com"
	resourceType           = "application"
	resourceName           = "ans-cli"
	eventTypeUser          = "ans-cli-user-notification"
	eventTypeSubaccount    = "ans-cli-subaccount-notification"
	eventTypeGlobalaccount = "ans-cli-globalaccount-notification"
)

var (
	credentialsFile string
	region          string
	subject         string
	body            string
)

func NewNotifyCommand() *cobra.Command {
	notifyCmd := &cobra.Command{
		Use:   "notify [type] [recipients...]",
		Short: "Send a notification event to ANS service",
		Long: `Send a notification event to the SAP BTP Alert Notification Service.

This command sends notifications to different recipient types:
- user: Send to user email addresses
- subaccount: Send to subaccount UUIDs  
- globalaccount: Send to global account UUIDs

Usage:
  ans-cli notify user email1@example.com [email2@example.com ...]
  ans-cli notify subaccount uuid1 [uuid2 ...]
  ans-cli notify globalaccount uuid1 [uuid2 ...]

The command uses OAuth authentication for the specified region.`,
		Args: cobra.MinimumNArgs(2), // At least type and one recipient
		RunE: func(cmd *cobra.Command, args []string) error {
			recipientType := args[0]
			recipients := args[1:]
			return runNotifyCommand(recipientType, recipients)
		},
	}

	notifyCmd.Flags().StringVarP(&region, "region", "r", "", "ANS service region (e.g., cf-eu10-canary, cf-eu12, cf-us31)")
	notifyCmd.MarkFlagRequired("region")

	notifyCmd.Flags().StringVarP(&subject, "subject", "s", "", "Event subject")
	notifyCmd.MarkFlagRequired("subject")

	notifyCmd.Flags().StringVarP(&body, "body", "b", "", "Event body/message")
	notifyCmd.MarkFlagRequired("body")

	notifyCmd.Flags().StringVarP(&credentialsFile, "credentials", "c", "credentials.json", "Path to credentials.json file")

	return notifyCmd
}

func runNotifyCommand(recipientType string, recipients []string) error {
	if err := validateRecipientType(recipientType); err != nil {
		return err
	}

	if err := validateRecipientsByType(recipientType, recipients); err != nil {
		return err
	}

	credentialsProvider, err := credentials.NewProviderFromJSON(credentialsFile)
	if err != nil {
		return fmt.Errorf("failed to create credentials provider: %w", err)
	}

	eventsClient, err := client.NewEventsClient(region, credentialsProvider)
	if err != nil {
		return fmt.Errorf("failed to create events client: %w", err)
	}

	resource := events.NewResource(resourceType, resourceName, "", "")

	var eventType string
	var visibility events.Visibility
	var eventRecipients *events.Recipients
	var userRecipients []events.UserRecipient
	var xsuaaRecipients []events.XsuaaRecipient

	switch recipientType {
	case "user":
		for _, email := range recipients {
			userRecipient := events.NewUserRecipient(email, iasHost)
			userRecipients = append(userRecipients, *userRecipient)
		}
		eventRecipients = events.NewRecipients(xsuaaRecipients, userRecipients)
		eventType = eventTypeUser
		visibility = events.VisibilityOwnerGlobalAccount
	case "subaccount":
		for _, uuid := range recipients {
			xsuaaRecipient := events.NewXsuaaRecipient(events.LevelSubaccount, uuid, []events.RoleName{events.RoleSubaccountAdministrator})
			xsuaaRecipients = append(xsuaaRecipients, *xsuaaRecipient)
		}
		eventRecipients = events.NewRecipients(xsuaaRecipients, userRecipients)
		eventType = eventTypeSubaccount
		visibility = events.VisibilityOwnerSubaccount
	case "globalaccount":
		for _, uuid := range recipients {
			xsuaaRecipient := events.NewXsuaaRecipient(events.LevelGlobalAccount, uuid, []events.RoleName{events.RoleGlobalAccountAdministrator})
			xsuaaRecipients = append(xsuaaRecipients, *xsuaaRecipient)
		}
		eventRecipients = events.NewRecipients(xsuaaRecipients, userRecipients)
		eventType = eventTypeGlobalaccount
		visibility = events.VisibilityOwnerGlobalAccount
	default:
		return fmt.Errorf("unsupported recipient type: %s", recipientType)
	}

	notificationMapping := events.NewNotificationMapping(notificationTypeKey, *eventRecipients)

	resourceEvent, err := events.NewResourceEvent(
		eventType,
		body,
		subject,
		resource,
		events.SeverityInfo,
		events.CategoryNotification,
		visibility,
		*notificationMapping,
	)
	if err != nil {
		return fmt.Errorf("failed to create resource event: %w", err)
	}

	if err := eventsClient.SendEvent(resourceEvent); err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	fmt.Printf("Event sent successfully to %s region\n", region)
	fmt.Printf("Recipients (%s): %v\n", recipientType, recipients)
	return nil
}

func validateRecipientType(recipientType string) error {
	if recipientType != "user" &&
		recipientType != "globalaccount" &&
		recipientType != "subaccount" {
		return fmt.Errorf("invalid recipient type: %s. Must be one of: user, globalaccount, subaccount", recipientType)
	}
	return nil
}

func validateRecipientsByType(recipientType string, recipients []string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("at least one recipient must be provided")
	}

	switch recipientType {
	case "user":
		return validateEmails(recipients)
	case "subaccount", "globalaccount":
		return validateUUIDs(recipients)
	default:
		return fmt.Errorf("unknown recipient type: %s", recipientType)
	}
}

func validateEmails(emails []string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	for _, email := range emails {
		if !emailRegex.MatchString(email) {
			return fmt.Errorf("invalid email address: %s", email)
		}
	}
	return nil
}

func validateUUIDs(uuids []string) error {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

	for _, uuid := range uuids {
		if !uuidRegex.MatchString(uuid) {
			return fmt.Errorf("invalid UUID format: %s", uuid)
		}
	}
	return nil
}
