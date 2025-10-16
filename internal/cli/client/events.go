package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kyma-project/ans-manager/events"
	"github.com/kyma-project/ans-manager/internal/cli/credentials"
	"golang.org/x/oauth2/clientcredentials"
)

const eventsEndpoint = "/cf/producer/service/v1/resource-events"

type EventsClient struct {
	httpClient *http.Client
	serviceURL string
}

func NewEventsClient(region string, credentialsProvider credentials.Provider) (*EventsClient, error) {
	creds, err := credentialsProvider.ProvideForRegion(region)
	if err != nil {
		return nil, fmt.Errorf("failed to load OAuth credentials for region %s: %w", region, err)
	}

	oauthCfg := &clientcredentials.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		TokenURL:     creds.TokenURL,
	}

	ctx := context.Background()
	httpClient := oauthCfg.Client(ctx)

	return &EventsClient{
		httpClient: httpClient,
		serviceURL: creds.ServiceURL,
	}, nil
}

func (c *EventsClient) SendEvent(event *events.ResourceEvent, debug bool) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if debug {
		fmt.Println("DEBUG: Request payload:")
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, eventJSON, "", "  "); err != nil {
			fmt.Printf("Raw JSON: %s\n", string(eventJSON))
		} else {
			fmt.Printf("%s\n", prettyJSON.String())
		}
		fmt.Println()
	}

	url := fmt.Sprintf("%s%s", c.serviceURL, eventsEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(eventJSON))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}
