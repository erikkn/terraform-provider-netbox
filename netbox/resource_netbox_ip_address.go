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
	var rawTags []interface{} = d.Get("tags").(*schema.Set).List()

	tags, err := nestedTagCreate(&rawTags, extrasClient)
	if err != nil {
		return fmt.Errorf("[ERROR] error creating tags: %s", err)
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

func resourceIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).Ipam

	ipID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("[ERROR] unable to convert ID '%s' to int64 with error message: %s", d.Id(), err)
	}

	params := ipam.NewIpamIPAddressesReadParams().WithID(ipID)
	fetch, err := conn.IpamIPAddressesRead(params, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] unable to fetch the IP address with error message: %s", err)
	}

	if fetch == nil {
		log.Printf("[WARN] IP address with ID '%d' not found in Netbox; Removing from state.", ipID)
		d.SetId("")
		return nil
	}

	tags := make([]string, 0, len(fetch.Payload.Tags))
	for _, v := range fetch.Payload.Tags {
		tags = append(tags, *v.Name)
	}

	d.Set("address", *fetch.Payload.Address)
	d.Set("description", fetch.Payload.Description)
	d.Set("dns_name", fetch.Payload.DNSName)
	d.Set("net_family", fetch.Payload.Family)
	d.Set("role", fetch.Payload.Role)
	d.Set("tags", tags)

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
	var rawTags []interface{} = d.Get("tags").(*schema.Set).List()

	tags, err := nestedTagCreate(&rawTags, extrasClient)
	if err != nil {
		return fmt.Errorf("[ERROR] error creating tags: %s", err)
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

	params := ipam.NewIpamIPAddressesDeleteParams().WithID(ipID)
	//params.SetID(ipID)

	if _, err := conn.IpamIPAddressesDelete(params, nil); err != nil {
		return fmt.Errorf("[ERROR] error while deleting the IP address with error: %s", err)
	}

	d.SetId("")

	return nil
}
