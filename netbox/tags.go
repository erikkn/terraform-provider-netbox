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

func tagCreate(t *string, conn extras.ClientService) error {
	exists, err := tagExists(t, conn)
	if err != nil {
		log.Printf("[ERROR] error checking if '%s' exists with error message: %s", *t, err)
		return err
	}

	if !exists {
		data := &models.Tag{
			Name: t,
			Slug: t,
		}

		params := extras.NewExtrasTagsCreateParams().WithData(data)
		if _, err := conn.ExtrasTagsCreate(params, nil); err != nil {
			log.Printf("[ERROR] error creating tag '%s' with error message: %s", *t, err)
			return err
		}
		log.Printf("[DEBUG] successfully created the non-existing tag '%s' ", *t)
	}

	return nil
}

func nestedTagCreate(t *[]interface{}, conn extras.ClientService) ([]*models.NestedTag, error) {
	nestedTags := make([]*models.NestedTag, 0, len(*t))

	for _, v := range *t {
		tag := fmt.Sprint(v)
		nestedTags = append(nestedTags, &models.NestedTag{Name: &tag, Slug: &tag})
		if err := tagCreate(&tag, conn); err != nil {
			return nil, err
		}
	}

	return nestedTags, nil
}

func tagListOfStrings(tags *[]interface{}) *[]string {
	newList := make([]string, len(*tags))

	for i, v := range *tags {
		newList[i] = fmt.Sprint(v)
	}

	return &newList
}
