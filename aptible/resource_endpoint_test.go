package aptible

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/aptible/go-deploy/aptible"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceEndpoint_customDomain(t *testing.T) {
	appHandle := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAptibleEndpointCustomDomain(appHandle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aptible_app.test", "handle", appHandle),
					resource.TestCheckResourceAttr("aptible_app.test", "env_id", strconv.Itoa(testEnvironmentId)),
					resource.TestCheckResourceAttrSet("aptible_app.test", "app_id"),
					resource.TestCheckResourceAttrSet("aptible_app.test", "git_repo"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "env_id", strconv.Itoa(testEnvironmentId)),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "endpoint_type", "https"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "internal", "true"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "domain", "www.aptible-test-demo.fake"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "platform", "alb"),
					resource.TestCheckResourceAttrSet("aptible_endpoint.test", "endpoint_id"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "virtual_domain", "www.aptible-test-demo.fake"),
					resource.TestMatchResourceAttr("aptible_endpoint.test", "external_hostname", regexp.MustCompile(`elb.*\.aptible\.in`)),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "dns_validation_record", "_acme-challenge.www.aptible-test-demo.fake"),
					resource.TestMatchResourceAttr("aptible_endpoint.test", "dns_validation_value", regexp.MustCompile(`acme\.elb.*\.aptible\.in`)),
				),
			},
			{
				ResourceName:      "aptible_endpoint.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceEndpoint_app(t *testing.T) {
	appHandle := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAptibleEndpointApp(appHandle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aptible_app.test", "handle", appHandle),
					resource.TestCheckResourceAttr("aptible_app.test", "env_id", strconv.Itoa(testEnvironmentId)),
					resource.TestCheckResourceAttrSet("aptible_app.test", "app_id"),
					resource.TestCheckResourceAttrSet("aptible_app.test", "git_repo"),

					resource.TestCheckResourceAttr("aptible_endpoint.test", "env_id", strconv.Itoa(testEnvironmentId)),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "resource_type", "app"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "endpoint_type", "https"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "internal", "true"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "platform", "alb"),
					resource.TestCheckResourceAttrSet("aptible_endpoint.test", "endpoint_id"),
					resource.TestMatchResourceAttr("aptible_endpoint.test", "virtual_domain", regexp.MustCompile(`app-.*\.on-aptible\.com`)),
					resource.TestMatchResourceAttr("aptible_endpoint.test", "external_hostname", regexp.MustCompile(`elb.*\.aptible\.in`)),
					resource.TestCheckNoResourceAttr("aptible_endpoint.test", "dns_validation_record"),
					resource.TestCheckNoResourceAttr("aptible_endpoint.test", "dns_validation_value"),
				),
			},
			{
				ResourceName:      "aptible_endpoint.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceEndpoint_db(t *testing.T) {
	dbHandle := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAptibleEndpointDatabase(dbHandle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aptible_database.test", "handle", dbHandle),
					resource.TestCheckResourceAttr("aptible_database.test", "env_id", strconv.Itoa(testEnvironmentId)),
					resource.TestCheckResourceAttrSet("aptible_database.test", "database_id"),
					resource.TestCheckResourceAttrSet("aptible_database.test", "default_connection_url"),

					resource.TestCheckResourceAttr("aptible_endpoint.test", "env_id", strconv.Itoa(testEnvironmentId)),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "resource_type", "database"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "endpoint_type", "tcp"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "internal", "false"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "platform", "elb"),
					resource.TestCheckResourceAttrSet("aptible_endpoint.test", "endpoint_id"),
				),
			},
			{
				ResourceName:      "aptible_endpoint.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceEndpoint_updateIPWhitelist(t *testing.T) {
	appHandle := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAptibleEndpointApp(appHandle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aptible_app.test", "handle", appHandle),
					resource.TestCheckResourceAttr("aptible_app.test", "env_id", strconv.Itoa(testEnvironmentId)),
					resource.TestCheckResourceAttrSet("aptible_app.test", "app_id"),
					resource.TestCheckResourceAttrSet("aptible_app.test", "git_repo"),

					resource.TestCheckResourceAttr("aptible_endpoint.test", "env_id", strconv.Itoa(testEnvironmentId)),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "resource_type", "app"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "endpoint_type", "https"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "internal", "true"),
					resource.TestCheckResourceAttr("aptible_endpoint.test", "platform", "alb"),
					resource.TestCheckResourceAttrSet("aptible_endpoint.test", "endpoint_id"),
				),
			},
			{
				ResourceName:      "aptible_endpoint.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAptibleEndpointUpdateIPWhitelist(appHandle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aptible_endpoint.test", "ip_filtering.0", "1.1.1.1/32"),
				),
			},
		},
	})
}

func TestAccResourceEndpoint_expectError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAptibleEndpointInvalidResourceType(),
				ExpectError: regexp.MustCompile(`expected resource_type to be one of .*, got should-error`),
			},
			{
				Config:      testAccAptibleEndpointInvalidEndpointType(),
				ExpectError: regexp.MustCompile(`expected endpoint_type to be one of .*, got should-error`),
			},
			{
				Config:      testAccAptibleEndpointInvalidPlatform(),
				ExpectError: regexp.MustCompile(`expected platform to be one of .*, got should-error`),
			},
			{
				Config:      testAccAptibleEndpointInvalidDomain(),
				ExpectError: regexp.MustCompile(`managed endpoints must specify a domain`),
			},
		},
	})
}

func testAccCheckEndpointDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*aptible.Client)
	// Allow time for deprovision operation to complete.
	// TODO: Replace this by waiting on the actual operation

	//lintignore:R018
	time.Sleep(30 * time.Second)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aptible_endpoint" {
			continue
		}

		res_id, err := strconv.Atoi(rs.Primary.Attributes["resource_id"])
		if err != nil {
			return err
		}

		res_typ := rs.Primary.Attributes["resource_type"]
		if err != nil {
			return err
		}

		if res_typ == "app" {
			endpoint, err := client.GetApp(int64(res_id))
			log.Println("Deleted? ", endpoint.Deleted)
			if !endpoint.Deleted {
				return fmt.Errorf("App %v not removed", res_id)
			}

			if err != nil {
				return err
			}

		} else {
			endpoint, err := client.GetDatabase(int64(res_id))
			log.Println("Deleted? ", endpoint.Deleted)
			if !endpoint.Deleted {
				return fmt.Errorf("Database %v not removed", res_id)
			}

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func testAccAptibleEndpointCustomDomain(appHandle string) string {
	output := fmt.Sprintf(`
resource "aptible_app" "test" {
	env_id = %d
	handle = "%v"
	config = {
		"APTIBLE_DOCKER_IMAGE" = "nginx"
	}
	service {
		process_type = "cmd"
		container_memory_limit = 512
		container_count = 1
	}
}

resource "aptible_endpoint" "test" {
	env_id = %d
	resource_id = aptible_app.test.app_id
	resource_type = "app"
	process_type = "cmd"
	endpoint_type = "https"
	managed = true
	domain = "www.aptible-test-demo.fake"
	internal = true
	platform = "alb"
}`, testEnvironmentId, appHandle, testEnvironmentId)
	log.Println("HCL generated: ", output)
	return output
}

func testAccAptibleEndpointApp(appHandle string) string {
	output := fmt.Sprintf(`
resource "aptible_app" "test" {
	env_id = %d
	handle = "%v"
	config = {
		"APTIBLE_DOCKER_IMAGE" = "nginx"
	}
	service {
		process_type = "cmd"
		container_memory_limit = 512
		container_count = 1
	}
}

resource "aptible_endpoint" "test" {
	env_id = %d
	resource_id = aptible_app.test.app_id
	resource_type = "app"
	process_type = "cmd"
	endpoint_type = "https"
	default_domain = true
	internal = true
	platform = "alb"
}`, testEnvironmentId, appHandle, testEnvironmentId)
	log.Println("HCL generated: ", output)
	return output
}

func testAccAptibleEndpointDatabase(dbHandle string) string {
	output := fmt.Sprintf(`
resource "aptible_database" "test" {
	env_id = %d
	handle = "%v"
	database_type = "postgresql"
	container_size = 1024
	disk_size = 10
}

resource "aptible_endpoint" "test" {
	env_id = %d
	resource_id = aptible_database.test.database_id
	resource_type = "database"
	endpoint_type = "tcp"
	internal = false
	platform = "elb"
}`, testEnvironmentId, dbHandle, testEnvironmentId)
	log.Println("HCL generated: ", output)
	return output
}

func testAccAptibleEndpointUpdateIPWhitelist(appHandle string) string {
	output := fmt.Sprintf(`
resource "aptible_app" "test" {
	env_id = %d
	handle = "%v"
	config = {
		"APTIBLE_DOCKER_IMAGE" = "nginx"
	}
	service {
		process_type = "cmd"
		container_memory_limit = 512
		container_count = 1
	}
}

resource "aptible_endpoint" "test" {
	env_id = %d
	resource_id = aptible_app.test.app_id
	resource_type = "app"
	process_type = "cmd"
	endpoint_type = "https"
	default_domain = true
	internal = true
	platform = "alb"
	ip_filtering = [
		"1.1.1.1/32",
	]
}`, testEnvironmentId, appHandle, testEnvironmentId)
	log.Println("HCL generated: ", output)
	return output
}

func testAccAptibleEndpointInvalidResourceType() string {
	output := fmt.Sprintf(`
resource "aptible_endpoint" "test" {
	env_id = %d
	resource_id = 1
	resource_type = "should-error"
	}`, testEnvironmentId)
	log.Println("HCL generated: ", output)
	return output
}

func testAccAptibleEndpointInvalidEndpointType() string {
	output := fmt.Sprintf(`
resource "aptible_endpoint" "test" {
	env_id = %d
	resource_id = 1
	resource_type = "app"
	process_type = "cmd"
	default_domain = true
	endpoint_type = "should-error"
	}`, testEnvironmentId)
	log.Println("HCL generated: ", output)
	return output
}

func testAccAptibleEndpointInvalidPlatform() string {
	output := fmt.Sprintf(`
resource "aptible_endpoint" "test" {
	env_id = %d
	resource_id = 1
	resource_type = "app"
	process_type = "cmd"
	default_domain = true
	platform = "should-error"
	}`, testEnvironmentId)
	log.Println("HCL generated: ", output)
	return output
}

func testAccAptibleEndpointInvalidDomain() string {
	output := fmt.Sprintf(`
resource "aptible_endpoint" "test" {
	env_id = %d
	resource_id = 1
	resource_type = "app"
	process_type = "cmd"
	default_domain = false
	platform = "alb"
	managed = true
	domain = ""
	}`, testEnvironmentId)
	log.Println("HCL generated: ", output)
	return output
}
