package acceptance

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	provider_test "github.com/sutrolabs/terraform-provider-census/census/tests/provider"
)

func TestAccResourceSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { provider_test.TestAccPreCheckIntegration(t) },
		Providers: provider_test.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourceConfig_redshift(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("census_source.test", "name", "Test Redshift Source"),
					resource.TestCheckResourceAttr("census_source.test", "type", "redshift"),
					resource.TestCheckResourceAttrSet("census_source.test", "id"),
					resource.TestCheckResourceAttrSet("census_source.test", "workspace_id"),
					resource.TestCheckResourceAttrSet("census_source.test", "created_at"),
				),
			},
		},
	})
}

func TestAccResourceSource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { provider_test.TestAccPreCheckIntegration(t) },
		Providers: provider_test.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourceConfig_redshift(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("census_source.test", "name", "Test Redshift Source"),
					resource.TestCheckResourceAttr("census_source.test", "auto_refresh_tables", "false"),
				),
			},
			{
				Config: testAccResourceSourceConfig_redshiftUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("census_source.test", "name", "Updated Redshift Source"),
					resource.TestCheckResourceAttr("census_source.test", "auto_refresh_tables", "true"),
				),
			},
		},
	})
}

func testAccResourceSourceConfig_redshift() string {
	return fmt.Sprintf(`
resource "census_workspace" "test" {
  name = "Test Workspace - Source"
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
`,
		os.Getenv("CENSUS_TEST_REDSHIFT_HOST"),
		getEnvOrDefault("CENSUS_TEST_REDSHIFT_PORT", "5439"),
		os.Getenv("CENSUS_TEST_REDSHIFT_DATABASE"),
		os.Getenv("CENSUS_TEST_REDSHIFT_USERNAME"),
		os.Getenv("CENSUS_TEST_REDSHIFT_PASSWORD"),
	)
}

func testAccResourceSourceConfig_redshiftUpdated() string {
	return fmt.Sprintf(`
resource "census_workspace" "test" {
  name = "Test Workspace - Source"
}

resource "census_source" "test" {
  workspace_id = census_workspace.test.id
  name = "Updated Redshift Source"
  type = "redshift"

  connection_config = {
    hostname = "%s"
    port     = "%s"
    database = "%s"
    user     = "%s"
    password = "%s"
  }

  auto_refresh_tables = true
}
`,
		os.Getenv("CENSUS_TEST_REDSHIFT_HOST"),
		getEnvOrDefault("CENSUS_TEST_REDSHIFT_PORT", "5439"),
		os.Getenv("CENSUS_TEST_REDSHIFT_DATABASE"),
		os.Getenv("CENSUS_TEST_REDSHIFT_USERNAME"),
		os.Getenv("CENSUS_TEST_REDSHIFT_PASSWORD"),
	)
}

func TestAccResourceSource_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { provider_test.TestAccPreCheckIntegration(t) },
		Providers: provider_test.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourceConfig_redshift(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("census_source.test", "name", "Test Redshift Source"),
					resource.TestCheckResourceAttr("census_source.test", "type", "redshift"),
				),
			},
			{
				ResourceName:      "census_source.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSourceImportStateIdFunc("census_source.test"),
				// connection_config contains sensitive data that isn't returned by the API
				ImportStateVerifyIgnore: []string{"connection_config", "auto_refresh_tables"},
			},
		},
	})
}

// Helper to construct composite ID for import (workspace_id:source_id)
func testAccSourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s:%s",
			rs.Primary.Attributes["workspace_id"],
			rs.Primary.ID), nil
	}
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
