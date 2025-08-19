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
- If you intend to use events create Service Instance of ANS Manager with plan `standard`. Then create a ticket for ANS team to elevate this Instance (see [this](https://jira.tools.sap/browse/HCPSR-52322)). 
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
- Download Service Bindings as JSON files. Provided shell scripts assume that the files are named `stage-e-sb.json` and `stage-n-sb.json`.

## Usage

### Get tokens
You need separate tokens for each Service Binding. The tokens are used to authenticate the requests to the ANS Manager API.
- Run the script `refresh-tokens.sh` to get the tokens for both Service Bindings. The script will use the Service Binding files you downloaded in the previous step.
- The script will create two files:  `current-e.token` and `current-n.token` with the tokens for both Service Bindings.


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
