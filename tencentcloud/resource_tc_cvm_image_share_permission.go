/*
Provides a resource to create a cvm image_share_permission

Example Usage

```hcl
resource "tencentcloud_cvm_image_share_permission" "image_share_permission" {
  image_id = "img-xxxxxx"
  account_ids = ["xxxxxx"]
}
```

Import

cvm image_share_permission can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_image_share_permission.image_share_permission image_share_permission_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCvmImageSharePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmImageSharePermissionCreate,
		Update: resourceTencentCloudCvmImageSharePermissionUpdate,
		Read:   resourceTencentCloudCvmImageSharePermissionRead,
		Delete: resourceTencentCloudCvmImageSharePermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID such as `img-gvbnzy6f`. You can only specify an image in the NORMAL state.",
			},

			"account_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of account IDs with which an image is shared.",
			},
		},
	}
}

func resourceTencentCloudCvmImageSharePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_image_share_permission.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = cvm.NewModifyImageSharePermissionRequest()
		imageId string
	)
	if v, ok := d.GetOk("image_id"); ok {
		imageId = v.(string)
		request.ImageId = helper.String(imageId)
	}

	if v, ok := d.GetOk("account_ids"); ok {
		accountIdsSet := v.(*schema.Set).List()
		for i := range accountIdsSet {
			accountIds := accountIdsSet[i].(string)
			request.AccountIds = append(request.AccountIds, &accountIds)
		}
	}

	request.Permission = helper.String(IMAGE_SHARE_PERMISSION_SHARE)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyImageSharePermission(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm modifyImageSharePermission failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(imageId)

	return resourceTencentCloudCvmImageSharePermissionRead(d, meta)
}

func resourceTencentCloudCvmImageSharePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_image_share_permission.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	if d.HasChange("account_ids") {
		old, new := d.GetChange("account_ids")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		add := newSet.Difference(oldSet).List()
		remove := oldSet.Difference(newSet).List()
		if len(add) > 0 {
			addError := service.ModifyImageSharePermission(ctx, d.Id(), IMAGE_SHARE_PERMISSION_SHARE, helper.InterfacesStrings(add))
			if addError != nil {
				return addError
			}
		}
		if len(remove) > 0 {
			removeError := service.ModifyImageSharePermission(ctx, d.Id(), IMAGE_SHARE_PERMISSION_CANCEL, helper.InterfacesStrings(remove))
			if removeError != nil {
				return removeError
			}
		}
	}

	return resourceTencentCloudCvmImageSharePermissionRead(d, meta)
}

func resourceTencentCloudCvmImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_image_share_permission.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var sharePermissionSet []*cvm.SharePermission

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImageSharePermissionByFilter(ctx, map[string]interface{}{"ImageId": helper.String(d.Id())})
		if e != nil {
			return retryError(e)
		}
		sharePermissionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	accountIds := make([]string, 0)
	for _, sharePermission := range sharePermissionSet {
		accountIds = append(accountIds, *sharePermission.AccountId)
	}

	_ = d.Set("account_ids", accountIds)
	_ = d.Set("image_id", d.Id())
	return nil
}

func resourceTencentCloudCvmImageSharePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_image_share_permission.delete")()
	defer inconsistentCheck(d, meta)()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var sharePermissionSet []*cvm.SharePermission

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImageSharePermissionByFilter(ctx, map[string]interface{}{"ImageId": helper.String(d.Id())})
		if e != nil {
			return retryError(e)
		}
		sharePermissionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	accountIds := make([]string, 0)
	for _, sharePermission := range sharePermissionSet {
		accountIds = append(accountIds, *sharePermission.AccountId)
	}

	err = service.ModifyImageSharePermission(ctx, d.Id(), IMAGE_SHARE_PERMISSION_CANCEL, accountIds)
	if err != nil {
		return err
	}

	return nil
}
