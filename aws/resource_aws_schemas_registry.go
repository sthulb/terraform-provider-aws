package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	"log"
)

func resourceAwsSchemasRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSchemasRegistryCreate,
		Read:   resourceAwsSchemasRegistryRead,
		Update: resourceAwsSchemasRegistryUpdate,
		Delete: resourceAwsSchemasRegistryDelete,

		Schema: map[string]*schema.Schema{
			"registry_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registry_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAwsSchemasRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	registryName := d.Get("registry_name").(string)

	createOpts := &schemas.CreateRegistryInput{
		RegistryName: aws.String(registryName),
	}

	if v, ok := d.GetOk("description"); ok {
		createOpts.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		createOpts.SetTags(keyvaluetags.New(v.(map[string]interface{})).IgnoreAws().SchemasTags())
	}

	log.Printf("[DEBUG] Schemas Create config: #{*createOpts}")
	createResp, err := conn.CreateRegistry(createOpts)
	if err != nil {
		return fmt.Errorf("error creating Schemas Registry: %s", err)
	}

	d.SetId(aws.StringValue(createResp.RegistryArn))

	return resourceAwsSchemasRegistryRead(d, meta)
}

func resourceAwsSchemasRegistryRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	describeInput := &schemas.DescribeRegistryInput{
		RegistryName: aws.String(d.Get("registry_name").(string)),
	}

	registry, err := conn.DescribeRegistry(describeInput)
	if err != nil {
		if isAWSErr(err, "", "") {
			log.Printf("[INFO] unable to find Schema Registry, removing from state: #{d.Id()}")
			d.SetId("")
			return nil
		}

		return err
	}

	if err := d.Set("registry_arn", aws.StringValue(registry.RegistryArn)); err != nil {
		return err
	}

	if err := d.Set("registry_name", aws.StringValue(registry.RegistryName)); err != nil {
		return err
	}

	if err := d.Set("description", aws.StringValue(registry.Description)); err != nil {
		return err
	}

	tags := keyvaluetags.SchemasKeyValueTags(registry.Tags)
	if err := d.Set("tags", tags.IgnoreAws().Map()); err != nil {
		return err
	}

	return nil
}

func resourceAwsSchemasRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	registryName := d.Get("registry_name").(string)

	if d.HasChange("description") {
		modifyOpts := &schemas.UpdateRegistryInput{
			RegistryName: aws.String(d.Get("registry_name").(string)),
			Description:  aws.String(d.Get("description").(string)),
		}

		log.Printf("[DEBUG] Schemas Update Registry: #{*modifyOpts}")
		if _, err := conn.UpdateRegistry(modifyOpts); err != nil {
			return fmt.Errorf("error updating Schemas Registry (%s): %s", registryName, err)
		}
	}
	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})

		log.Printf("[DEBUG] Schemas Update Registry Tags: #{*tags}")

		if _, err := conn.TagResource(&schemas.TagResourceInput{
			ResourceArn: aws.String(d.Id()),
			Tags:        keyvaluetags.New(tags).IgnoreAws().SchemasTags(),
		}); err != nil {
			return fmt.Errorf("error updating Schemas Registry tags (%s): %s", registryName, err)
		}
	}

	return resourceAwsSagemakerEndpointRead(d, meta)
}

func resourceAwsSchemasRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	registryName := d.Get("registry_name").(string)

	deleteOpts := &schemas.DeleteRegistryInput{
		RegistryName: aws.String(registryName),
	}

	log.Printf("[DEBUG] Schemas Delete Registry: #{*deleteOpts}")
	if _, err := conn.DeleteRegistry(deleteOpts); err != nil {
		if isAWSErr(err, "ValidationException", "") {
			return nil
		}

		return fmt.Errorf("error deleting Schemas Registry (%s): %s", registryName, err)
	}

	return nil
}
