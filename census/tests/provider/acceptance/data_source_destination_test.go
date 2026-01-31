package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	provider_test "github.com/sutrolabs/terraform-provider-census/census/tests/provider"
)

func TestAccDataSourceDestination_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { provider_test.TestAccPreCheckIntegration(t) },
		Providers: provider_test.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDestinationConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.census_destination.test", "name", "Test Salesforce Destination"),
					resource.TestCheckResourceAttr("data.census_destination.test", "type", "salesforce"),
					resource.TestCheckResourceAttrSet("data.census_destination.test", "id"),
					resource.TestCheckResourceAttrSet("data.census_destination.test", "created_at"),
					resource.TestCheckResourceAttrPair("data.census_destination.test", "id", "census_destination.test", "id"),
					resource.TestCheckResourceAttrPair("data.census_destination.test", "workspace_id", "census_destination.test", "workspace_id"),
				),
			},
		},
	})
}

func testAccDataSourceDestinationConfig_basic() string {
	return fmt.Sprintf(`
resource "census_workspace" "test" {
  name = "Test Workspace - Data Source"
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
}

data "census_destination" "test" {
  id           = census_destination.test.id
  workspace_id = census_workspace.test.id
}
`,
		os.Getenv("CENSUS_TEST_SALESFORCE_USERNAME"),
		os.Getenv("CENSUS_TEST_SALESFORCE_INSTANCE_URL"),
		os.Getenv("CENSUS_TEST_SALESFORCE_CLIENT_ID"),
		os.Getenv("CENSUS_TEST_SALESFORCE_JWT_SIGNING_KEY"),
		os.Getenv("CENSUS_TEST_SALESFORCE_DOMAIN"),
	)
}
