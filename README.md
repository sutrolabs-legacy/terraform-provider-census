# Terraform Provider for Census

A Terraform provider for managing [Census](https://getcensus.com) resources. Census enables you to sync data from your warehouse to all your operational tools, and this provider allows you to manage Census infrastructure as code.

## Features

- **Multi-region support**: Works with US, EU, and AU Census regions
- **Complete Census workflow**: Manage workspaces, sources, datasets, destinations, and syncs
- **PAT-only authentication**: Uses personal access tokens with dynamic workspace token retrieval
- **Import support**: Import existing Census resources into Terraform state
- **Staging environment support**: Configure custom base URLs for testing

## Installation

```hcl
terraform {
  required_providers {
    census = {
      source  = "sutrolabs/census"
      version = "~> 0.2.0"
    }
  }
}

provider "census" {
  personal_access_token = var.census_personal_token
  region                = "us"  # or "eu", "au"
}
```

## Usage

```hcl
resource "census_workspace" "data_team" {
  name = "Data Team Workspace"
  notification_emails = ["data-alerts@company.com"]
}

resource "census_source" "warehouse" {
  workspace_id = census_workspace.data_team.id
  name         = "Production Warehouse"
  type         = "snowflake"

  connection_config = {
    account   = "xy12345.us-east-1"
    database  = "ANALYTICS"
    warehouse = "COMPUTE_WH"
    role      = "CENSUS_ROLE"
    username  = "census_user"
    password  = var.snowflake_password
  }
}

resource "census_sync" "users_to_crm" {
  workspace_id = census_workspace.data_team.id
  label        = "Users to CRM"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type         = "table"
      table_name   = "users"
      table_schema = "public"
    }
  }

  destination_attributes {
    connection_id = census_destination.crm.id
    object        = "Contact"
  }

  operation = "upsert"

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    from = "first_name"
    to   = "FirstName"
  }

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "daily"
        hour      = 8
        minute    = 0
      }
    }
  }
}
```

## Resources

- `census_workspace` - Manage Census workspaces
- `census_source` - Data warehouse connections (Snowflake, BigQuery, Postgres, etc.)
- `census_destination` - Business tool integrations (Salesforce, HubSpot, etc.)
- `census_dataset` - SQL datasets for data transformation
- `census_sync` - Data syncs between sources and destinations

## Data Sources

All resources have corresponding data sources for read-only operations. See [documentation](docs/) for details.

## Documentation

- [Resource Documentation](docs/resources/) - Detailed documentation for each resource
- [Data Source Documentation](docs/data-sources/) - Read-only data source documentation
- [Examples](examples/) - Complete working examples
- [CHANGELOG](CHANGELOG.md) - Version history and changes
- [Census API Documentation](https://developers.getcensus.com/api-reference/introduction/overview)

## Troubleshooting

### Enabling Debug Logging

When encountering errors or unexpected behavior, enable debug logging to see detailed information about API requests and responses:

```bash
# Enable debug logging for all Terraform operations
export TF_LOG=DEBUG
terraform apply

# Or for a single command
TF_LOG=DEBUG terraform apply

# Save logs to a file
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform-debug.log
terraform apply
```

Debug logs include:
- API response bodies (used for diagnosing API format issues)
- HTTP status codes and headers
- JSON unmarshaling details
- Resource state transitions

> **Note**: The provider does not log API request bodies (which contain your Terraform configuration).
> Debug logs primarily show API responses, which generally do not contain sensitive credentials.
> However, responses may include connection IDs and data schema information.

### Common Issues

#### "Provider produced inconsistent result after apply"

This error typically indicates an API response format mismatch between the Census API and the provider. To diagnose:

1. **Enable debug logging**:
   ```bash
   TF_LOG=DEBUG terraform apply 2>&1 | tee terraform-debug.log
   ```

2. **Look for the underlying error**:
   ```bash
   grep -A 10 "failed to decode response JSON" terraform-debug.log
   ```

3. **Check the raw API response**:
   ```bash
   grep -A 5 "Raw API response:" terraform-debug.log
   ```

4. **Report the issue** with:
   - The full error message from debug logs
   - The raw API response showing the format mismatch
   - Your Terraform configuration (with sensitive values redacted)

### Reporting Issues

When reporting issues, please include:

1. **Debug logs**:
   ```bash
   # Capture logs (safe to share - only contains API responses)
   TF_LOG=DEBUG terraform apply 2>&1 | tee terraform-debug.log
   ```

2. **Provider version**:
   ```bash
   terraform version
   ```

3. **Terraform configuration** (redact sensitive values):
   ```hcl
   resource "census_source" "example" {
     connection_config = {
       username = "census_user"
       password = "REDACTED"  # ← Redact before sharing
     }
   }
   ```

4. **Expected vs actual behavior**

5. **Steps to reproduce**

Submit issues to: [GitHub Issues](https://github.com/sutrolabs/terraform-provider-census/issues)

## Contributing

At this time we are not accepting external contributions to the provider. Please contact Census Support with feature requests or bug reports.

## Support

- [GitHub Issues](https://github.com/sutrolabs/terraform-provider-census/issues)
- [Census Documentation](https://docs.getcensus.com/)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
