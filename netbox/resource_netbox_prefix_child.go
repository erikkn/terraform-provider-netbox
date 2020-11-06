package netbox

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourcePrefixChild() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrefixChildCreate,
		Read:   resourcePrefixChildRead,
		Update: resourcePrefixChildUpdate,
		Delete: resourcePrefixChildDelete,

		Schema: map[string]*schema.Schema{
			"parent_prefix_tags": {
				Description: "(Required) The tags that are configured on the parent Prefix from which you want to create this child Prefix.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},

			"cidr_prefix_length": {
				Description: "(Required) The CIDR Prefix length you require, e.g. 20 or 22. Don't pass `/` only the number of the desired CIDR.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"description": {
				Description: "(Optional) Description of this new child Prefix",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"is_pool": {
				Description: "(Optional) All IP addresses within this prefix are considered usable",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},

			"role": {
				Description: "A role indicates the function of the prefix, e.g. the VPC name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"tags": {
				Description: "(Optional) A map of tags to assign to the resource.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourcePrefixChildCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client)
	ipamClient := conn.Ipam

	var parentPrefixTags []interface{} = d.Get("parent_prefix_tags").(*schema.Set).List()
	requestedParentPrefixTags := tagListOfStrings(&parentPrefixTags)
	//description := d.Get("description").(string)
	prefixLength, err := strconv.ParseInt(d.Get("cidr_prefix_length").(string), 10, 64)
	if err != nil {
		return fmt.Errorf("[ERROR] error conversing `cidr_prefix_length` to int64 with error message: %s", err)
	}

	prefixListParams := ipam.NewIpamPrefixesListParams()
	prefixesList, err := ipamClient.IpamPrefixesList(prefixListParams, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] error fetching the list of all Prefixes with error message: %s", err)
	}

	parentPrefixID, err := getRequestedPrefixID(prefixesList, requestedParentPrefixTags)
	if err != nil {
		return fmt.Errorf("[ERROR] error fetching the Parent Prefix with error message: %s", err)
	}

	data := &models.PrefixLength{PrefixLength: &prefixLength}
	prefixCreateParams := ipam.NewIpamPrefixesAvailablePrefixesCreateParams().WithID(*parentPrefixID).WithData(data)
	prefix, err := ipamClient.IpamPrefixesAvailablePrefixesCreate(prefixCreateParams, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] error creating `AvailablePrefixesCreate` with error message: %s ", err)
	}

	d.SetId(strconv.FormatInt(prefix.Payload.ID, 10))
	return nil
}

func resourcePrefixChildRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourcePrefixChildUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourcePrefixChildDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
