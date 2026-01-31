package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	provider_test "github.com/sutrolabs/terraform-provider-census/census/tests/provider"
)

func TestAccDataSourceWorkspace_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { provider_test.TestAccPreCheckIntegration(t) },
		Providers: provider_test.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkspaceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Verify data source attributes
					resource.TestCheckResourceAttr("data.census_workspace.test", "name", "Test Workspace - Data Source"),
					resource.TestCheckResourceAttrSet("data.census_workspace.test", "id"),
					resource.TestCheckResourceAttrSet("data.census_workspace.test", "organization_id"),
					resource.TestCheckResourceAttrSet("data.census_workspace.test", "created_at"),

					// Verify data source matches resource
					resource.TestCheckResourceAttrPair(
						"data.census_workspace.test", "id",
						"census_workspace.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.census_workspace.test", "name",
						"census_workspace.test", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.census_workspace.test", "organization_id",
						"census_workspace.test", "organization_id",
					),
				),
			},
		},
	})
}

func testAccDataSourceWorkspaceConfig_basic() string {
	return `
resource "census_workspace" "test" {
  name = "Test Workspace - Data Source"
}

data "census_workspace" "test" {
  id = census_workspace.test.id
}
`
}
