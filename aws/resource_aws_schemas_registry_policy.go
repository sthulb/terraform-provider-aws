package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsSchemasRegistryPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSchemasRegistryPolicyCreate,
		Read: resourceAwsSchemasRegistryPolicyRead,
		Update: resourceAwsSchemasRegistryPolicyUpdate,
		Delete: resourceAwsSchemasRegistryPolicyDelete,

		Schema: map[string]*schema.Schema{
			"policy": {
				Type: schema.TypeString,
				Required: true,
				ValidateFunc: validateIAMPolicyJson,
				DiffSuppressFunc: suppressEquivalentAwsPolicyDiffs,
			},
			"registry_name": {
				Type: schema.TypeString,
				Required: true,
			},
			"revision_id": {
				Type: schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAwsSchemasRegistryPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	registryName := d.Get("registry_name").(string)

	createOpts := &schemas.PutResourcePolicyInput{
		RegistryName: aws.String(registryName),
	}

	if v, ok := d.GetOk("revision_id"); ok {
		createOpts.SetRevisionId(v.(string))
	}

	if _, err := conn.PutResourcePolicy(createOpts); err != nil {
		return fmt.Errorf("error creating Schemas Registry Policy")
	}

	d.SetId(registryName)

	return resourceAwsSchemasRegistryPolicyRead(d, meta)
}

func resourceAwsSchemasRegistryPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsSchemasRegistryPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsSchemasRegistryPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}