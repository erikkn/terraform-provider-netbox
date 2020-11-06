package netbox

import (
	"fmt"
	"log"

	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func prefixTagisInRequestedTags(prefixesNestedTagElement *string, requestedTags *[]string) bool {
	for _, v := range *requestedTags {
		if v == *prefixesNestedTagElement {
			return true
		}
	}
	return false
}

func prefixesTagsMatchRequested(prefixesNestedTags []*models.NestedTag, requestedTags *[]string) bool {
	var temp bool

	if len(prefixesNestedTags) < 1 {
		log.Print("[LOG] Prefix does not have any tag set!")
		return false
	} else if len(*requestedTags) > len(prefixesNestedTags) {
		log.Printf("[INFO] trying to lookup more tags than are set on the actual Prefix!")
		return false
	}

	for _, v := range prefixesNestedTags {
		matchingTags := prefixTagisInRequestedTags(v.Name, requestedTags)
		if !matchingTags {
			temp = false
			break
		}

		if matchingTags {
			temp = true
		}

	}

	if temp {
		return true
	}

	return false
}

func getRequestedPrefixID(IpamPrefixesListOK *ipam.IpamPrefixesListOK, requestedTags *[]string) (*int64, error) {
	var ids []int64

	for _, v := range IpamPrefixesListOK.Payload.Results {
		matchingTags := prefixesTagsMatchRequested(v.Tags, requestedTags)
		if !matchingTags {
			continue
		}

		if matchingTags {
			ids = append(ids, v.ID)
		}
	}

	if len(ids) > 1 {
		return new(int64), fmt.Errorf("[ERROR] Duplicated Prefixes found with the requested tags combination")
	}

	return &ids[0], nil

}
