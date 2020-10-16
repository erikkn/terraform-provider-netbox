package netbox

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPAddressCreate,
		Read:   resourceIPAddressRead,
		Update: resourceIPAddressUpdate,
		Delete: resourceIPAddressDelete,

		Schema: map[string]*schema.Schema{
			"address": {
				Description: "(Required) The address that you would like to create, e.g. 192.168.96.0/20.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"description": {
				Description: "(Optional) The address's description",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"dns_name": {
				Description: "(Optional) String value of the DNS name you want to associate with this address",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},

			"net_family": {
				Description: "(Optional) Net Family number (e.g. IPv4 or IPv6) of your address. Default is IPv4",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "IPv4",
			},

			// TODO: Make this required?
			"tenant_group": {
				Description: "(Required) The tenant group organizes the individual tenants. An example of such group is `payment partners` and the tenant is a specific partner, e.g. `example` (see the `tenant` attribute).",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},

			// TODO: Make this required?
			"tenant": {
				Description: "(Required) The tenant field is used to signify the ownership of the address. Typically, tenants are used to represent internal departments, partners, AWS account, etc.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},

			"role": {
				Description: "(Optional) The functional role of this address; Roles are used to indicate some special attribute to the IP address. Valid values are `Loopback`, `Secondary`, `Anycast`, `VIP`, `VRRP`, `HSRP`, `GLBP`",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},

			"tags": {
				Description: "(Optional) A map of tags to assign to the resource.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client)
	ipamClient := conn.Ipam
	extrasClient := conn.Extras

	address := d.Get("address").(string)
	description := d.Get("description").(string)
	dnsName := d.Get("dns_name").(string)
	netFamily := d.Get("net_family").(string)
	role := d.Get("role").(string)
	rawTags := d.Get("tags").(*schema.Set).List()

	tags, err := createTag(&rawTags, extrasClient)
	if err != nil {
		fmt.Errorf("[ERROR] error creating tags: %s", err)
	}

	params := ipam.NewIpamIPAddressesCreateParams()
	params.WithData(&models.WritableIPAddress{
		Address:     &address,
		Description: description,
		DNSName:     dnsName,
		Family:      netFamily,
		Role:        role,
		Tags:        tags,
	})

	ip, err := ipamClient.IpamIPAddressesCreate(params, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating the IP address: %s", err)
	}
	log.Printf("[INFO]: Creation of the IP address is successful!")

	// `ID` returns an int64, while `SetId` requires a string as input.
	d.SetId(strconv.FormatInt(ip.Payload.ID, 10))

	return nil
}

// TODO: Double check the READ function, do I need the `HasChange` method?
func resourceIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).Ipam

	ipID := d.Id()

	params := ipam.NewIpamIPAddressesListParams().WithID(&ipID)
	fetch, err := conn.IpamIPAddressesList(params, nil)
	if err != nil {
		log.Printf("[WARN] IP address not found, removing from state")
		d.SetId("")
		return nil
	}
	log.Printf("[INFO]: Fetching IP address successful")

	// TODO: Add tags to the d.Set
	for _, v := range fetch.Payload.Results {
		d.Set("address", *v.Address)
		d.Set("description", v.Description)
		d.Set("dns_name", v.DNSName)
		d.Set("net_family", v.Family)
		d.Set("role", v.Role)
	}

	return nil
}

func resourceIPAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client)
	ipamClient := conn.Ipam
	extrasClient := conn.Extras

	ipID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("[ERROR] error while converting ID '%d' with error message: %s", ipID, err)
	}

	address := d.Get("address").(string)
	description := d.Get("description").(string)
	dnsName := d.Get("dns_name").(string)
	role := d.Get("role").(string)
	rawTags := d.Get("tags").(*schema.Set).List()

	tags, err := createTag(&rawTags, extrasClient)
	if err != nil {
		fmt.Errorf("[ERROR] error creating tags: %s", err)
	}

	params := ipam.NewIpamIPAddressesPartialUpdateParams().WithID(ipID)
	params.WithData(&models.WritableIPAddress{
		Address:     &address,
		Description: description,
		DNSName:     dnsName,
		Role:        role,
		Tags:        tags,
	})

	if _, err := ipamClient.IpamIPAddressesPartialUpdate(params, nil); err != nil {
		return fmt.Errorf("[ERROR] error while updating address '%s' with error message: %s", address, err)
	}

	return resourceIPAddressRead(d, meta)
}

func resourceIPAddressDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).Ipam

	ipID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("[ERROR] error conversing ID '%d' with error: %s", ipID, err)
	}

	params := ipam.NewIpamIPAddressesDeleteParams()
	params.SetID(ipID)

	if _, err := conn.IpamIPAddressesDelete(params, nil); err != nil {
		return fmt.Errorf("[ERROR] error while deleting the IP address with error: %s", err)
	}

	d.SetId("")

	return nil
}
