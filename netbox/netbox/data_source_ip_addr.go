package netbox

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIPAddr() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
			},
		},
	}
}
