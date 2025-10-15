package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Provider interface {
	ProvideForRegion(region string) (*OAuthCredentials, error)
}

type OAuthCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TokenURL     string `json:"token_url"`
	ServiceURL   string `json:"service_url"`
}

type ProviderFromJSON struct {
	JSONFilePath string
}

func NewProviderFromJSON(jsonFilePath string) (Provider, error) {
	if !filepath.IsAbs(jsonFilePath) {
		absPath, err := filepath.Abs(jsonFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for credentials JSON file: %w", err)
		}
		jsonFilePath = absPath
	}
	return &ProviderFromJSON{
		JSONFilePath: jsonFilePath,
	}, nil
}

func (p *ProviderFromJSON) ProvideForRegion(region string) (*OAuthCredentials, error) {
	return p.getOAuthCredsFromJSON(region)
}

func (p *ProviderFromJSON) getOAuthCredsFromJSON(region string) (*OAuthCredentials, error) {
	file, err := os.Open(p.JSONFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open credentials file %s: %w", p.JSONFilePath, err)
	}
	defer file.Close()

	regionCredentials, err := decodeRegionCredentials(file, region)
	if err != nil {
		return nil, err
	}

	if err := validateRegionCredentials(regionCredentials); err != nil {
		return nil, fmt.Errorf("invalid credentials for region %s: %w", region, err)
	}

	return regionCredentials, nil
}

func decodeRegionCredentials(file *os.File, region string) (*OAuthCredentials, error) {
	var regionCredentials OAuthCredentials
	dec := json.NewDecoder(file)
	t, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON: %w", err)
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return nil, fmt.Errorf("expected start of a JSON object ('{' delimiter)")
	}

	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("failed to read key: %w", err)
		}
		key, ok := keyToken.(string)
		if !ok {
			return nil, fmt.Errorf("expected string key, got %T", keyToken)
		}
		if key == region {
			if err := dec.Decode(&regionCredentials); err != nil {
				return nil, fmt.Errorf("failed to decode OAuth credentials: %w", err)
			}
			return &regionCredentials, nil
		} else {
			var skip json.RawMessage
			if err := dec.Decode(&skip); err != nil {
				return nil, fmt.Errorf("failed to skip value: %w", err)
			}
		}
	}

	return nil, fmt.Errorf("region %s not found in credentials file", region)
}

func validateRegionCredentials(regionCredentials *OAuthCredentials) error {
	var errs []error
	if regionCredentials.ClientID == "" {
		errs = append(errs, fmt.Errorf("client_id is empty"))
	}
	if regionCredentials.ClientSecret == "" {
		errs = append(errs, fmt.Errorf("client_secret is empty"))
	}
	if regionCredentials.TokenURL == "" {
		errs = append(errs, fmt.Errorf("token_url is empty"))
	}
	if regionCredentials.ServiceURL == "" {
		errs = append(errs, fmt.Errorf("service_url is empty"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
