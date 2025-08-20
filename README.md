# ANS Manager
<!--- mandatory --->
Alert Notifications Service manager

## Overview
<!--- mandatory section --->

Simple client using ANS services to post notifications and events created for the sake of Proof of Concept, to explore possibilities of API, and get the knowledge of the process required to acquire the necessary 
resources, permissions.

## Prerequisites

- Create Service Instance of ANS Manager with plan `xproduct-single-tenant`. For Staging and Canary Service Instance must be created in the subaccount region `cf-us10`. This Instance is needed to process notifications. 
- Ask ANS team to assign storage passing them the Service Instance ID. The storage is required to store the notifications.
- If you intend to use events create Service Instance of ANS Manager with plan `standard`. Then create a ticket for ANS team to elevate this Instance. 
- Update the `standard` Service Instance with the `xproduct-single-tenant` instance ID.
```json
{
    "configuration": {
        "actions": [],
        "conditions": [],
        "subscriptions": []
    },
    "enableCloudControllerAuditEvents": false,
    "notificationServiceInstanceId": "5928f39e-fab3-44b2-a62a-82f3eb3bbf55"
}
```
> Note: In the ANS documentation there is an example with `enableCloudControllerAuditEvents` set to `true`. This didn't work on Staging.

## Installation

- Create Service Bindings for both Service Instances. The `xproduct-single-tenant` Service Binding is needed to post notifications, the `standard` Service Binding is needed to post events.
- Download Service Bindings as JSON files. Provided shell scripts assume that the files are named `stage-e-sb.json` (binding for `standard` plan) and `stage-n-sb.json` (binding for `xproduct-single-tenant` plan).

## Go service and client

Module `ans-manager` contains a simple Go service and client that can be used to post notifications and events. 
To create the service, you need to provide configuration for two endpoints: one for notifications and one for events.
```go
	EndpointConfig struct {
		ClientID               string
		ClientSecret           string
		AuthURL                string
		ServiceURL             string
		RateLimitingInterval   time.Duration `envconfig:"default=2s,optional"`
		MaxRequestsPerInterval int           `envconfig:"default=5,optional"`
	}

```
ClientID and ClientSecret are the credentials for the Service Binding. AuthURL is the URL for the authentication service, and ServiceURL is the URL for the ANS Manager API, both provided in Service Binding.
Since the ANS Manager API is rate-limited, you can configure the rate limiting for the client using `RateLimitingInterval` and `MaxRequestsPerInterval`. The default values are set to 2 seconds and 5 requests per interval, respectively.

There are two end-to-end tests in the `service-test.go` file that demonstrate how to use the clients to post notifications and events.
The tests do not use the service itself, but rather the clients directly. To run tests, you need to set the environment variables for the Service Binding credentials,
and you need to comment out the `t.Skip()` lines in the tests. The tests will post a notification and an event using the clients, and they will check if the notification and event are delivered correctly.

## Development

> Add instructions on how to develop the project or example. It must be clear what to do and, for example, how to trigger the tests so that other contributors know how to make their pull requests acceptable. Include the instructions or provide links to related documentation.

## Contributing
<!--- mandatory section - do not change this! --->

See the [Contributing Rules](CONTRIBUTING.md).

## Code of Conduct
<!--- mandatory section - do not change this! --->

See the [Code of Conduct](CODE_OF_CONDUCT.md) document.

## Licensing
<!--- mandatory section - do not change this! --->

See the [license](./LICENSE) file.
