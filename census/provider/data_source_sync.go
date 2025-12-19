package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/sutrolabs/terraform-provider-census/census/client"
)

func dataSourceSync() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about a Census sync.",
		ReadContext: dataSourceSyncRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the sync.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace this sync belongs to.",
			},
			"label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name/label of the sync.",
			},
			"source_attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source-specific configuration (e.g., SQL query, table selection).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the source connection.",
						},
						"object": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source object configuration (JSON string for complex objects like table/dataset).",
						},
					},
				},
			},
			"destination_attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Destination-specific configuration (e.g., object, operation mode).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the destination connection.",
						},
						"object": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination object name.",
						},
					},
				},
			},
			"field_mapping": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Field mappings between source and destination.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source field name.",
						},
						"to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination field name.",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mapping operation (direct, constant, sync_metadata, segment_membership, liquid_template).",
						},
						"constant": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Constant value for constant operations.",
						},
						"is_primary_identifier": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this field is the primary identifier (sync key).",
						},
					},
				},
			},
			"paused": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the sync is paused.",
			},
			"run_mode": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Run mode configuration for the sync (live vs triggered with various trigger types).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mode type: 'live' for continuous syncing or 'triggered' for event-based syncing.",
						},
						"triggers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Trigger configurations (only for 'triggered' mode).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Schedule-based trigger configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"frequency": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Sync frequency: never, continuous, quarter_hourly, hourly, daily, weekly, or expression (for cron).",
												},
												"day": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Day of week (Sunday-Saturday, for weekly schedules).",
												},
												"hour": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Hour to run (0-24).",
												},
												"minute": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Minute to run (0-59).",
												},
												"cron_expression": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cron expression (only valid when frequency is 'expression').",
												},
											},
										},
									},
									"dbt_cloud": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "dbt Cloud job trigger configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"project_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "dbt Cloud project ID.",
												},
												"job_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "dbt Cloud job ID.",
												},
											},
										},
									},
									"fivetran": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Fivetran connector trigger configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"job_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Fivetran job ID.",
												},
												"job_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Fivetran job name.",
												},
											},
										},
									},
									"sync_sequence": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Sync dependency trigger configuration (triggers after another sync completes).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sync_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ID of the sync to trigger after.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			// Computed status fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the sync.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the sync was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the sync was last updated.",
			},
			"last_run_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the sync was last executed.",
			},
			"next_run_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the sync is scheduled to run next.",
			},
			"last_run_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the last sync run.",
			},
		},
	}
}

func dataSourceSyncRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.Client)

	syncID, err := strconv.Atoi(d.Get("id").(string))
	if err != nil {
		return diag.Errorf("invalid sync ID: %s", d.Get("id").(string))
	}

	workspaceId := d.Get("workspace_id").(string)
	workspaceIdInt, err := strconv.Atoi(workspaceId)
	if err != nil {
		return diag.Errorf("invalid workspace ID: %s", workspaceId)
	}

	workspaceToken, err := apiClient.GetWorkspaceAPIKey(ctx, workspaceIdInt)
	if err != nil {
		return diag.FromErr(err)
	}

	sync, err := apiClient.GetSyncWithToken(ctx, syncID, workspaceToken)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if sync is nil (API returned successfully but with nil data)
	if sync == nil {
		return diag.Errorf("sync not found: %d", syncID)
	}

	d.SetId(strconv.Itoa(sync.ID))
	d.Set("workspace_id", workspaceId)
	d.Set("label", sync.Label)
	d.Set("status", sync.Status)
	d.Set("paused", sync.Paused)
	d.Set("created_at", sync.CreatedAt.Format("2006-01-02T15:04:05Z"))
	d.Set("updated_at", sync.UpdatedAt.Format("2006-01-02T15:04:05Z"))

	if sync.LastRunAt != nil {
		d.Set("last_run_at", sync.LastRunAt.Format("2006-01-02T15:04:05Z"))
	}
	if sync.NextRunAt != nil {
		d.Set("next_run_at", sync.NextRunAt.Format("2006-01-02T15:04:05Z"))
	}
	if sync.LastRunID != nil {
		d.Set("last_run_id", *sync.LastRunID)
	}

	// Set complex attributes
	if err := d.Set("source_attributes", flattenSourceAttributesForDataSource(sync.SourceAttributes)); err != nil {
		return diag.Errorf("failed to set source_attributes: %v", err)
	}
	if err := d.Set("destination_attributes", flattenDestinationAttributesForDataSource(sync.DestinationAttributes)); err != nil {
		return diag.Errorf("failed to set destination_attributes: %v", err)
	}
	if err := d.Set("field_mapping", FlattenFieldMappings(sync.FieldMappings)); err != nil {
		return diag.Errorf("failed to set field_mapping: %v", err)
	}
	if sync.Mode != nil {
		if err := d.Set("run_mode", FlattenRunMode(sync.Mode)); err != nil {
			return diag.Errorf("failed to set run_mode: %v", err)
		}
	}

	return nil
}

// flattenSourceAttributesForDataSource converts API source_attributes map to Terraform TypeList format
func flattenSourceAttributesForDataSource(attrs map[string]interface{}) []interface{} {
	if attrs == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Convert connection_id to string
	if connID, ok := attrs["connection_id"]; ok {
		switch v := connID.(type) {
		case float64:
			result["connection_id"] = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			result["connection_id"] = strconv.Itoa(v)
		case string:
			result["connection_id"] = v
		}
	}

	// Serialize object as JSON string
	if obj, ok := attrs["object"]; ok {
		if objMap, isMap := obj.(map[string]interface{}); isMap {
			jsonBytes, err := json.Marshal(objMap)
			if err == nil {
				result["object"] = string(jsonBytes)
			}
		} else {
			result["object"] = fmt.Sprintf("%v", obj)
		}
	}

	return []interface{}{result}
}

// flattenDestinationAttributesForDataSource converts API destination_attributes map to Terraform TypeList format
func flattenDestinationAttributesForDataSource(attrs map[string]interface{}) []interface{} {
	if attrs == nil {
		return nil
	}

	result := make(map[string]interface{})

	// Convert connection_id to string
	if connID, ok := attrs["connection_id"]; ok {
		switch v := connID.(type) {
		case float64:
			result["connection_id"] = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			result["connection_id"] = strconv.Itoa(v)
		case string:
			result["connection_id"] = v
		}
	}

	// Set object as string
	if obj, ok := attrs["object"]; ok {
		result["object"] = fmt.Sprintf("%v", obj)
	}

	return []interface{}{result}
}
