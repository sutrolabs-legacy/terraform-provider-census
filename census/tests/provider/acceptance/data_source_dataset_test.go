package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	provider_test "github.com/sutrolabs/terraform-provider-census/census/tests/provider"
)

func TestAccDataSourceDataset_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { provider_test.TestAccPreCheckIntegration(t) },
		Providers: provider_test.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatasetConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Verify data source attributes
					resource.TestCheckResourceAttr("data.census_dataset.test", "name", "Test Dataset - Data Source"),
					resource.TestCheckResourceAttr("data.census_dataset.test", "type", "sql"),
					resource.TestCheckResourceAttrSet("data.census_dataset.test", "id"),
					resource.TestCheckResourceAttrSet("data.census_dataset.test", "workspace_id"),
					resource.TestCheckResourceAttrSet("data.census_dataset.test", "source_id"),
					resource.TestCheckResourceAttrSet("data.census_dataset.test", "query"),

					// Verify data source matches resource
					resource.TestCheckResourceAttrPair(
						"data.census_dataset.test", "id",
						"census_dataset.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.census_dataset.test", "workspace_id",
						"census_dataset.test", "workspace_id",
					),
					resource.TestCheckResourceAttrPair(
						"data.census_dataset.test", "name",
						"census_dataset.test", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.census_dataset.test", "query",
						"census_dataset.test", "query",
					),
				),
			},
		},
	})
}

func testAccDataSourceDatasetConfig_basic() string {
	return fmt.Sprintf(`
resource "census_workspace" "test" {
  name = "Test Workspace - Dataset Data Source"
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

resource "census_dataset" "test" {
  workspace_id = census_workspace.test.id
  name         = "Test Dataset - Data Source"
  type         = "sql"
  description  = "Test dataset for data source testing"
  source_id    = census_source.test.id

  query = <<-SQL
    SELECT
      id,
      email,
      first_name,
      last_name
    FROM users
    WHERE active = true
  SQL
}

data "census_dataset" "test" {
  id           = census_dataset.test.id
  workspace_id = census_workspace.test.id
}
`,
		os.Getenv("CENSUS_TEST_REDSHIFT_HOST"),
		getEnvOrDefault("CENSUS_TEST_REDSHIFT_PORT", "5439"),
		os.Getenv("CENSUS_TEST_REDSHIFT_DATABASE"),
		os.Getenv("CENSUS_TEST_REDSHIFT_USERNAME"),
		os.Getenv("CENSUS_TEST_REDSHIFT_PASSWORD"),
	)
}
