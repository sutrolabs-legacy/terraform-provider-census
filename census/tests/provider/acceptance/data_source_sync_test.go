package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	provider_test "github.com/sutrolabs/terraform-provider-census/census/tests/provider"
)

func TestAccDataSourceSync_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { provider_test.TestAccPreCheckIntegration(t) },
		Providers: provider_test.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSyncConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Verify data source attributes
					resource.TestCheckResourceAttr("data.census_sync.test", "label", "Test Sync - Data Source"),
					resource.TestCheckResourceAttrSet("data.census_sync.test", "id"),
					resource.TestCheckResourceAttrSet("data.census_sync.test", "workspace_id"),
					resource.TestCheckResourceAttr("data.census_sync.test", "paused", "true"),

					// Verify data source matches resource
					resource.TestCheckResourceAttrPair(
						"data.census_sync.test", "id",
						"census_sync.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.census_sync.test", "workspace_id",
						"census_sync.test", "workspace_id",
					),
					resource.TestCheckResourceAttrPair(
						"data.census_sync.test", "label",
						"census_sync.test", "label",
					),
				),
			},
		},
	})
}

func testAccDataSourceSyncConfig_basic() string {
	return fmt.Sprintf(`
resource "census_workspace" "test" {
  name = "Test Workspace Sync DS Query"
}

resource "census_source" "test" {
  workspace_id = census_workspace.test.id
  name = "Test Redshift Source"
  type = "redshift"

  connection_config = {
    hostname = "%s"
    port     = "%s"
    database = "%s"
    user     = "%s"
    password = "%s"
  }

  auto_refresh_tables = false
}

resource "census_destination" "test" {
  workspace_id = census_workspace.test.id
  name = "Test Salesforce Destination"
  type = "salesforce"

  connection_config = {
    username        = "%s"
    instance_url    = "%s"
    client_id       = "%s"
    jwt_signing_key = "%s"
    domain          = "%s"
  }

  auto_refresh_objects = false
}

resource "census_sync" "test" {
  workspace_id = census_workspace.test.id
  label        = "Test Sync - Data Source"

  source_attributes {
    connection_id = census_source.test.id
    object {
      type          = "table"
      table_catalog = "dev"
      table_schema  = "public"
      table_name    = "users"
    }
  }

  destination_attributes {
    connection_id = census_destination.test.id
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

  field_mapping {
    from = "last_name"
    to   = "LastName"
  }

  field_mapping {
    from = "id"
    to   = "Census_ID__c"
  }

  run_mode {
    type = "triggered"
    triggers {
      schedule {
        frequency = "never"
      }
    }
  }

  alert {
    type                 = "FailureAlertConfiguration"
    send_for             = "first_time"
    should_send_recovery = true
  }

  paused = true
}

data "census_sync" "test" {
  id           = census_sync.test.id
  workspace_id = census_workspace.test.id
}
`,
		os.Getenv("CENSUS_TEST_REDSHIFT_HOST"),
		getEnvOrDefault("CENSUS_TEST_REDSHIFT_PORT", "5439"),
		os.Getenv("CENSUS_TEST_REDSHIFT_DATABASE"),
		os.Getenv("CENSUS_TEST_REDSHIFT_USERNAME"),
		os.Getenv("CENSUS_TEST_REDSHIFT_PASSWORD"),
		os.Getenv("CENSUS_TEST_SALESFORCE_USERNAME"),
		os.Getenv("CENSUS_TEST_SALESFORCE_INSTANCE_URL"),
		os.Getenv("CENSUS_TEST_SALESFORCE_CLIENT_ID"),
		os.Getenv("CENSUS_TEST_SALESFORCE_JWT_SIGNING_KEY"),
		os.Getenv("CENSUS_TEST_SALESFORCE_DOMAIN"),
	)
}
