---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cyberarkoss_dbaccount Resource - cyberarkoss"
subcategory: ""
description: |-
  Database Account Resource
---

# cyberarkoss_dbaccount (Resource)

Database Account Resource

## Example Usage

```terraform
resource "cyberarkoss_dbaccount" "pgdb" {
  name             = "user-db"
  address          = "1.2.3.4"
  username         = "user-db"
  platform         = "PROD_PostgreSQL"
  safe             = "TF_TEST_SAFE"
  secrettype       = "password"
  secret           = "SincerelySecure2#24!"
  sm_manage        = false
  sm_manage_reason = "No CPM Associated with Safe."
  db_port          = "8432"
  db_dsn           = "dsn"
  dbname           = "dbo.services"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String) URI, URL or IP associated with the credential.
- `name` (String) Custom Account Name for customizing the object name in a safe.
- `platform` (String) Management Platform associated with the Database Credential.
- `safe` (String) Target Safe where the credential object will be onboarded.
- `secret` (String, Sensitive) Password of the credential object.
- `secrettype` (String) Secret type of credential, should always be password unless working with AWS Keys.
- `username` (String) Username of the Credential object.

### Optional

- `db_dsn` (String) Database data source name.
- `db_port` (String) Database connection port.
- `dbname` (String) Database name.
- `sm_manage` (Boolean) Automatic Management of a credential. Optional Value.
- `sm_manage_reason` (String) If sm_manage is false, provide reason why credential is not managed.

### Read-Only

- `id` (String) CyberArk Privilege Cloud Credential ID- Generated from CyberArk after onboarding account into a safe.
- `last_updated` (String)
