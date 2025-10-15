package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Provider interface {
	ProvideForRegion(region string) (*OAuth, error)
}

type OAuth struct {
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

func (p *ProviderFromJSON) ProvideForRegion(region string) (*OAuth, error) {
	return p.getOAuthConfigFromJSON(region)
}

func (p *ProviderFromJSON) getOAuthConfigFromJSON(region string) (*OAuth, error) {
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

func decodeRegionCredentials(file *os.File, region string) (*OAuth, error) {
	decoder := json.NewDecoder(file)
	var regionCredentials OAuth
	t, err := decoder.Token()
	if err != nil || t != json.Delim('{') {
		return nil, fmt.Errorf("expected start of a JSON object ('{' delimiter)")
	}

	for decoder.More() {
		var key string
		if err := decoder.Decode(&key); err != nil {
			return nil, err
		}
		if key == region {
			if err := decoder.Decode(&regionCredentials); err != nil {
				return nil, err
			}
			break
		}
	}
	return &regionCredentials, nil
}

func validateRegionCredentials(regionCredentials *OAuth) error {
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
