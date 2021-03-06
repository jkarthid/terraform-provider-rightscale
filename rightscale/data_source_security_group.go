package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

// Example:
//
// data "rightscale_security_group" "ssh" {
//   filter {
//     resource_uid = "sg-c31ee987"
//   }
//   cloud_href = ${data.rightscale_cloud.ec2_us_east_1.id}
// }

func dataSourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: resourceSecurityGroupRead,

		Schema: map[string]*schema.Schema{
			"cloud_href": {
				Type:        schema.TypeString,
				Description: "ID of the security group cloud",
				Required:    true,
				ForceNew:    true,
			},
			"filter": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "name of security group, uses partial match",
							Optional:    true,
							ForceNew:    true,
						},
						"deployment_href": {
							Type:        schema.TypeString,
							Description: "ID of deployment resource that owns security group",
							Optional:    true,
							ForceNew:    true,
						},
						"resource_uid": {
							Type:        schema.TypeString,
							Description: "cloud ID of security group, e.g. 'sg-c31ee987'",
							Optional:    true,
							ForceNew:    true,
						},
						"network_href": {
							Type:        schema.TypeString,
							Description: "ID of the security group network resource",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},

			// Read-only fields
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(rsc.Client)
	cloud := d.Get("cloud_href").(string)
	loc := &rsc.Locator{Namespace: "rs_cm", Href: cloud}

	res, err := client.List(loc, "security_groups", cmFilters(d))
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return nil
	}
	for k, v := range res[0].Fields {
		d.Set(k, v)
	}
	d.SetId(res[0].Locator.Href)
	return nil
}
