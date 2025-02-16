---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_renew_db_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_renew_db_instance"
description: |-
  Provides a resource to create a sqlserver renew_db_instance
---

# tencentcloud_sqlserver_renew_db_instance

Provides a resource to create a sqlserver renew_db_instance

## Example Usage

```hcl
resource "tencentcloud_sqlserver_renew_db_instance" "renew_db_instance" {
  instance_id = "mssql-i1z41iwd"
  period      = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `period` - (Optional, Int) How many months to renew, the value range is 1-48, the default is 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver renew_db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_renew_db_instance.renew_db_instance renew_db_instance_id
```

