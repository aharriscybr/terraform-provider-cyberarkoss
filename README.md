[![Release](https://github.com/aharriscybr/terraform-provider-cyberarkoss/actions/workflows/release.yml/badge.svg)](https://github.com/aharriscybr/terraform-provider-cyberarkoss/actions/workflows/release.yml)

# This is an unofficial Terraform Provider for interacting with CyberArk Privilege Cloud Resources.
# This is not production, or even development code. This should **only** be used in a test environment.


This release is currently in **beta**. 

Please see [Docs](/docs/index.md) for current supported resources and data providers


# Configurable Environment Variables
*We do not support configuring optional properties via environment variables, these must be implicitly defined.*

## Provider
- CYBERARK_PROVIDER_TENANT
- CYBERARK_PROVIDER_CLIENT_ID
- CYBERARK_PROVIDER_CLIENT_SECRET
- CYBERARK_PROVIDER_DOMAIN
## Accounts
- CYBERARK_ACCOUNT_CUSTOM_NAME
- CYBERARK_ACCOUNT_USERNAME
- CYBERARK_ACCOUNT_PLATFORM
- CYBERARK_ACCOUNT_SAFE
- CYBERARK_ACCOUNT_SECRETTYPE
- CYBERARK_ACCOUNT_SECRET

# Set up Terraform User
- Log into Identity Administration and navigate to the Users Widget

<img src="img/users-widget.png" width="60%" height="30%">

- Create New User

<img src="img/add-user-widget.png"  width="60%" height="30%">

- Populate User Data

<img src="img/terraform-user.png"  width="60%" height="30%">

- Navigate to the Roles Widget

<img src="img/roles-widget.png" width="60%" height="30%">

- Add the new user to the Privilege Cloud Safe Managers Role

<img src="img/priv-safe-manager.png" width="60%" height="30%">

- Search for the Terraform User and Add

<img src="img/add-terraform-user.png" width="60%" height="30%">

## Usage instructions

**TBD**
