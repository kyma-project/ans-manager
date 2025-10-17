package credentials

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewProviderFromJSON(t *testing.T) {
	t.Run("absolute path", func(t *testing.T) {
		absPath, _ := filepath.Abs("testdata/credentials.json")
		provider, err := NewProviderFromJSON(absPath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if provider == nil {
			t.Fatal("expected provider to be non-nil")
		}
		p := provider.(*ProviderFromJSON)
		if p.JSONFilePath != absPath {
			t.Errorf("expected path %s, got %s", absPath, p.JSONFilePath)
		}
	})

	t.Run("relative path", func(t *testing.T) {
		relPath := "testdata/credentials.json"
		provider, err := NewProviderFromJSON(relPath)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		p := provider.(*ProviderFromJSON)
		if !filepath.IsAbs(p.JSONFilePath) {
			t.Error("expected absolute path")
		}
	})
}

func TestProviderFromJSON_ProvideForRegion(t *testing.T) {
	testDataPath := "testdata/credentials.json"

	t.Run("valid region cf-eu10-canary", func(t *testing.T) {
		provider, err := NewProviderFromJSON(testDataPath)
		if err != nil {
			t.Fatalf("failed to create provider: %v", err)
		}

		oauth, err := provider.ProvideForRegion("cf-eu10-canary")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expected := &RegionAccess{
			ServiceEndpoint: ServiceEndpoint{
				ServiceURL:   "https://notifications.cf-eu10-canary.test.com/api/v1",
				SubaccountID: "your-subaccount-id-eu10",
			},
			OAuthCredentials: OAuthCredentials{
				ClientID:     "your-client-id-eu10",
				ClientSecret: "your-client-secret-eu10",
				TokenURL:     "https://oauth.cf-eu10-canary.test.com/oauth/token",
			},
		}

		if *oauth != *expected {
			t.Errorf("expected %+v, got %+v", expected, oauth)
		}
	})

	t.Run("valid region cf-us31", func(t *testing.T) {
		provider, err := NewProviderFromJSON(testDataPath)
		if err != nil {
			t.Fatalf("failed to create provider: %v", err)
		}

		oauth, err := provider.ProvideForRegion("cf-us31")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expected := &RegionAccess{
			ServiceEndpoint: ServiceEndpoint{
				ServiceURL:   "https://notifications.cf-us31.test.com/api/v1",
				SubaccountID: "your-subaccount-id-us31",
			},
			OAuthCredentials: OAuthCredentials{
				ClientID:     "your-client-id-us31",
				ClientSecret: "your-client-secret-us31",
				TokenURL:     "https://oauth.cf-us31.test.com/oauth/token",
			},
		}

		if *oauth != *expected {
			t.Errorf("expected %+v, got %+v", expected, oauth)
		}
	})

	t.Run("nonexistent region", func(t *testing.T) {
		provider, err := NewProviderFromJSON(testDataPath)
		if err != nil {
			t.Fatalf("failed to create provider: %v", err)
		}

		oauth, err := provider.ProvideForRegion("nonexistent")
		if err == nil {
			t.Error("expected error for empty credentials")
		}

		if oauth != nil {
			t.Errorf("expected nil oauth for nonexistent region, got %+v", oauth)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		provider, err := NewProviderFromJSON("/nonexistent/file.json")
		if err != nil {
			t.Fatalf("failed to create provider: %v", err)
		}

		_, err = provider.ProvideForRegion("cf-eu10-canary")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		invalidJSON := `{"invalid": json}`
		tmpFile := createTempFile(t, invalidJSON)
		defer os.Remove(tmpFile)

		provider, err := NewProviderFromJSON(tmpFile)
		if err != nil {
			t.Fatalf("failed to create provider: %v", err)
		}

		_, err = provider.ProvideForRegion("cf-eu10-canary")
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("missing credentials in valid JSON", func(t *testing.T) {
		incompleteJSON := `{
			"cf-test": {
				"client_id": "test-id"
			}
		}`
		tmpFile := createTempFile(t, incompleteJSON)
		defer os.Remove(tmpFile)

		provider, err := NewProviderFromJSON(tmpFile)
		if err != nil {
			t.Fatalf("failed to create provider: %v", err)
		}

		_, err = provider.ProvideForRegion("cf-test")
		if err == nil {
			t.Error("expected error for incomplete credentials")
		}
	})
}

func TestValidateRegionCredentials(t *testing.T) {
	t.Run("valid credentials", func(t *testing.T) {
		oauth := &RegionAccess{
			ServiceEndpoint: ServiceEndpoint{
				ServiceURL:   "https://api.example.com",
				SubaccountID: "test-subaccount-id",
			},
			OAuthCredentials: OAuthCredentials{
				ClientID:     "test-id",
				ClientSecret: "test-secret",
				TokenURL:     "https://oauth.example.com/token",
			},
		}

		err := validateRegionCredentials(oauth)
		if err != nil {
			t.Errorf("expected no error for valid credentials, got %v", err)
		}
	})

	t.Run("missing client_id", func(t *testing.T) {
		oauth := &RegionAccess{
			ServiceEndpoint: ServiceEndpoint{
				ServiceURL:   "https://api.example.com",
				SubaccountID: "test-subaccount-id",
			},
			OAuthCredentials: OAuthCredentials{
				ClientSecret: "test-secret",
				TokenURL:     "https://oauth.example.com/token",
			},
		}

		err := validateRegionCredentials(oauth)
		if err == nil {
			t.Error("expected error for missing client_id")
		}
	})

	t.Run("missing client_secret", func(t *testing.T) {
		oauth := &RegionAccess{
			ServiceEndpoint: ServiceEndpoint{
				ServiceURL:   "https://api.example.com",
				SubaccountID: "test-subaccount-id",
			},
			OAuthCredentials: OAuthCredentials{
				ClientID: "test-id",
				TokenURL: "https://oauth.example.com/token",
			},
		}

		err := validateRegionCredentials(oauth)
		if err == nil {
			t.Error("expected error for missing client_secret")
		}
	})

	t.Run("missing token_url", func(t *testing.T) {
		oauth := &RegionAccess{
			ServiceEndpoint: ServiceEndpoint{
				ServiceURL:   "https://api.example.com",
				SubaccountID: "test-subaccount-id",
			},
			OAuthCredentials: OAuthCredentials{
				ClientID:     "test-id",
				ClientSecret: "test-secret",
			},
		}

		err := validateRegionCredentials(oauth)
		if err == nil {
			t.Error("expected error for missing token_url")
		}
	})

	t.Run("missing service_url", func(t *testing.T) {
		oauth := &RegionAccess{
			ServiceEndpoint: ServiceEndpoint{
				SubaccountID: "test-subaccount-id",
			},
			OAuthCredentials: OAuthCredentials{
				ClientID:     "test-id",
				ClientSecret: "test-secret",
				TokenURL:     "https://oauth.example.com/token",
			},
		}

		err := validateRegionCredentials(oauth)
		if err == nil {
			t.Error("expected error for missing service_url")
		}
	})

	t.Run("multiple missing fields", func(t *testing.T) {
		oauth := &RegionAccess{
			OAuthCredentials: OAuthCredentials{
				ClientID: "test-id",
			},
		}

		err := validateRegionCredentials(oauth)
		if err == nil {
			t.Error("expected error for multiple missing fields")
		}
	})
}

func TestDecodeRegionCredentials(t *testing.T) {
	t.Run("valid JSON structure with testdata", func(t *testing.T) {
		file, err := os.Open("testdata/credentials.json")
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer file.Close()

		oauth, err := decodeRegionCredentials(file, "cf-eu12")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if oauth.ClientID != "your-client-id-eu12" {
			t.Errorf("expected client_id 'your-client-id-eu12', got '%s'", oauth.ClientID)
		}

		if oauth.TokenURL != "https://oauth.cf-eu12.test.com/oauth/token" {
			t.Errorf("expected correct token_url, got '%s'", oauth.TokenURL)
		}
	})

	t.Run("invalid JSON structure", func(t *testing.T) {
		invalidJSON := `["not", "an", "object"]`
		tmpFile := createTempFile(t, invalidJSON)
		defer os.Remove(tmpFile)

		file, err := os.Open(tmpFile)
		if err != nil {
			t.Fatalf("failed to open test file: %v", err)
		}
		defer file.Close()

		_, err = decodeRegionCredentials(file, "cf-eu10-canary")
		if err == nil {
			t.Error("expected error for invalid JSON structure")
		}
	})
}

func createTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test-credentials-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}
