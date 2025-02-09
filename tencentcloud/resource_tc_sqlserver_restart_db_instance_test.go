package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverRestartDBInstanceResource_basic -v
func TestAccTencentCloudSqlserverRestartDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverRestartDBInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_restart_db_instance.restart_db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_restart_db_instance.restart_db_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverRestartDBInstance = testAccSqlserverBasicInstanceNetwork + `
resource "tencentcloud_sqlserver_instance" "test" {
  name                          = "tf_sqlserver_instance"
  availability_zone             = var.default_az
  charge_type                   = "POSTPAID_BY_HOUR"
  vpc_id                        = local.vpc_id
  subnet_id                     = local.subnet_id
  security_groups               = [local.sg_id]
  project_id                    = 0
  memory                        = 2
  storage                       = 10
  maintenance_week_set          = [1,2,3]
  maintenance_start_time        = "09:00"
  maintenance_time_span         = 3
  tags = {
    "test"                      = "test"
  }
}

resource "tencentcloud_sqlserver_restart_db_instance" "restart_db_instance" {
  instance_id = tencentcloud_sqlserver_instance.test.id
}
`
