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
				Sensitive:   true, // Makes sure that the API key is not shown in the log.
			},
			"host": {
				Description: "Address of your Netbox endpoint, e.g. 127.0.0.1 or example.com",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   false,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"netbox_ip_address": resourceIPAddress(),
			//"netbox_prefix":       resourcePrefix(),
			"netbox_child_prefix": resourcePrefixChild(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"netbox_ip_address": dataSourceIPAddress(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		api_key: d.Get("api_key").(string),
		host:    d.Get("host").(string),
	}
	return config.Client()
}
