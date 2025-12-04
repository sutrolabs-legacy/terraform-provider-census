# Census Staging Environment Example

This example demonstrates how to use the Census Terraform provider with a **staging environment**. It creates a complete data pipeline with source, destination, and sync resources.

## Why Staging?

Testing Census configurations in staging allows you to:
- Test syncs without affecting production data
- Validate field mappings and transformations
- Experiment with new sync patterns safely
- Connect to sandbox versions of your destinations (e.g., Salesforce sandboxes)

## What This Example Creates

1. **Workspace** - A staging workspace for testing
2. **Source** - Connection to your staging data warehouse
3. **Destination** - Connection to a staging destination (e.g., Salesforce sandbox)
4. **Sync** - A data sync from source to destination (starts paused)

## Prerequisites

- Census staging account access
- Census Personal Access Token
- Staging data warehouse credentials
- Staging destination credentials (e.g., Salesforce sandbox)

## Quick Start

### 1. Configure Your Credentials

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your staging credentials
```

### 2. Key Configuration: Staging URL

The critical setting for staging is the `base_url` parameter:

```hcl
provider "census" {
  personal_access_token = var.census_personal_token
  base_url              = "https://app.staging.getcensus.com/api/v1"
}
```

This overrides the default production URL and routes all API calls to Census staging.

### 3. Deploy

```bash
terraform init
terraform plan
terraform apply
```

## Staging vs Production

### Provider Configuration

**Production:**
```hcl
provider "census" {
  personal_access_token = var.census_token
  # region defaults to "us" → https://app.getcensus.com/api/v1
}
```

**Staging:**
```hcl
provider "census" {
  personal_access_token = var.census_token
  base_url              = "https://app.staging.getcensus.com/api/v1"
}
```

### Environment Variables Approach

You can also use environment variables:

```bash
# Set your token
export CENSUS_PERSONAL_ACCESS_TOKEN="your-token"

# For production (default)
terraform apply

# For staging (override in terraform.tfvars or provider block)
```

## Configuration Options

### Source Types Supported

- `postgres` - PostgreSQL
- `snowflake` - Snowflake
- `bigquery` - Google BigQuery
- `redshift` - Amazon Redshift
- `databricks` - Databricks

### Destination Types Supported

- `salesforce` - Salesforce (set `sandbox = true` for sandboxes, use boolean not string)
- `hubspot` - HubSpot
- `intercom` - Intercom
- `postgres` - PostgreSQL (can be used as destination)

## Example Configurations

### Snowflake Source

```hcl
source_type = "snowflake"
source_credentials = {
  account   = "staging-account.snowflakecomputing.com"
  username  = "census_staging"
  password  = "your-password"
  warehouse = "STAGING_WH"
  database  = "STAGING_DB"
  schema    = "PUBLIC"
}
```

### BigQuery Source

```hcl
source_type = "bigquery"
source_credentials = {
  project_id   = "staging-project-id"
  private_key  = file("staging-service-account.json")
  client_email = "census@staging-project.iam.gserviceaccount.com"
}
```

### HubSpot Destination

```hcl
destination_type = "hubspot"
destination_credentials = {
  access_token = "your-hubspot-staging-token"
}
```

## Testing Workflow

### 1. Create Infrastructure (Paused)

The example creates the sync in a **paused** state by default:

```hcl
sync_paused = true
```

This allows you to verify the configuration before running any syncs.

### 2. Verify in Census UI

After applying, visit Census staging dashboard to verify:
- Source connection test passes
- Destination connection test passes
- Field mappings look correct

### 3. Run Test Sync

Once verified, you can either:

**Option A: Update Terraform**
```hcl
sync_paused = false
```
Then run `terraform apply`

**Option B: Trigger in Census UI**
Manually trigger a test run from the Census staging dashboard

### 4. Promote to Production

Once tested in staging, you can:
1. Create a separate production configuration
2. Change `base_url` to production (or remove it)
3. Update credentials to production values
4. Apply to production environment

## Important Notes

### Separate State Files

Keep staging and production state files separate:

```bash
# Staging
terraform workspace new staging
terraform apply

# Production
terraform workspace new production
terraform apply
```

Or use separate directories:
```
environments/
  ├── staging/
  │   ├── main.tf
  │   └── terraform.tfvars
  └── production/
      ├── main.tf
      └── terraform.tfvars
```

### Credential Management

- Never commit `terraform.tfvars` to version control
- Use different credentials for staging vs production
- Consider using Terraform Cloud or similar for secure credential storage

### Syncs Start Paused

By default, syncs are created in a paused state (`sync_paused = true`) to prevent accidental data sync runs during testing.

## Troubleshooting

### Connection Test Failures

If source or destination connections fail:
1. Verify credentials are correct for staging environment
2. Check network connectivity and firewall rules
3. For Salesforce sandboxes, ensure `sandbox = true` is set (boolean, not string)
4. Verify the staging API URL is accessible

### API Errors

If you get API errors:
- Verify your token has access to Census staging
- Check that `base_url` is exactly: `https://app.staging.getcensus.com/api/v1`
- Ensure token has required permissions

### Field Mapping Issues

If field mappings fail validation:
- Verify source table exists and has the specified columns
- Check destination object accepts the mapped fields
- Test the connection and run "Test Connection" in Census UI

## Cleaning Up

To remove all resources:

```bash
terraform destroy
```

This will delete:
- The sync
- The destination connection
- The source connection
- The workspace

## Next Steps

After testing in staging:
1. Document any configuration changes needed
2. Create production configuration
3. Review and approve production deployment
4. Monitor syncs in production

## Support

- Census Documentation: https://docs.getcensus.com/
- Provider Issues: https://github.com/sutrolabs/terraform-provider-census/issues