package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbExportInstanceSlowQueriesResource_basic -v
func TestAccTencentCloudCynosdbExportInstanceSlowQueriesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbExportInstanceSlowQueries,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_slow_queries.export_instance_slow_queries", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_export_instance_slow_queries.export_instance_slow_queries", "file_content"),
				),
			},
		},
	})
}

const testAccCynosdbExportInstanceSlowQueries = CommonCynosdb + `

resource "tencentcloud_cynosdb_export_instance_slow_queries" "export_instance_slow_queries" {
	instance_id = var.cynosdb_cluster_instance_id
}

`
