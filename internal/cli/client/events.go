package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kyma-project/ans-manager/internal/cli/credentials"
	"golang.org/x/oauth2/clientcredentials"
)

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
