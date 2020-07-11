package netbox

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"api_key": {
				Description: "API key of your Netbox user",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true, // Make sure that the API key is not shown in the log.
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"netbox_ip_addr": dataSourceIPAddr(),
		},
	}
}
