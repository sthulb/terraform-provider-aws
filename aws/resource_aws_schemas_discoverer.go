package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	"log"
)

func resourceAwsSchemasDiscoverer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSchemasDiscovererCreate,
		Read:   resourceAwsSchemasDiscovererRead,
		Update: resourceAwsSchemasDiscovererUpdate,
		Delete: resourceAwsSchemasDiscovererDelete,

		Schema: map[string]*schema.Schema{
			"discovery_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discovery_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"source_arn": {
				Type:     schema.TypeString,
				Computed: true,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAwsSchemasDiscovererCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	createOpts := &schemas.CreateDiscovererInput{}
	if v, ok := d.GetOk("description"); ok {
		createOpts.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		createOpts.SetTags(keyvaluetags.New(v.(map[string]interface{})).IgnoreAws().SchemasTags())
	}

	log.Printf("[DEBUG] Schemas Create Discoverer: #{*createOpts}")
	createResp, err := conn.CreateDiscoverer(createOpts)
	if err != nil {
		return fmt.Errorf("error creating Schemas Discoverer: #{err}")
	}

	d.SetId(aws.StringValue(createResp.DiscovererId))

	return resourceAwsSchemasDiscovererRead(d, meta)
}

func resourceAwsSchemasDiscovererRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	describeOpts := &schemas.DescribeDiscovererInput{
		DiscovererId: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Schemas Describe Discoverer: #{*describeOpts}")
	discoverer, err := conn.DescribeDiscoverer(describeOpts)
	if err != nil {
		if isAWSErr(err, "", "") {
			log.Printf("[INFO] unable to find Schemas Discoverer, removing from state: #{d.Id()}")
			d.SetId("")
			return nil
		}

		return err
	}

	if err := d.Set("discovery_arn", aws.StringValue(discoverer.DiscovererArn)); err != nil {
		return err
	}

	if err := d.Set("discovery_id", aws.StringValue(discoverer.DiscovererId)); err != nil {
		return err
	}

	if err := d.Set("description", aws.StringValue(discoverer.Description)); err != nil {
		return err
	}

	if err := d.Set("source_arn", aws.StringValue(discoverer.SourceArn)); err != nil {
		return err
	}

	tags := keyvaluetags.SchemasKeyValueTags(discoverer.Tags)
	if err := d.Set("tags", tags.IgnoreAws().Map()); err != nil {
		return err
	}

	return nil
}

func resourceAwsSchemasDiscovererUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	discoveryId := d.Id()

	modifyOpts := &schemas.UpdateDiscovererInput{
		DiscovererId: aws.String(discoveryId),
		Description:  aws.String(d.Get("description").(string)),
	}

	log.Printf("[DEBUG] Schemas Describe Discoverer: #{*describeOpts}")

	if _, err := conn.UpdateDiscoverer(modifyOpts); err != nil {
		return fmt.Errorf("error updating Schemas Discoverer (%s): %s", discoveryId, err)
	}

	return resourceAwsSchemasRegistryRead(d, meta)
}

func resourceAwsSchemasDiscovererDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).schemasconn

	deleteOpts := &schemas.DeleteDiscovererInput{
		DiscovererId: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Schemas Delete Discoverer: #{*deleteOpts}")
	if _, err := conn.DeleteDiscoverer(deleteOpts); err != nil {
		if isAWSErr(err, "ValidationException", "") {
			return nil
		}

		return fmt.Errorf("error deleting Schemas Discoverer (%s): %s", d.Id(), err)
	}

	return nil
}
