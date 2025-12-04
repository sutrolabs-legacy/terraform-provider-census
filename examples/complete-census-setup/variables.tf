# ==============================================================================
# AUTHENTICATION VARIABLES
# ==============================================================================

variable "census_personal_token" {
  description = "Census personal access token for all operations"
  type        = string
  sensitive   = true
}

variable "census_region" {
  description = "Census region (us or eu)"
  type        = string
  default     = "us"

  validation {
    condition     = contains(["us", "eu"], var.census_region)
    error_message = "Region must be either 'us' or 'eu'."
  }
}

# ==============================================================================
# CONNECTION CONFIGURATION
# ==============================================================================

variable "postgres_warehouse_connection" {
  description = "Postgres data warehouse connection configuration (supports booleans, numbers, and strings)"
  type        = map(any)
  sensitive   = true
}

variable "redshift_warehouse_connection" {
  description = "Redshift data warehouse connection configuration (supports booleans, numbers, and strings)"
  type        = map(any)
  sensitive   = true
}

variable "snowflake_warehouse_connection" {
  description = "Snowflake data warehouse connection configuration (supports booleans, numbers, and strings)"
  type        = map(any)
  sensitive   = true
}

variable "salesforce_prod_connection" {
  description = "Salesforce production environment connection configuration (supports booleans, numbers, and strings)"
  type        = map(any)
  sensitive   = true
}

variable "salesforce_staging_connection" {
  description = "Salesforce staging environment connection configuration (supports booleans, numbers, and strings)"
  type        = map(any)
  sensitive   = true
}

# ==============================================================================
# OPTIONAL CONFIGURATION
# ==============================================================================

variable "enable_auto_refresh" {
  description = "Enable automatic metadata refresh for sources and destinations after creation/updates"
  type        = bool
  default     = true
}

