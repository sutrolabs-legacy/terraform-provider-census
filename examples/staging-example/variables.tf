variable "census_personal_token" {
  description = "Census Personal Access Token for authentication"
  type        = string
  sensitive   = true
}

variable "census_base_url" {
  description = "Census API base URL (use staging URL for testing)"
  type        = string
  default     = "https://app.staging.getcensus.com/api/v1"
}

variable "workspace_name" {
  description = "Name for the Census workspace"
  type        = string
  default     = "Staging Test Workspace"
}

variable "notification_emails" {
  description = "List of email addresses to receive notifications"
  type        = list(string)
  default     = []
}

# Source Configuration
variable "source_label" {
  description = "Label for the data source"
  type        = string
  default     = "Staging Data Warehouse"
}

variable "source_type" {
  description = "Type of data source (e.g., postgres, snowflake, bigquery)"
  type        = string
}

variable "source_credentials" {
  description = "Credentials for the data source (supports booleans, numbers, and strings)"
  type        = map(any)
  sensitive   = true
}

variable "source_table" {
  description = "Source table name for the sync"
  type        = string
}

# Destination Configuration
variable "destination_label" {
  description = "Label for the destination"
  type        = string
  default     = "Staging CRM"
}

variable "destination_type" {
  description = "Type of destination (e.g., salesforce, hubspot)"
  type        = string
}

variable "destination_credentials" {
  description = "Credentials for the destination (supports booleans, numbers, and strings)"
  type        = map(any)
  sensitive   = true
}

variable "destination_object" {
  description = "Destination object (e.g., Salesforce object name)"
  type        = string
}

# Sync Configuration
variable "sync_label" {
  description = "Label for the sync"
  type        = string
  default     = "Staging Test Sync"
}

variable "field_mapping" {
  description = "Field mappings for the sync"
  type = list(object({
    from      = string
    to        = string
    operation = optional(string)
  }))
}

variable "sync_frequency" {
  description = "Sync frequency (hour, day, week, month)"
  type        = string
  default     = "day"
}

variable "sync_hour" {
  description = "Hour of day to run sync (0-23)"
  type        = number
  default     = 9
}

variable "sync_minute" {
  description = "Minute of hour to run sync (0-59)"
  type        = number
  default     = 0
}

variable "sync_paused" {
  description = "Whether the sync is paused"
  type        = bool
  default     = true
}