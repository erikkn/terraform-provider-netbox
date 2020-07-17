package netbox

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
)

func dataSourceIPAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIPAddressRead,

		Schema: map[string]*schema.Schema{

			"cidr_block": {
				Description: "The CIDR block of the desired object.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"dns_name": {
				Description: "DNS Name associated with the IP address object.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"role_name": {
				Description: "The name of the role that is associated with the CIDR",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
			},

			"role_id": {
				Description: "The ID of the role that is associated with the CIDR",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
			},
		},
	}
}

// Function types; This function shares the same signature and is therefore of type ReadFunc.
func dataSourceIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).Ipam

	address := d.Get("cidr_block").(string)
	dnsName := d.Get("dns_name").(string)

	params := ipam.IpamIPAddressesListParams{
		Address: &address,
		DNSName: &dnsName,
		Context: context.Background(),
	}

	ipLookup, err := conn.IpamIPAddressesList(&params, nil)
	if err != nil {
		fmt.Errorf("[ERROR]: ", err)
	}

	for _, v := range ipLookup.Payload.Results {
		d.SetId("123") // !: Change needed!
		d.Set("cidr_block", *v.Address)
		d.Set("dns_name", v.DNSName)

		if v.Role == nil {
			log.Printf("[INFO]: Role returned nil value")
			d.Set("role_name", "")
			d.Set("role_id", "")
		} else {
			d.Set("role_name", *v.Role.Label)
			d.Set("role_id", *v.Role.Value)
		}

	}

	return nil
}