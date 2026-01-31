package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	provider_test "github.com/sutrolabs/terraform-provider-census/census/tests/provider"
)

func TestAccDataSourceSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { provider_test.TestAccPreCheckIntegration(t) },
		Providers: provider_test.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.census_source.test", "name", "Test Redshift Source"),
					resource.TestCheckResourceAttr("data.census_source.test", "type", "redshift"),
					resource.TestCheckResourceAttrSet("data.census_source.test", "id"),
					resource.TestCheckResourceAttrSet("data.census_source.test", "created_at"),
					resource.TestCheckResourceAttrPair("data.census_source.test", "id", "census_source.test", "id"),
					resource.TestCheckResourceAttrPair("data.census_source.test", "workspace_id", "census_source.test", "workspace_id"),
				),
			},
		},
	})
}

func testAccDataSourceSourceConfig_basic() string {
	return fmt.Sprintf(`
resource "census_workspace" "test" {
  name = "Test Workspace - Data Source"
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
}

data "census_source" "test" {
  id           = census_source.test.id
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
