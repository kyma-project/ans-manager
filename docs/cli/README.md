# ANS Manager CLI

A simplified command-line interface for sending notifications to the SAP BTP Alert Notification Service (ANS).

## Features

- Simple interface with only essential parameters
- OAuth2 authentication with automatic token management
- Support for multiple regions (cf-eu10-canary, cf-eu12, cf-us31, etc.)
- Secure credential management via external JSON file
- Support for three recipient types: user, subaccount, globalaccount

## Installation

Build the CLI from source:

```bash
make build
```

Or manually:

```bash
go build -o bin/ans-cli ./cmd/ans-cli
```

## Configuration

Create a `credentials.json` file with your OAuth configuration:

```json
{
  "cf-eu10-canary": {
    "subaccount_id": "your-subaccount-id-eu10",
    "client_id": "your-client-id-eu10",
    "client_secret": "your-client-secret-eu10",
    "token_url": "https://oauth.cf-eu10-canary.test.local/oauth/token",
    "service_url": "https://notifications.cf-eu10-canary.test.local/api/v1"
  },
  "cf-eu12": {
    "subaccount_id": "your-subaccount-id-eu12",
    "client_id": "your-client-id-eu12",
    "client_secret": "your-client-secret-eu12",
    "token_url": "https://oauth.cf-eu12.test.local/oauth/token",
    "service_url": "https://notifications.cf-eu12.test.local/api/v1"
  },
  "cf-us31": {
    "subaccount_id": "your-subaccount-id-us31",
    "client_id": "your-client-id-us31",
    "client_secret": "your-client-secret-us31",
    "token_url": "https://oauth.cf-us31.test.local/oauth/token",
    "service_url": "https://notifications.cf-us31.test.local/api/v1"
  }
}
```

## Usage

### Basic Command Structure

```bash
./bin/ans-cli notify [type] [recipients...] [flags]
```

### Required Parameters

- `type`: Recipient type (user, globalaccount, subaccount) - first argument
- `recipients`: One or more recipients - remaining arguments
- `--region, -r`: ANS service region (e.g., cf-eu10-canary, cf-eu12, cf-us31)
- `--subject, -s`: Event subject/title
- `--body, -b`: Event body/message

### Optional Parameters

- `--credentials, -c`: Path to credentials.json (default: credentials.json)

### Recipient Types and Validation

#### user
Send notifications to user email addresses. Recipients must be valid email addresses.

#### subaccount
Send notifications to subaccount UUIDs. Recipients must be valid UUID format.

#### globalaccount
Send notifications to global account UUIDs. Recipients must be valid UUID format.

### Examples

#### Send notification to users via email:

```bash
./bin/ans-cli notify user admin@example.com developer@example.com \
  --region cf-eu12 \
  --subject "Application Deployed Successfully" \
  --body "Your application has been deployed to production"
```

#### Send notification to single user:

```bash
./bin/ans-cli notify user admin@example.com \
  --region cf-eu12 \
  --subject "System Alert" \
  --body "System maintenance will begin soon"
```

#### Send notification to subaccounts:

```bash
./bin/ans-cli notify subaccount 12345678-1234-1234-1234-123456789012 87654321-4321-4321-4321-210987654321 \
  --region cf-eu12 \
  --subject "Resource Alert" \
  --body "Resource usage has exceeded threshold"
```

#### Send notification to global accounts:

```bash
./bin/ans-cli notify globalaccount 11111111-2222-3333-4444-555555555555 \
  --region cf-eu10-canary \
  --subject "System Maintenance" \
  --body "Scheduled maintenance will begin at 2 AM UTC"
```

#### Using custom credentials file:

```bash
./bin/ans-cli notify user dba@example.com \
  --credentials /path/to/my/credentials.json \
  --region cf-us31 \
  --subject "Backup Completed" \
  --body "Database backup completed successfully"
```

## Security Features

- Credentials are loaded only when needed and not kept in memory
- OAuth tokens are automatically managed and refreshed
- Supports secure client credentials flow
- No sensitive data is logged or exposed

## Error Handling

The CLI provides detailed error messages for:
- Invalid recipient types
- Invalid email address formats (for user type)
- Invalid UUID formats (for subaccount/globalaccount types)
- Missing required parameters or arguments
- Authentication failures
- Network connectivity issues
- Invalid configuration
- Event validation errors
- API response errors

## Help

Get help for the CLI:

```bash
./bin/ans-cli --help
./bin/ans-cli notify --help
```