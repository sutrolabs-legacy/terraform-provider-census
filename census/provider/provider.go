package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/sutrolabs/terraform-provider-census/census/client"
)

const (
	// DefaultRegion is the default Census region
	DefaultRegion = "us"
	// DefaultBaseURL is the default base URL for US region
	DefaultBaseURL = "https://app.getcensus.com/api/v1"
)

// Provider returns the Census Terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"personal_access_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CENSUS_PERSONAL_ACCESS_TOKEN", nil),
				Description: "Personal access token for Census APIs. Used for all operations including dynamic workspace token retrieval. Can also be set via CENSUS_PERSONAL_ACCESS_TOKEN environment variable.",
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DefaultRegion,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^(us|eu)$`),
					"region must be either 'us' or 'eu'",
				),
				Description: "Census region to use (us or eu). Defaults to 'us'.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CENSUS_BASE_URL", ""),
				Description: "Base URL for Census API. If not provided, will be determined based on region. Can also be set via CENSUS_BASE_URL environment variable.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"census_workspace":   resourceWorkspace(),
			"census_source":      resourceSource(),
			"census_destination": resourceDestination(),
			"census_sync":        ResourceSync(),
			"census_dataset":     resourceDataset(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"census_workspace":   dataSourceWorkspace(),
			"census_source":      dataSourceSource(),
			"census_destination": dataSourceDestination(),
			"census_sync":        dataSourceSync(),
			"census_dataset":     dataSourceDataset(),
		},
		ConfigureContextFunc: configure,
	}
}

func configure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	personalToken := d.Get("personal_access_token").(string)
	region := d.Get("region").(string)
	baseURL := d.Get("base_url").(string)

	// Validate that personal access token is provided
	if personalToken == "" {
		return nil, diag.Errorf("personal_access_token is required")
	}

	// Determine base URL if not explicitly provided
	if baseURL == "" {
		switch region {
		case "us":
			baseURL = "https://app.getcensus.com/api/v1"
		case "eu":
			baseURL = "https://app-eu.getcensus.com/api/v1"
		default:
			return nil, diag.Errorf("unsupported region: %s", region)
		}
	}

	config := &client.Config{
		PersonalAccessToken:  personalToken,
		WorkspaceAccessToken: "", // No longer used - workspace tokens are fetched dynamically
		BaseURL:              baseURL,
		Region:               region,
	}

	client, err := client.NewClient(config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, diags
}
