# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased - v0.3.0]

### Changed
- **[PLANNED]** Migrate `field_mapping` from TypeList to TypeSet for true order independence. This will cause a one-time diff showing all mappings as "replaced" during upgrade, but eliminates all future order-related drift. The migration is automatic and safe - mappings are keyed by destination field ("to") which is unique per sync. Users will see a large but harmless diff on first upgrade.

## [0.2.7] - 2025-12-16

### Added
- **429 Rate Limit Handling**: Automatic retry with exponential backoff for HTTP 429 (Too Many Requests) responses. The provider now intelligently handles rate limiting by:
  - Respecting `Retry-After` header when provided by the Census API (supports both delay-seconds and HTTP-date formats per RFC 7231)
  - Using exponential backoff with jitter when Retry-After is not present (1s → 2s → 4s → 8s → 16s → 32s → 64s → 90s max)
  - Continuing retries for up to 5 minutes (separate from the 1-minute per-request timeout)
  - Only retrying 429 errors (other errors like 5xx, network failures are not retried and fail immediately)
  - Providing debug and warning logs for retry attempts via Terraform's structured logging
  - Gracefully handling context cancellation and deadline exhaustion

## [0.2.6] - 2025-12-16

### Fixed
- **Sync Schedule Type Mismatch**: Fixed critical bug where `schedule_day` field type mismatch caused "Provider produced inconsistent result after apply" errors when creating syncs with cron-based schedules. The Census API returns `schedule_day` as a string (e.g., "Monday"), but the provider struct expected `*int`, causing JSON unmarshaling to silently fail. Changed `schedule_day` to `*string` to match the API contract.

## [0.2.5] - 2025-12-12

### Changed
- Increased HTTP client timeout from 30 seconds to 8 minutes to better handle high-latency scenarios, particularly for creating syncs, which can take a few minutes. This prevents premature client-side timeouts while Census API operations are still processing, reducing state drift and improving reliability when creating large numbers of resources.

## [0.2.4] - 2025-12-03

### Fixed
- **Sync Update API Format**: Fixed sync updates to use OpenAPI-compliant `mappings` field ([]MappingAttributes) instead of `field_mappings` ([]FieldMapping). This ensures sync updates are properly formatted for the Census API.
- **Liquid Template Mapping Format**: Fixed liquid_template mappings to send the template string directly to the API (`"{{ template }}"`) instead of incorrectly wrapping it in a hash (`{"liquid_template": "{{ template }}"}`). This resolves errors when creating or updating syncs with liquid template mappings.
- **Alert ID Preservation**: Fixed alert updates to preserve existing alert IDs, preventing API errors during sync updates.

### Added
- Comprehensive unit tests for mapping format conversions (9 new test cases)
- Debug logging to help diagnose sync update issues
- Export conversion functions (ConvertFieldMappingsToMappingAttributes, ConvertMappingAttributesToFieldMappings) for testability

## [0.2.3] - 2025-12-02

### Fixed
- Fix ghost alert entries appearing when removing alerts from syncs. This was caused by a Terraform SDK bug with TypeSet that creates zero-value entries during removal operations (N→N-1 or N→0). The fix migrates alert storage from TypeSet to TypeList, addressing the root cause. State is automatically compatible between versions - no manual migration needed.

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