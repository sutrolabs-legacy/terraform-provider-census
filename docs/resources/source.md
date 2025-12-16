# census_source Resource

Manages a Census data source connection. Sources connect to data warehouses like Snowflake, BigQuery, Postgres, and others.

## Example Usage

### Snowflake Source (Password Authentication)

```hcl
resource "census_source" "warehouse" {
  workspace_id = census_workspace.main.id
  name         = "Production Warehouse"
  type         = "snowflake"

  connection_config = {
    account        = "abc12345.us-east-1"
    warehouse      = "COMPUTE_WH"
    database       = "PRODUCTION"
    username       = "census_user"
    password       = var.snowflake_password
    role           = "CENSUS_ROLE"
  }
}
```

### Snowflake Source (Keypair Authentication)

```hcl
resource "census_source" "warehouse_keypair" {
  workspace_id = census_workspace.main.id
  name         = "Production Warehouse (Keypair)"
  type         = "snowflake"

  connection_config = {
    account                = "abc12345.us-east-1"
    warehouse              = "COMPUTE_WH"
    database               = "PRODUCTION"
    username               = "census_user"
    role                   = "CENSUS_ROLE"
    use_keypair            = true  # Boolean works directly!
    private_key_pkcs8      = var.snowflake_private_key
    private_key_passphrase = var.snowflake_key_passphrase  # Optional, omit if key is not encrypted
  }
}
```

### BigQuery Source

```hcl
resource "census_source" "bigquery" {
  workspace_id = census_workspace.main.id
  name         = "Analytics BigQuery"
  type         = "big_query"

  connection_config = {
    project_id = "my-gcp-project"
    dataset_id = "analytics"
    private_key = var.gcp_service_account_key
  }
}
```

### Postgres Source

```hcl
resource "census_source" "postgres" {
  workspace_id = census_workspace.main.id
  name         = "Production Database"
  type         = "postgres"

  connection_config = {
    host     = "postgres.example.com"
    port     = 5432  # Numbers work directly
    database = "production"
    username = "census"
    password = var.postgres_password
  }
}
```

## Argument Reference

* `workspace_id` - (Required, Forces new resource) The ID of the workspace this source belongs to.
* `name` - (Required) The name of the source.
* `type` - (Required, Forces new resource) The type of data source connector. Supported types include:
  - `snowflake`
  - `big_query`
  - `postgres`
  - `redshift`
  - `databricks`
  - `mysql`
  - And many more... (validated against Census API)
* `connection_config` - (Required, Sensitive) Map of credentials for connecting to the source. Supports strings, numbers, and booleans. The required fields vary by source type and are validated against the Census API schema.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the source.
* `status` - The current status of the source.
* `test_status` - The test status of the source connection.

## Import

Sources can be imported using the workspace ID and source ID separated by a colon:

```shell
terraform import census_source.warehouse "workspace_id:source_id"
```

For example:

```shell
terraform import census_source.warehouse "12345:67890"
```

## Notes

* The `connection_config` field is marked as sensitive and will not be displayed in Terraform output.
* Source types and required credential fields are validated against the Census API's `/source_types` endpoint.
* After creation, the provider automatically triggers a table refresh to discover available tables.