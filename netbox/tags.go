package netbox

import (
	"fmt"
	"log"

	"github.com/netbox-community/go-netbox/netbox/client/extras"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func tagExists(t *string, conn extras.ClientService) (exists bool, err error) {
	// Using the `ListParams` method here instead of `Read`, because we want to lookup the tag based on the `Name`.
	params := extras.NewExtrasTagsListParams().WithName(t)
	fetch, err := conn.ExtrasTagsList(params, nil)
	if err != nil {
		log.Printf("[ERROR] error fetching tag '%s' with error message: '%s'", *t, err)
		return false, err
	}
	log.Printf("[INFO] successfully fetched tag '%s' from upstream", *t)

	if len(fetch.Payload.Results) < 1 {
		log.Printf("[DEBUG] tag '%s' not found, meaning it doesn't exist. Returning False.", *t)
		return false, nil
	}

	return true, nil
}

func createTag(t *[]interface{}, conn extras.ClientService) (nestedTags []*models.NestedTag, err error) {
	var nestedTag []*models.NestedTag

	for _, v := range *t {
		tag := fmt.Sprint(v)
		nestedTag = append(nestedTag, &models.NestedTag{Name: &tag, Slug: &tag})
		exists, err := tagExists(&tag, conn)
		if err != nil {
			log.Printf("[ERROR] error checking if '%s' exists with error message: %s", tag, err)
			return nil, err
		}

		if !exists {
			data := &models.Tag{
				Name: &tag,
				Slug: &tag,
			}

			params := extras.NewExtrasTagsCreateParams().WithData(data)
			if _, err := conn.ExtrasTagsCreate(params, nil); err != nil {
				log.Printf("[ERROR] error creating tag '%s' with error message: %s", tag, err)
				return nil, err
			}
			log.Printf("[DEBUG] successfully created the non-existing tag '%s' ", tag)
		}
	}

	return nestedTag, nil
}
