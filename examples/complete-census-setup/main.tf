terraform {
  required_providers {
    census = {
      source  = "sutrolabs/census"
      version = "~> 0.2.0"
    }
  }
}

# Provider configuration
provider "census" {
  personal_access_token = var.census_personal_token # For workspace management
  region                = var.census_region

  # Note: No single workspace_access_token - workspace-specific tokens are handled
  # per resource using the workspace API keys from the workspaces configuration
}

# ==============================================================================
# PHASE 1: WORKSPACE MANAGEMENT (✅ Available)
# ==============================================================================

# Create workspaces as Terraform resources
resource "census_workspace" "marketing_prod" {
  name                     = "Marketing Production"
  notification_emails      = ["marketing-data@company.com", "data-alerts@company.com"]
  return_workspace_api_key = true
}

resource "census_workspace" "marketing_staging" {
  name                     = "Marketing Staging"
  notification_emails      = ["marketing-dev@company.com"]
  return_workspace_api_key = true
}

resource "census_workspace" "revops_prod" {
  name                     = "Revenue Operations Production"
  notification_emails      = ["revops@company.com", "finance-alerts@company.com"]
  return_workspace_api_key = true
}

# Example: Reference existing workspaces with data sources
# Uncomment and configure these if you have existing workspaces to reference
# data "census_workspace" "existing_sales" {
#   id = "12345"  # Your existing workspace ID
# }
# 
# data "census_workspace" "existing_finance" {
#   id = "67890"  # Another existing workspace ID
# }

# ==============================================================================
# PHASE 2: WORKSPACE-SCOPED DATA CONNECTIONS (🚧 Updated!)
# ==============================================================================

# Create sources explicitly tied to specific workspaces
# Each source belongs to exactly one workspace

# Postgres sources (original)
resource "census_source" "marketing_prod_postgres" {
  workspace_id = census_workspace.marketing_prod.id

  name = "Marketing Production Postgres Warehouse"
  type = "postgres"

  connection_config   = var.postgres_warehouse_connection
  auto_refresh_tables = var.enable_auto_refresh
}

resource "census_source" "revops_prod_postgres" {
  workspace_id = census_workspace.revops_prod.id

  name = "RevOps Production Postgres Warehouse"
  type = "postgres"

  connection_config   = var.postgres_warehouse_connection
  auto_refresh_tables = var.enable_auto_refresh
}

# Add staging warehouse sources for the staging sync

# Postgres staging source (original)
resource "census_source" "marketing_staging_postgres" {
  workspace_id = census_workspace.marketing_staging.id

  name = "Marketing Staging Postgres Warehouse"
  type = "postgres"

  connection_config   = var.postgres_warehouse_connection
  auto_refresh_tables = var.enable_auto_refresh
}

# Redshift staging source (working connection)
resource "census_source" "marketing_staging_warehouse" {
  workspace_id = census_workspace.marketing_staging.id

  name = "Marketing Staging Redshift Warehouse"
  type = "redshift"

  connection_config   = var.redshift_warehouse_connection
  auto_refresh_tables = var.enable_auto_refresh
}

# Redshift sources (working connection)
resource "census_source" "marketing_prod_warehouse" {
  workspace_id = census_workspace.marketing_prod.id

  name = "Marketing Production Redshift Warehouse"
  type = "redshift"

  connection_config   = var.redshift_warehouse_connection
  auto_refresh_tables = var.enable_auto_refresh
}

resource "census_source" "revops_prod_warehouse" {
  workspace_id = census_workspace.revops_prod.id

  name = "RevOps Production Redshift Warehouse"
  type = "redshift"

  connection_config   = var.redshift_warehouse_connection
  auto_refresh_tables = var.enable_auto_refresh
}

# Create destinations explicitly tied to specific workspaces
resource "census_destination" "marketing_prod_crm" {
  workspace_id = census_workspace.marketing_prod.id

  name = "Marketing Production CRM"
  type = "salesforce"

  connection_config    = var.salesforce_prod_connection
  auto_refresh_objects = var.enable_auto_refresh
}

resource "census_destination" "marketing_staging_crm" {
  workspace_id = census_workspace.marketing_staging.id

  name = "Marketing Staging CRM"
  type = "salesforce"

  connection_config    = var.salesforce_staging_connection
  auto_refresh_objects = var.enable_auto_refresh
}

# Example: Reference sources/destinations from existing workspaces
# Uncomment if you have existing workspaces with existing sources/destinations
# data "census_source" "existing_sales_warehouse" {
#   id = "98765"  # ID of existing source in existing workspace
# }
#
# data "census_destination" "existing_finance_system" {
#   id = "54321"  # ID of existing destination in existing workspace
# }

# ==============================================================================
# PHASE 2.5: SQL DATASETS (✅ Available!)
# ==============================================================================

# Create SQL datasets for data transformation and modeling
# Datasets allow you to define custom SQL queries against your sources

resource "census_dataset" "active_users" {
  workspace_id = census_workspace.marketing_prod.id
  name         = "Active Users"
  type         = "sql"
  description  = "Filtered dataset of active users for marketing campaigns"

  source_id = census_source.marketing_prod_warehouse.id

  query = <<-SQL
    SELECT
      id,
      email,
      first_name,
      last_name,
      created_at,
      last_login_at
    FROM users
    WHERE active = true
      AND last_login_at > CURRENT_DATE - INTERVAL '30 days'
  SQL
}

resource "census_dataset" "high_value_customers" {
  workspace_id = census_workspace.marketing_prod.id
  name         = "High Value Customers"
  type         = "sql"
  description  = "Customers with lifetime value > $1000"

  source_id = census_source.marketing_prod_warehouse.id

  query = <<-SQL
    SELECT
      u.id,
      u.email,
      u.first_name,
      u.last_name,
      SUM(o.amount) as lifetime_value,
      COUNT(o.id) as order_count
    FROM users u
    JOIN orders o ON u.id = o.user_id
    GROUP BY u.id, u.email, u.first_name, u.last_name
    HAVING SUM(o.amount) > 1000
  SQL
}

resource "census_dataset" "all_users" {
  workspace_id = census_workspace.marketing_prod.id
  name         = "All Users"
  type         = "sql"
  description  = "Simple dataset with all user data for syncing"

  source_id = census_source.marketing_prod_warehouse.id

  query = <<-SQL
    SELECT * FROM users
  SQL
}

# Example: Use data source to reference an existing dataset
# data "census_dataset" "existing_dataset" {
#   id           = "123"
#   workspace_id = census_workspace.marketing_prod.id
# }

# ==============================================================================
# PHASE 3: DATA SYNCS (✅ Available!)
# ==============================================================================

# Create syncs that connect sources to destinations within the same workspace
# Each sync defines how data flows from a source to a destination

resource "census_sync" "marketing_contact_sync" {
  workspace_id = census_workspace.marketing_prod.id
  label        = "REAL Marketing Contact Sync"

  # Source configuration - use a table from the warehouse
  source_attributes {
    connection_id = census_source.marketing_prod_warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "public"
      table_catalog = "dev" # Use database name as catalog
    }
  }

  # Destination configuration - specify Salesforce object
  destination_attributes {
    connection_id = census_destination.marketing_prod_crm.id
    object        = "Contact"
  }

  # Operation mode for the sync
  operation = "upsert"

  # Field mappings between source and destination
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
    from = "id"
    to   = "Census_ID__c"
  }

  field_mapping {
    from = "last_name"
    to   = "LastName"
  }

  # Run mode - triggered on daily schedule at 6 AM UTC
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

  # Alert configurations for production sync monitoring
  alert {
    type                 = "FailureAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options              = {}
  }

  alert {
    type                 = "InvalidRecordPercentAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options = {
      threshold = "75"
    }
  }

  # Start paused until ready to go live
  paused = true
}

resource "census_sync" "marketing_contact_sync_2" {
  workspace_id = census_workspace.marketing_prod.id
  label        = "REAL Marketing Contact Sync 2"

  # Source configuration - use a table from the warehouse
  source_attributes {
    connection_id = census_source.marketing_prod_warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "public"
      table_catalog = "dev" # Use database name as catalog
    }
  }

  # Destination configuration - specify Salesforce object
  destination_attributes {
    connection_id = census_destination.marketing_prod_crm.id
    object        = "Contact"
  }

  # Operation mode for the sync
  operation = "upsert"

  alert {
    type                 = "FailureAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options              = {}
  }

  alert {
    type                 = "InvalidRecordPercentAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options = {
      threshold = "75"
    }
  }


  # Field mappings between source and destination
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
    from = "id"
    to   = "Census_ID__c"
  }

  field_mapping {
    from = "last_name"
    to   = "LastName"
  }

  # Run mode - triggered on daily schedule at 6 AM UTC
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

  # Start paused until ready to go live
  paused = true
}

resource "census_sync" "marketing_contact_sync_3" {
  workspace_id = census_workspace.marketing_prod.id
  label        = "Marketing Contact Sync 3"

  # Source configuration - use a table from the warehouse
  source_attributes {
    connection_id = census_source.marketing_prod_warehouse.id
    object {
      type          = "table"
      table_name    = "users"
      table_schema  = "public"
      table_catalog = "dev" # Use database name as catalog
    }
  }

  # Destination configuration - specify Salesforce object
  destination_attributes {
    connection_id = census_destination.marketing_prod_crm.id
    object        = "Contact"
  }

  # Operation mode for the sync
  operation = "upsert"

  # Field mappings between source and destination
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
    from = "id"
    to   = "Census_ID__c"
  }


  field_mapping {
    from = "last_name"
    to   = "LastName"
  }
  alert {
    type                 = "FailureAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options              = {}
  }

  alert {
    type                 = "InvalidRecordPercentAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options = {
      threshold = "75"
    }
  }


  # Run mode - triggered on daily schedule at 6 AM UTC
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

  # Start paused until ready to go live
  paused = true
}

# Sync using a dataset as the source
resource "census_sync" "dataset_contact_sync" {
  workspace_id = census_workspace.marketing_prod.id
  label        = "Dataset to Contacts Sync"

  # Source configuration - use a dataset instead of table
  source_attributes {
    connection_id = census_source.marketing_prod_warehouse.id
    object {
      type = "dataset"
      id   = census_dataset.all_users.id
    }
  }

  # Destination configuration - Salesforce Contacts
  destination_attributes {
    connection_id = census_destination.marketing_prod_crm.id
    object        = "Contact"
  }

  operation = "upsert"

  advanced_configuration = jsonencode({
    bulk_id_lookup = true
  })

  # Field mappings using dataset columns

  field_mapping {
    from = "last_name"
    to   = "LastName"
  }

  # Sync metadata mapping - track sync run IDs
  field_mapping {
    type              = "sync_metadata"
    sync_metadata_key = "sync_run_id"
    to                = "another__c"
  }

  # Segment membership mapping
  field_mapping {
    type                = "segment_membership"
    segment_identify_by = "name"
    to                  = "MailingStreet"
  }

  # Liquid template transformation
  field_mapping {
    type            = "liquid_template"
    liquid_template = "{{ record['_fivetran_synced'] | upcase }}"
    to              = "OtherPostalCode"
  }

  field_mapping {
    from = "user_id"
    to   = "Census_ID__c"
  }

  field_mapping {
    type     = "constant"
    constant = "HERE is my constant value"
    to       = "AssistantName"
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

  # Run mode - triggered hourly at minute 10
  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "hourly"
        minute    = 10
      }
    }
  }

  # Comprehensive alerting for high-frequency sync
  alert {
    type                 = "FailureAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options              = {}
  }

  alert {
    type                 = "InvalidRecordPercentAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
    options = {
      threshold = "75"
    }
  }

  alert {
    type                 = "RuntimeAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = false
    options = {
      threshold  = "45"
      unit       = "minutes"
      start_type = "actual"
    }
  }

  alert {
    type                 = "StatusAlertConfiguration"
    send_for             = "every_time"
    should_send_recovery = false
    options = {
      status_name = "completed"
    }
  }
}

# resource "census_sync" "marketing_staging_test_sync" {
#   workspace_id = census_workspace.marketing_staging.id
#   label        = "Marketing Staging Test Sync"


#   # Test sync with limited data  
#   source_attributes {
#     connection_id = census_source.marketing_staging_warehouse.id
#     object {
#       type          = "table"
#       table_name    = "users"
#       table_schema  = "public"
#       table_catalog = "dev" # Use database name as catalog
#     }
#   }

#   destination_attributes {
#     connection_id = census_destination.marketing_staging_crm.id
#     object        = "Contact"
#   }

#   # Operation mode for the sync
#   operation = "upsert"

#   field_mapping {
#     from                  = "email"
#     to                    = "Email"
#     is_primary_identifier = true
#   }

#   field_mapping {
#     from = "first_name"
#     to   = "FirstName"
#   }

#   field_mapping {
#     from = "last_name"
#     to   = "LastName"
#   }


#   # Manual sync for testing
#   run_mode {
#     type = "triggered"
#     triggers {
#       schedule {
#         frequency = "manual"
#       }
#     }
#   }

#   paused = false # Ready for testing
# }

# # RevOps sync for account data
# resource "census_sync" "revops_account_sync" {
#   workspace_id = census_workspace.revops_prod.id
#   label        = "RevOps Account Sync"


#   # Source configuration - account enrichment data
#   source_attributes {
#     connection_id = census_source.revops_prod_warehouse.id
#     object {
#       type          = "table"
#       table_name    = "users"
#       table_schema  = "public"
#       table_catalog = "dev" # Use database name as catalog
#     }
#   }

#   # Destination configuration
#   destination_attributes {
#     connection_id = census_destination.marketing_prod_crm.id
#     object        = "Account"
#   }

#   # Operation mode for the sync
#   operation = "upsert"

#   # Field mappings for account data
#   field_mapping {
#     from                  = "company_domain"
#     to                    = "Website"
#     is_primary_identifier = true
#   }

#   field_mapping {
#     from = "company_name"
#     to   = "Name"
#   }

#   field_mapping {
#     from = "industry"
#     to   = "Industry"
#   }

#   field_mapping {
#     from = "employee_count"
#     to   = "NumberOfEmployees"
#   }

#   field_mapping {
#     from = "annual_revenue"
#     to   = "AnnualRevenue"
#   }

#   field_mapping {
#     from = "country"
#     to   = "BillingCountry"
#   }


#   # Run weekly on Sundays at 2 AM
#   run_mode {
#     type = "triggered"
#     triggers {
#       schedule {
#         frequency = "weekly"
#         day       = "Sunday"
#         hour      = 2
#         minute    = 0
#       }
#     }
#   }

#   paused = true # Start paused for review
# }

# Example of using data source to reference an existing sync
# data "census_sync" "existing_sync" {
#   id           = "123456"
#   workspace_id = census_workspace.marketing_prod.id
# }

# ==============================================================================
# PHASE 4: ADVANCED FEATURES (📈 Future updates)
# ==============================================================================

# Dataset definitions, sync runs, webhooks, etc. will be added here
# as we implement more resources

# ==============================================================================
# OUTPUTS - Information about created resources
# ==============================================================================

output "workspaces_info" {
  description = "Information about all created workspaces"
  value = {
    marketing_prod = {
      id              = census_workspace.marketing_prod.id
      name            = census_workspace.marketing_prod.name
      organization_id = census_workspace.marketing_prod.organization_id
      created_at      = census_workspace.marketing_prod.created_at
      api_key         = census_workspace.marketing_prod.api_key
    }
    marketing_staging = {
      id              = census_workspace.marketing_staging.id
      name            = census_workspace.marketing_staging.name
      organization_id = census_workspace.marketing_staging.organization_id
      created_at      = census_workspace.marketing_staging.created_at
      api_key         = census_workspace.marketing_staging.api_key
    }
    revops_prod = {
      id              = census_workspace.revops_prod.id
      name            = census_workspace.revops_prod.name
      organization_id = census_workspace.revops_prod.organization_id
      created_at      = census_workspace.revops_prod.created_at
      api_key         = census_workspace.revops_prod.api_key
    }
  }
  sensitive = true
}

output "sources_info" {
  description = "Information about all created sources"
  value = {
    marketing_prod_warehouse = {
      id           = census_source.marketing_prod_warehouse.id
      name         = census_source.marketing_prod_warehouse.name
      type         = census_source.marketing_prod_warehouse.type
      status       = census_source.marketing_prod_warehouse.status
      created_at   = census_source.marketing_prod_warehouse.created_at
      workspace_id = census_source.marketing_prod_warehouse.workspace_id
    }
    revops_prod_warehouse = {
      id           = census_source.revops_prod_warehouse.id
      name         = census_source.revops_prod_warehouse.name
      type         = census_source.revops_prod_warehouse.type
      status       = census_source.revops_prod_warehouse.status
      created_at   = census_source.revops_prod_warehouse.created_at
      workspace_id = census_source.revops_prod_warehouse.workspace_id
    }
  }
}

output "destinations_info" {
  description = "Information about all created destinations"
  value = {
    marketing_prod_crm = {
      id           = census_destination.marketing_prod_crm.id
      name         = census_destination.marketing_prod_crm.name
      type         = census_destination.marketing_prod_crm.type
      status       = census_destination.marketing_prod_crm.status
      created_at   = census_destination.marketing_prod_crm.created_at
      workspace_id = census_destination.marketing_prod_crm.workspace_id
    }
    marketing_staging_crm = {
      id           = census_destination.marketing_staging_crm.id
      name         = census_destination.marketing_staging_crm.name
      type         = census_destination.marketing_staging_crm.type
      status       = census_destination.marketing_staging_crm.status
      created_at   = census_destination.marketing_staging_crm.created_at
      workspace_id = census_destination.marketing_staging_crm.workspace_id
    }
  }
}

output "setup_summary" {
  description = "Summary of the complete workspace-scoped Census setup"
  sensitive   = true
  value = {
    workspaces_ready = {
      marketing_prod    = census_workspace.marketing_prod.id != ""
      marketing_staging = census_workspace.marketing_staging.id != ""
      revops_prod       = census_workspace.revops_prod.id != ""
    }
    sources_ready = {
      marketing_prod_warehouse = census_source.marketing_prod_warehouse.id != ""
      revops_prod_warehouse    = census_source.revops_prod_warehouse.id != ""
    }
    destinations_ready = {
      marketing_prod_crm    = census_destination.marketing_prod_crm.id != ""
      marketing_staging_crm = census_destination.marketing_staging_crm.id != ""
    }

    total_workspaces   = 3 # marketing_prod + marketing_staging + revops_prod
    total_sources      = 2 # marketing_prod_warehouse + revops_prod_warehouse
    total_destinations = 2 # marketing_prod_crm + marketing_staging_crm
    total_resources    = 7 # 3 workspaces + 2 sources + 2 destinations

    workspace_structure = [
      "marketing_prod: ${census_workspace.marketing_prod.name} (ID: ${census_workspace.marketing_prod.id})",
      "marketing_staging: ${census_workspace.marketing_staging.name} (ID: ${census_workspace.marketing_staging.id})",
      "revops_prod: ${census_workspace.revops_prod.name} (ID: ${census_workspace.revops_prod.id})"
    ]

    resource_relationships = [
      "✅ Each source/destination explicitly belongs to one workspace",
      "✅ No ambiguous global configurations",
      "✅ Clear workspace isolation and security boundaries",
      "✅ Proper Terraform resource relationships"
    ]

    authentication_model = [
      "Personal Access Token (PAT): Used for workspace management operations",
      "Each workspace has its own API key stored in Terraform state",
      "Sources and destinations are scoped to specific workspaces",
      "Use data sources to reference existing workspaces/sources/destinations"
    ]

    terraform_patterns = [
      "✅ Resources for managing Census objects with Terraform",
      "✅ Data sources for referencing existing Census objects",
      "✅ Explicit workspace_id parameters for clear relationships",
      "✅ Standard Terraform import patterns for existing resources"
    ]

    next_steps = [
      "Verify all workspace connections in Census dashboard",
      "Use 'terraform import' to bring existing workspaces under Terraform management",
      "Add data sources to reference existing sources/destinations if needed",
      "Configure workspace-specific sync rules (coming in next update)"
    ]
  }
}
