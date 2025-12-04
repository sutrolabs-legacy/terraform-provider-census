# Complete Census Setup Example

This example demonstrates how to use the Census Terraform provider to manage a complete multi-workspace Census setup. It supports flexible workspace organization patterns and handles both new and existing workspaces.

## Features

- ✅ **Flexible Workspace Structure**: Define any workspace pattern (marketing_prod, revops_staging, etc.)
- ✅ **Mixed Management**: Combine Terraform-created and existing workspaces
- ✅ **API Key Handling**: Automatic key management for created workspaces, manual for existing
- ✅ **Multi-Environment Support**: Production, staging, development, or any custom environments
- ✅ **Environment-Specific Configurations**: Different connections for staging vs production

## Quick Start

### 1. Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- Census account with appropriate permissions
- Census Personal Access Token (for workspace management)
- Census Workspace Access Token (for sources/destinations/syncs)

### 2. Get Your Census Tokens

#### Personal Access Token (PAT)
1. Go to Census Dashboard → Settings → Developer → Personal Access Tokens
2. Create a new token with workspace management permissions
3. Copy the token (starts with `census_pat_...`)

#### Workspace API Keys (Per Workspace)
- **For new workspaces**: API keys are automatically generated when Terraform creates the workspace
- **For existing workspaces**: Get the API key from Census Dashboard → Workspace Settings → API → Access Token for each workspace you want to use

### 3. Configure Your Setup

```bash
# Clone or copy this example
cp terraform.tfvars.example terraform.tfvars

# Edit terraform.tfvars with your tokens and configuration
vim terraform.tfvars
```

### 4. Deploy

```bash
# Build the provider (if using development version)
cd ../../ && make dev

# Return to example directory
cd examples/complete-census-setup/

# Initialize and deploy
terraform init
terraform plan
terraform apply
```

## Configuration Options

### Basic Configuration

The minimal required configuration in `terraform.tfvars`:

```hcl
# Authentication
census_personal_token = "census_pat_your_personal_token_here"
census_region = "us"  # or "eu"

# No single workspace token needed - tokens are handled per workspace!

# Workspaces are configured in the workspaces section below

# Data connections
source_type = "postgres"  # or snowflake, bigquery, etc.
source_connection = {
  host     = "your-database.amazonaws.com"
  port     = 5432  # Numbers work directly, no need for quotes
  username = "census_user"
  password = "your-password"
  database = "analytics"
}

destination_type = "salesforce"  # or hubspot, postgres, etc.
destination_connection = {
  username       = "your-sf-user@company.com"
  password       = "your-sf-password"
  security_token = "your-sf-security-token"
  sandbox        = false  # Booleans work directly!
}

# Auto-refresh metadata after creating/updating connections
enable_auto_refresh = true  # Set to false to manually refresh schemas in Census UI
```

### Flexible Workspace Configuration

Configure workspaces for your organization structure:

#### Option 1: Marketing + RevOps Teams
```hcl
workspaces = {
  marketing_prod = {
    create = true
    name   = "Marketing Production"
    notification_emails = ["marketing-data@company.com"]
  }
  marketing_staging = {
    create = true
    name   = "Marketing Staging"
    notification_emails = ["marketing-dev@company.com"]
  }
  revops_prod = {
    create = true
    name   = "Revenue Operations Production"
    notification_emails = ["revops@company.com"]
  }
}
```

#### Option 2: Mixed Existing + New Workspaces
```hcl
workspaces = {
  # Existing production workspace (managed outside Terraform)
  marketing_prod = {
    create      = false
    existing_id = "12345"
    api_key     = "ws_live_abc123..."  # Get from Census Dashboard → Workspace Settings → API
  }
  # New staging workspace (created by Terraform)
  marketing_staging = {
    create = true
    name   = "Marketing Staging Environment"
    notification_emails = ["marketing-dev@company.com"]
    # api_key will be auto-generated on creation and stored in Terraform state
  }
}
```

#### Option 3: Enterprise Multi-Team Structure
```hcl
workspaces = {
  sales_prod      = { create = true, name = "Sales Production" }
  sales_staging   = { create = true, name = "Sales Staging" }
  marketing_prod  = { create = true, name = "Marketing Production" }
  marketing_dev   = { create = true, name = "Marketing Development" }
  finance_prod    = { create = true, name = "Finance Production" }
  support_prod    = { create = true, name = "Customer Support Production" }
  analytics_prod  = { create = true, name = "Analytics Team Production" }
}
```

### Environment-Specific Connections

Configure different connections for staging environments:

```hcl
# Production database
source_connection = {
  host     = "prod-db.company.com"
  port     = 5432
  username = "census_prod"
  password = "prod-password"
  database = "analytics"
}

# Staging database (optional - falls back to production config if not set)
staging_source_connection = {
  host     = "staging-db.company.com"
  port     = 5432
  username = "census_staging"
  password = "staging-password"
  database = "staging_analytics"
}

# Production Salesforce
destination_connection = {
  username       = "census@company.com"
  password       = "prod-password"
  security_token = "prod-token"
  sandbox        = false
}

# Staging Salesforce (optional)
staging_destination_connection = {
  username       = "census@company.com.staging"
  password       = "staging-password"
  security_token = "staging-token"
  sandbox        = true
}
```

## Supported Source Types

- `snowflake`
- `bigquery`
- `postgres`
- `redshift`
- `databricks`
- `mysql`
- `sql_server`
- `oracle`
- `clickhouse`

## Supported Destination Types

- `salesforce`
- `hubspot`
- `postgres`
- `airtable`
- `intercom`
- `zendesk`
- `mailchimp`
- `segment`
- `mixpanel`
- `amplitude`

## Understanding API Key Management

### For New Workspaces (create = true)
- API key is automatically generated and captured during workspace creation
- Key is stored in Terraform state
- Subsequent `terraform plan/apply` operations use the key from state
- **Important**: The API key is only returned on creation - if you lose state, you'll need the key from Census dashboard

### For Existing Workspaces (create = false)
- You must provide the `api_key` in your workspace configuration
- Get the key from Census Dashboard → Workspace Settings → API → Access Token
- Each existing workspace needs its own API key

### API Key Persistence
- After workspace creation, the API key persists in Terraform state
- Team members sharing state will have access to the keys
- For production deployments, consider:
  - Using Terraform Cloud/Enterprise for secure state storage
  - Implementing key rotation procedures
  - Using workspace-specific API keys for enhanced isolation

## Common Connection Examples

### Snowflake (Password Authentication)
```hcl
source_type = "snowflake"
source_connection = {
  account   = "your-account.snowflakecomputing.com"
  username  = "census_user"
  password  = "your-password"
  warehouse = "COMPUTE_WH"
  database  = "ANALYTICS"
  schema    = "PUBLIC"
}
```

### Snowflake (Keypair Authentication)
```hcl
source_type = "snowflake"
source_connection = {
  account                = "abc12345.us-east-1"
  warehouse              = "COMPUTE_WH"
  database               = "PRODUCTION"
  username               = "census_user"
  role                   = "CENSUS_ROLE"
  use_keypair            = true  # Boolean works directly!
  private_key_pkcs8      = var.snowflake_private_key
  private_key_passphrase = var.snowflake_key_passphrase  # Optional
}
```

### BigQuery
```hcl
source_type = "bigquery"
source_connection = {
  project_id   = "your-gcp-project-id"
  private_key  = jsonencode(var.bigquery_service_account_key)
  client_email = "census@your-project.iam.gserviceaccount.com"
}
```

### HubSpot
```hcl
destination_type = "hubspot"
destination_connection = {
  access_token = "your-hubspot-access-token"
}
```

## Troubleshooting

### Common Issues

**Provider not found**
```bash
# Build the development provider
cd ../../ && make dev
cd examples/complete-census-setup/
terraform init
```

**Token permissions errors**
- Ensure Personal Access Token has workspace management permissions
- Ensure Workspace Access Token has appropriate permissions for the target workspace

**API key errors for existing workspaces**
- Verify `api_key` is correct and from the right workspace
- Check that `existing_id` matches the workspace ID in Census

**Connection test failures**
- Verify database/service credentials are correct
- Check network connectivity and firewall rules
- For staging environments, ensure sandbox settings are correct

### Validation

```bash
# Validate configuration syntax
terraform validate

# Plan without applying to check for issues
terraform plan
```

### Outputs

After successful deployment, you'll see outputs with:
- **workspaces_info**: Details about all workspaces (IDs, names, API keys)
- **sources_info**: Information about created data sources
- **destinations_info**: Information about created destinations  
- **setup_summary**: Overview of the complete setup with next steps

## Next Steps

1. **Verify Connections**: Check Census dashboard to ensure all sources and destinations are working
2. **Configure Syncs**: Set up data syncs between your sources and destinations
3. **Monitor**: Set up monitoring and alerting for your data pipelines
4. **Scale**: Add more workspaces, sources, and destinations as needed

## Support

- **Provider Issues**: [GitHub Issues](https://github.com/sutrolabs/terraform-provider-census/issues)
- **Census Documentation**: [Census Docs](https://docs.getcensus.com/)
- **Terraform Documentation**: [Terraform Docs](https://www.terraform.io/docs/)

## License

This example is provided as-is for demonstration purposes.