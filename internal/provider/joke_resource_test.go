package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJokeResource(t *testing.T) {
	resourceName := "montana_joke.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJokeResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "text", "one"),
					resource.TestCheckResourceAttr(resourceName, "id", "joke-id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"text"},
			},
			{
				Config: testAccJokeResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "text", "two"),
				),
			},
		},
	})
}

func testAccJokeResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "montana_joke" "test" {
  text = %[1]q
}
`, configurableAttribute)
}
