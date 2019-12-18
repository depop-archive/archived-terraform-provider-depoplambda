package depoplambda

import (
	"fmt"
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// S3 should contain a manifest file with the following
// [{
//	name: "test_lambda"
//	schedule: "* * * * *"
// 	version: "0.1"
// }]
func TestAccLambdaManifests(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLambdaManifestsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceLambdaManifests("data.depoplambda_lambdas.lambdas"),
				),
			},
		},
	})
}

func testAccDataSourceLambdaManifests(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Printf("[DEBUG] name: %s\n", name)

		ds, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if ds.Primary.ID == "" {
			return fmt.Errorf("Manifest data source ID is not set")
		}

		return nil
	}
}

const testAccLambdaManifestsConfig = `
data "depoplambda_lambdas" "lambdas" {
	s3_bucket = "terraform-provider-depoplambda-test"
}`
