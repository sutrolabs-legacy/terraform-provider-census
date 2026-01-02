# census_sync Resource

Manages a Census sync that moves data from a source (table, dataset, model, etc.) to a destination (Salesforce, HubSpot, etc.) with configurable field mappings and scheduling.

## Example Usage

### Basic Sync with Field Mappings

```hcl
resource "census_sync" "user_sync" {
  workspace_id = census_workspace.main.id
  label        = "Users to Salesforce"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Contact"
  }

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    from = "first_name"
    to   = "FirstName"
  }

  field_mapping {
    from = "last_name"
    to   = "LastName"
  }

  operation = "upsert"

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "hourly"
        minute    = 0
      }
    }
  }
}
```

### Sync with Dataset Source

```hcl
resource "census_sync" "high_value_sync" {
  workspace_id = census_workspace.main.id
  label        = "High Value Customers to HubSpot"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type = "dataset"
      id   = census_dataset.high_value_customers.id
    }
  }

  destination_attributes {
    connection_id = census_destination.hubspot.id
    object        = "contacts"
  }

  field_mapping {
    from                  = "email"
    to                    = "email"
    is_primary_identifier = true
  }

  field_mapping {
    from = "lifetime_value"
    to   = "lifetime_value"
  }

  operation = "upsert"

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

### Sync with Segment Source

```hcl
resource "census_sync" "vip_segment_sync" {
  workspace_id = census_workspace.main.id
  label        = "VIP Users Segment to Salesforce"

  source_attributes {
    connection_id = 829
    object {
      type       = "segment"
      id         = "3060"      # The segment ID
      dataset_id = "5951"      # The dataset ID that the segment belongs to
    }
  }

  destination_attributes {
    connection_id = 456
    object        = "Contact"
  }

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    from = "first_name"
    to   = "FirstName"
  }

  field_mapping {
    from = "last_name"
    to   = "LastName"
  }

  operation = "upsert"

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "daily"
        hour      = 2
        minute    = 0
      }
    }
  }
}
```

### Sync with Constant Value

```hcl
resource "census_sync" "tagged_sync" {
  workspace_id = census_workspace.main.id
  label        = "Tagged Contact Sync"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Contact"
  }

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    type     = "constant"
    constant = "Terraform Managed"
    to       = "LeadSource"
  }

  operation = "upsert"
}
```

### Sync with Sync Metadata Mapping

```hcl
resource "census_sync" "metadata_sync" {
  workspace_id = census_workspace.main.id
  label        = "Sync with Metadata Tracking"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Contact"
  }

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    from = "first_name"
    to   = "FirstName"
  }

  field_mapping {
    # Map Census sync_run_id to a custom field
    type              = "sync_metadata"
    sync_metadata_key = "sync_run_id"
    to                = "Last_Sync_Run_ID__c"
  }

  operation = "upsert"
}
```

### Sync with Segment Membership

```hcl
resource "census_sync" "segment_sync" {
  workspace_id = census_workspace.main.id
  label        = "Sync with Segment Data"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Contact"
  }

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    # Map segment membership information
    type                = "segment_membership"
    segment_identify_by = "name"
    to                  = "Active_Segments__c"
  }

  operation = "upsert"
}
```

### Sync with Liquid Template Transformation

```hcl
resource "census_sync" "template_sync" {
  workspace_id = census_workspace.main.id
  label        = "Sync with Field Transformations"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Contact"
  }

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    from = "first_name"
    to   = "FirstName"
  }

  field_mapping {
    # Use Liquid template to transform data
    type            = "liquid_template"
    liquid_template = "{{ record['status'] | upcase }}"
    to              = "Account_Status__c"
  }

  operation = "upsert"
}
```

### Sync with Automatic Field Mapping (Sync All Properties)

```hcl
resource "census_sync" "auto_sync" {
  workspace_id = census_workspace.main.id
  label        = "Auto-Mapped Users Sync"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Contact"
  }

  # Automatically sync all properties from source to destination
  field_behavior      = "sync_all_properties"
  field_normalization = "snake_case"  # Format field names in snake_case
  field_order         = "mapping_order"

  # Only need to define the primary identifier when using sync_all_properties
  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  operation = "upsert"

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

### Sync to Google Sheets (Mirror with Auto-Mapping)

```hcl
resource "census_sync" "sheets_sync" {
  workspace_id = census_workspace.main.id
  label        = "Data Export to Google Sheets"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.google_sheets.id
    # Google Sheets object format: JSON with spreadsheet_id and sheet_id
    object        = "{\"spreadsheet_id\":\"1r9HavWIo-CS14sblFl-Kxa2kVbMGlqWIBDU6Hb_W8sA\",\"sheet_id\":0}"
  }

  operation = "mirror"

  # Automatically sync all properties - no primary identifier needed for Google Sheets mirror
  field_behavior = "sync_all_properties"
  field_normalization = "match_source_names"

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "daily"
        hour      = 6
        minute    = 0
      }
    }
  }
}
```

### Sync with Lookup Field (Foreign Key Relationship)

```hcl
resource "census_sync" "user_list_sync" {
  workspace_id = census_workspace.main.id
  label        = "Users to Google Ads Customer Match"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.google_ads.id
    object        = "user_data"
  }

  field_mapping {
    from                  = "email"
    to                    = "user_identifier.hashed_email"
    is_primary_identifier = true
  }

  field_mapping {
    # Map a constant value to user_list_id via lookup
    # This looks up the user_list record where id = "6600827417"
    type          = "constant"
    constant      = "6600827417"
    to            = "user_list_id"
    lookup_object = "user_list"
    lookup_field  = "id"
  }

  operation = "mirror"

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "hourly"
        minute    = 10
      }
    }
  }
}
```

### Sync with Advanced Configuration (File Export)

```hcl
resource "census_sync" "blob_storage_sync" {
  workspace_id = census_workspace.main.id
  label        = "Users to Azure Blob Storage"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type = "model"
      id   = "21130"
    }
  }

  destination_attributes {
    connection_id = census_destination.azure_blob.id
    object        = "path_to_file/data_%m-%d-%y.parquet"
  }

  field_mapping {
    from = "email"
    to   = "EMAIL"
  }

  operation           = "mirror"
  field_behavior      = "sync_all_properties"
  field_normalization = "match_source_names"

  # Advanced configuration for file export
  advanced_configuration = jsonencode({
    file_format        = "Parquet"
    csv_delimiter      = ","
    csv_include_header = true
  })

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "hourly"
        minute    = 35
      }
    }
  }
}
```

### Sync with Alert Configurations

```hcl
resource "census_sync" "monitored_sync" {
  workspace_id = census_workspace.main.id
  label        = "High-Priority Customer Sync with Alerts"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Contact"
  }

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    from = "name"
    to   = "Name"
  }

  operation = "upsert"

  # Configure multiple alerts
  alert {
    # Alert when sync fails completely
    type                 = "FailureAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options              = {}
  }

  alert {
    # Alert when more than 50% of records are invalid
    type                 = "InvalidRecordPercentAlertConfiguration"
    send_for             = "every_time"
    should_send_recovery = true
    options = {
      threshold = "50"
    }
  }

  alert {
    # Alert when sync runtime exceeds 30 minutes
    type                 = "RuntimeAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = false
    options = {
      threshold  = "30"
      unit       = "minutes"
      start_type = "actual"
    }
  }

  alert {
    # Alert on sync completion
    type                 = "StatusAlertConfiguration"
    send_for             = "every_time"
    should_send_recovery = false
    options = {
      status_name = "completed"
    }
  }

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "hourly"
        minute    = 0
      }
    }
  }
}
```

### Mirror Sync (Replace All)

```hcl
resource "census_sync" "mirror_sync" {
  workspace_id = census_workspace.main.id
  label        = "Product Catalog Mirror"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Product2"
  }

  field_mapping {
    from                  = "product_id"
    to                    = "ProductCode"
    is_primary_identifier = true
  }

  field_mapping {
    from = "name"
    to   = "Name"
  }

  operation = "mirror"

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "daily"
        hour      = 2
        minute    = 0
      }
    }
  }
}
```

### Append Sync with High Water Mark

```hcl
resource "census_sync" "incremental_append" {
  workspace_id = census_workspace.main.id
  label        = "Incremental Event Log Sync"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Event__c"
  }

  field_mapping {
    from                  = "event_id"
    to                    = "Event_ID__c"
    is_primary_identifier = true
  }

  field_mapping {
    from = "event_name"
    to   = "Name"
  }

  field_mapping {
    from = "updated_at"
    to   = "Updated_At__c"
  }

  operation = "append"

  # Use high water mark to only sync new records based on timestamp
  # This is more efficient than Census's default diff engine for append operations
  high_water_mark_attribute = "updated_at"

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "hourly"
        minute    = 15
      }
    }
  }
}
```

### Sync with Field Preservation and Null Value Control

```hcl
resource "census_sync" "preserve_example" {
  workspace_id = census_workspace.main.id
  label        = "Customer Sync with Field Preservation"

  source_attributes {
    connection_id = census_source.warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "demo"
      table_catalog = "dev"
    }
  }

  destination_attributes {
    connection_id = census_destination.salesforce.id
    object        = "Contact"
  }

  field_mapping {
    from                  = "email"
    to                    = "Email"
    is_primary_identifier = true
  }

  field_mapping {
    from = "first_name"
    to   = "FirstName"
  }

  field_mapping {
    # Don't overwrite existing phone numbers in destination
    from             = "phone"
    to               = "Phone"
    preserve_values  = true # Don't overwrite existing data in the destination
    sync_null_values = false  # Don't sync null phone values
  }

  field_mapping {
    # Generate a custom field in the destination
    from           = "customer_tier"
    to             = "Customer_Tier__c"
    generate_field = true
  }

  operation = "upsert"

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "daily"
        hour      = 9
        minute    = 0
      }
    }
  }
}
```

## Argument Reference

* `workspace_id` - (Required, Forces new resource) The ID of the workspace this sync belongs to.
* `label` - (Required) The label of the sync.
* `source_attributes` - (Required) Configuration block for the source. Block contains:
  * `connection_id` - (Required) The source connection ID
  * `object` - (Required) Object configuration block:
    * `type` - (Required) Source type: `"table"`, `"dataset"`, `"model"`, `"topic"`, `"segment"`, or `"cohort"`
    * For table sources: `table_name`, `table_schema`, and `table_catalog`
    * For dataset/model sources: `id` of the dataset/model
    * For segment sources: use `type="segment"`, provide the segment `id`, and specify `dataset_id` for the dataset the segment belongs to
    * For cohort sources: use `type="cohort"`, provide the cohort `id`, and specify `dataset_id` for the dataset the cohort belongs to
* `destination_attributes` - (Required) Destination configuration block:
  * `connection_id` - (Required) The destination connection ID
  * `object` - (Required) The destination object name (e.g., "Contact" for Salesforce, "contacts" for HubSpot)
  * `lead_union_insert_to` - (Optional) Where to insert a union object (for Salesforce connections only)
* `field_mapping` - (Optional) Field mappings between source and destination. Define multiple `field_mapping` blocks for multiple mappings. Each mapping block includes:
  * `from` - Source field name (required for `type="direct"`). Omit for `constant`, `sync_metadata`, `segment_membership`, and `liquid_template` mappings.
  * `to` - Destination field name (required)
  * `type` - Mapping type: `"direct"` (default), `"constant"`, `"sync_metadata"`, `"segment_membership"`, or `"liquid_template"`.
  * `constant` - Constant value (must also set `type="constant"`)
  * `sync_metadata_key` - Sync metadata key (e.g., `"sync_run_id"`). Must also set `type="sync_metadata"`.
  * `segment_identify_by` - How to identify segments (e.g., `"name"`). Must also set `type="segment_membership"`.
  * `liquid_template` - Liquid template for data transformation (e.g., `"{{ record['field'] | upcase }}"`). Must also set `type="liquid_template"`.
  * `is_primary_identifier` - (Optional) Boolean indicating if this field is the primary identifier for matching records. Exactly one field_mapping must have this set to `true`. Defaults to `false`.
  * `lookup_object` - (Optional) Object to lookup for relationship mapping (e.g., `"user_list"`). Used with `lookup_field` for foreign key lookups.
  * `lookup_field` - (Optional) Field to use for lookup in the `lookup_object` (e.g., `"id"`). Used with `lookup_object` for foreign key lookups.
  * `preserve_values` - (Optional) If true, preserves existing values in the destination field and prevents Census from overwriting them. Defaults to `false`.
  * `generate_field` - (Optional) If true, Census will generate/create this field in the destination. Defaults to `false`.
  * `sync_null_values` - (Optional) If true (default), null values in the source will be synced to the destination. Set to false to skip syncing null values. Defaults to `true`.
  * `array_field` - (Optional) Whether the destination field is an array type. Only applicable when `generate_field` is true (for user-defined fields). Defaults to `false`.
  * `field_type` - (Optional) The type of the destination field. Only applicable when `generate_field` is true (for user-defined fields). Available types depend on the destination (e.g., "text", "number", "boolean", "date").
  * `follow_source_type` - (Optional) Whether the destination field type should automatically follow changes to the source column type. Defaults to `false`.
* `operation` - (Optional) Sync mode: `"upsert"`, `"append"`, or `"mirror"`. Defaults to `"upsert"`.
* `field_behavior` - (Optional) Controls how fields are synced:
  * `"specific_properties"` (default) - Use only the field mappings defined in `field_mapping`
  * `"sync_all_properties"` - Automatically sync all properties from source to destination
* `field_normalization` - (Optional) When `field_behavior` is `"sync_all_properties"`, specifies how automatic field names should be normalized:
  * `"start_case"` - Start Case (e.g., "First Name")
  * `"lower_case"` - lower case (e.g., "first name")
  * `"upper_case"` - UPPER CASE (e.g., "FIRST NAME")
  * `"camel_case"` - camelCase (e.g., "firstName")
  * `"snake_case"` - snake_case (e.g., "first_name")
  * `"match_source_names"` - Use exact source field names
* `field_order` - (Optional) Specifies how destination fields should be ordered. Only applicable for destinations that support field ordering:
  * `"alphabetical_column_name"` (default) - Sort fields alphabetically
  * `"mapping_order"` - Use the order fields are defined in `field_mapping`
* `sync_behavior_family` - (Optional) Specifies the behavior family for the sync:
  * `"activateEvents"` - For event-based activation syncs (only supported for live syncs from Kafka/streaming sources)
  * `"mapRecords"` - For record mapping syncs (not supported for live syncs from Materialize)
* `advanced_configuration` - (Optional) Advanced configuration options specific to the destination type as JSON string. Use `jsonencode()` to specify values. Available options vary by destination (e.g., file format for file exports, bulk settings for APIs). Values can be strings, numbers, or booleans. Refer to destination-specific Census documentation for available options.
* `high_water_mark_attribute` - (Optional) The name of the timestamp column to use for high water mark diffing strategy. When set, append syncs will use this column to identify new records instead of the default Census diff engine (using primary keys). This is more efficient for append operations with timestamp-based data. Example: `"updated_at"`.
* `historical_sync_operation` - (Optional) Specifies how the first sync should handle historical records when using append operation. Only applicable for append syncs:
  * `"skip_current_records"` - Skip existing records on first sync, only sync new records going forward
  * `"backfill_all_records"` - Include all existing records on first sync (full backfill)
* `mirror_strategy` - (Optional, Computed) Specifies the strategy for mirror syncs. Only applicable when `operation` is set to `"mirror"`. Determines how Census keeps the destination in sync with the source data:
  * `"sync_updates_and_deletes"` - Incrementally syncs changes by inserting new records, updating modified records, and deleting records that no longer exist in the source. This is the most common and efficient strategy for keeping destinations in sync (default).
  * `"sync_updates_and_nulls"` - Updates existing records and sets fields to null when the source contains null values, without performing deletes.
  * `"upload_and_swap"` - Replaces the entire destination table with the current source snapshot. Useful for destinations that don't support incremental updates or when you need a complete refresh.
* `alert` - (Optional) Alert configurations for monitoring sync health. Define multiple `alert` blocks to configure multiple alerts. If no `alert` blocks are specified, the sync will be created with no alerts. Each alert block includes:
  * `type` - (Required) Type of alert. Valid values:
    * `"FailureAlertConfiguration"` - Alert when sync fails completely
    * `"InvalidRecordPercentAlertConfiguration"` - Alert when invalid/rejected records exceed threshold
    * `"FullSyncTriggerAlertConfiguration"` - Alert when a full sync is triggered
    * `"RecordCountDeviationAlertConfiguration"` - Alert when record counts deviate from expected
    * `"RuntimeAlertConfiguration"` - Alert when sync runtime exceeds threshold
    * `"StatusAlertConfiguration"` - Alert on sync status changes (started, completed)
  * `send_for` - (Optional) When to send alerts: `"first_time"` (default, only first violation) or `"every_time"` (every violation)
  * `should_send_recovery` - (Optional) Whether to send recovery notification when condition resolves. Defaults to `true`.
  * `options` - (Optional) Alert-specific configuration options (values as strings):
    * For `InvalidRecordPercentAlertConfiguration`:
      * `threshold` - Percentage (0-100) of invalid records that triggers alert
    * For `RecordCountDeviationAlertConfiguration`:
      * `threshold` - Percentage (0-100) deviation from expected count
      * `record_type` - Type to monitor: `source_record_count`, `records_updates`, `records_deletes`, `records_invalid`, `records_processed`, `records_updated`, or `records_failed`
    * For `RuntimeAlertConfiguration`:
      * `threshold` - Number of time units before alert
      * `unit` - Time unit: `"minutes"` or `"hours"`
      * `start_type` - When to start measuring: `"actual"` (when sync actually starts) or `"scheduled"` (from scheduled time)
    * For `StatusAlertConfiguration`:
      * `status_name` - Status to alert on: `"started"` or `"completed"`
  * `id` - (Computed) The alert configuration ID assigned by Census
* `run_mode` - (Optional) Run mode configuration block for controlling how and when the sync runs:
  * `type` - (Required) Mode type:
    * `"live"` - Continuous syncing for streaming sources (Kafka, Materialize)
    * `"triggered"` - Event-based syncing with configured triggers
  * `triggers` - (Optional) Trigger configurations (only for `triggered` mode). Multiple triggers can be configured simultaneously:
    * `schedule` - (Optional) Schedule-based trigger configuration block:
      * `frequency` - (Required) How often to run: `"never"`, `"continuous"`, `"quarter_hourly"`, `"hourly"`, `"daily"`, `"weekly"`, or `"expression"` (for cron)
      * `day` - (Optional) Day of week for weekly schedules: `"Sunday"`, `"Monday"`, `"Tuesday"`, `"Wednesday"`, `"Thursday"`, `"Friday"`, or `"Saturday"`
      * `hour` - (Optional) Hour to run (0-24) for daily/weekly schedules
      * `minute` - (Optional) Minute to run (0-59)
      * `cron_expression` - (Optional) Cron expression when `frequency` is `"expression"`. Mutually exclusive with hour/day settings
    * `dbt_cloud` - (Optional) dbt Cloud job trigger configuration block:
      * `project_id` - (Required) dbt Cloud project ID
      * `job_id` - (Required) dbt Cloud job ID
    * `fivetran` - (Optional) Fivetran connector trigger configuration block:
      * `job_id` - (Required) Fivetran job ID
      * `job_name` - (Required) Fivetran job name
    * `sync_sequence` - (Optional) Sync dependency trigger configuration block (triggers after another sync completes):
      * `sync_id` - (Required) ID of the sync to trigger after
## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the sync.
* `paused` - Whether the sync is currently paused.
* `status` - The current status of the sync.

## Import

Syncs can be imported using the workspace ID and sync ID separated by a colon:

```shell
terraform import census_sync.user_sync "workspace_id:sync_id"
```

For example:

```shell
terraform import census_sync.user_sync "12345:67890"
```

## Notes

* Field mappings are defined as multiple `field_mapping` blocks. The order of field mappings is preserved.
* The `source_attributes` and `destination_attributes` are configuration blocks, not JSON-encoded strings.
* Sync operations:
  * `upsert` - Insert new records and update existing ones
  * `append` - Only insert new records, never update
  * `mirror` - Replace all destination records with source data
* Manual syncs (frequency="never") must be triggered externally.
* Source types determine which fields are required in `source_attributes.object`.
  * For segment sources, use `type="segment"` with `id` (segment ID) and `dataset_id` (parent dataset ID)
  * For cohort sources, use `type="cohort"` with `id` (cohort ID) and `dataset_id` (parent dataset ID)