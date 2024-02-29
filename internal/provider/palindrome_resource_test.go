package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPalindromeResource(t *testing.T) {
	resourceName := "montana_palindrome.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPalindromeResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "configurable_attribute", "one"),
					resource.TestCheckResourceAttr(resourceName, "defaulted", "palindrome value when not configured"),
					resource.TestCheckResourceAttr(resourceName, "id", "palindrome-id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			{
				Config: testAccPalindromeResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "configurable_attribute", "two"),
				),
			},
		},
	})
}

func testAccPalindromeResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "montana_palindrome" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}
