package rightscale

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rightscale/rsc/cm15"
	"github.com/rightscale/rsc/rsapi"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccRightScaleServer_basic(t *testing.T) {
	t.Parallel()

	var (
		instanceName   = "terraform-test-instance-" + testString + "-" + acctest.RandString(10)
		serverName     = "terraform-test-server-" + testString + "-" + acctest.RandString(10)
		imageHref      = getTestImageFromEnv()
		typeHref       = getTestInstanceTypeFromEnv()
		cloudHref      = getTestCloudFromEnv()
		templateHref   = getTestTemplateFromEnv()
		deploymentHref = getTestDeploymentFromEnv()
		server         cm15.Server
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServer_basic(serverName, instanceName, cloudHref, imageHref, typeHref, deploymentHref, templateHref),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerExists("rightscale_server.test-server", &server),
				),
			},
		},
	})
}

func testAccCheckServerExists(n string, server *cm15.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		loc := getCMClient().ServerLocator(getHrefFromID(rs.Primary.ID))

		found, err := loc.Show(rsapi.APIParams{"view": "default"})
		if err != nil {
			return err
		}

		*server = *found

		return nil
	}
}

func testAccCheckServerDestroy(s *terraform.State) error {
	c := getCMClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rightscale_server" {
			continue
		}

		loc := c.ServerLocator(getHrefFromID(rs.Primary.ID))
		servers, err := loc.Index(nil)
		if err != nil {
			return fmt.Errorf("failed to check for existence of server: %s", err)
		}
		found := false
		self := strings.Split(rs.Primary.ID, ":")[1]
		for _, server := range servers {
			if string(server.Locator(c).Href) == self && server.State != "terminated" {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("server still exists")
		}
	}

	return nil
}

func testAccServer_basic(serverName string, instanceName string, cloudHref string, imageHref string, typeHref string, deploymentHref string, templateHref string) string {
	return fmt.Sprintf(`
resource "rightscale_server" "test-server" {
  name                   = %q
  deployment_href        = %q
  instance {
    cloud_href           = %q
    image_href           = %q
    instance_type_href   = %q
    name                 = %q
    server_template_href = %q
  }
}
  `, serverName, deploymentHref, cloudHref, imageHref, typeHref, instanceName, templateHref)
}
