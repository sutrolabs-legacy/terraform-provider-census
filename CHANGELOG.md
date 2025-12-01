# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.2] - 2025-12-01

### Fixed
- Fix Snowflake source validation to support keypair authentication. The provider now correctly handles `show.unless` conditions in field metadata, allowing sources to be created with `use_keypair=true` without requiring a password field.

## [0.2.1] - 2025-11-25

### Fixed
- Remove client-side validation for primary identifiers in sync resources. The Census API now handles this validation, fixing issues with certain destination types like Google Sheets that don't require explicit primary identifier mappings. ([#1](https://github.com/sutrolabs/terraform-provider-census/issues/1))
- Fix sync `paused` attribute not updating correctly when changing from `true` to `false`. Removed `omitempty` JSON tag from boolean field to ensure the value is always sent to the API.

## [0.2.0] - 2025-10-23 - Initial Public Release

This is the first official release of the Census Terraform Provider on the [Terraform Registry](https://registry.terraform.io/providers/sutrolabs/census/latest).

### Provider Features

Complete Census data pipeline management from sources to syncs with infrastructure-as-code.

#### Resources

- **`census_workspace`** - Manage Census workspaces
  - Notification emails configuration
  - API key retrieval on creation
  - Full CRUD operations with import support

- **`census_source`** - Data warehouse connections
  - Support for all Census-supported databases (Snowflake, BigQuery, Postgres, Redshift, etc.)
  - Connection credential management with validation
  - Auto table refresh functionality

- **`census_destination`** - Business tool integrations
  - Support for all Census-supported destinations (Salesforce, HubSpot, etc.)
  - Dynamic connector type validation via Census API
  - Connection testing and credential management

- **`census_dataset`** - SQL datasets for data transformation
  - Multi-line SQL query support with heredoc syntax
  - Column schema discovery (computed fields)
  - Source connection reference and validation

- **`census_sync`** - Data syncs between sources and destinations
  - Field mapping configuration (direct, hash, constant operations)
  - Sync scheduling (hourly, daily, weekly, manual modes)
  - Sync mode support (upsert, append, mirror)
  - Support for all source types (table, dataset, model, topic, segment, cohort)

#### Data Sources

All resources have corresponding data sources for read-only operations: `census_workspace`, `census_source`, `census_destination`, `census_dataset`, `census_sync`

#### Authentication & Configuration

- **PAT-only authentication** with dynamic workspace token retrieval
- **Multi-region support**: US, EU, and AU regions with automatic endpoint configuration
- **Environment variable support**: `CENSUS_PERSONAL_ACCESS_TOKEN`, `CENSUS_REGION`, `CENSUS_BASE_URL`
- **Staging environment support**: Custom base URL configuration for testing

#### Import Support

- All resources support Terraform import
- Composite import format for workspace-scoped resources: `workspace_id:resource_id`
- Example: `terraform import census_source.example 69962:828`

### Getting Started

Install the provider from the Terraform Registry:

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

For detailed documentation and examples, visit the [Terraform Registry](https://registry.terraform.io/providers/sutrolabs/census/latest/docs).