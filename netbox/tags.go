package netbox

import (
	"fmt"
	"log"

	"github.com/netbox-community/go-netbox/netbox/client/extras"
	"github.com/netbox-community/go-netbox/netbox/models"
)

// Checks whether a tag already exists in Netbox.
func tagExists(tag string, meta interface{}) (bool, error) {
	conn := meta.(*Client).Extras

	params := extras.NewExtrasTagsListParams().WithName(&tag)
	exist, err := conn.ExtrasTagsList(params, nil)
	if err != nil {
		log.Printf("[ERROR] Error while fetching the tag '%s' with error message: %s", tag, err)
		return false, err
	}

	if *exist.Payload.Count < 1 {
		return false, nil
	}

	return true, nil
}

// Creates the tags that don't exist yet (by callin the `tagExist`).
func createTag(tags []string, meta interface{}) (nestedTags []*models.NestedTag, err error) {
	conn := meta.(*Client).Extras
	tagsToCreate := []string{}
	// TODO: Intializing an empty slice will result in an memory error. Adding a default value works for now, because NetBox automatically gets rid of duplicated tags.
	returnTags := []*models.NestedTag{{Name: &tags[0], Slug: &tags[0]}}

	for _, v := range tags {
		t := fmt.Sprint(v)
		returnTags = append(returnTags, &models.NestedTag{Name: &t, Slug: &t})

		ok, err := tagExists(v, meta)
		if err != nil {
			log.Printf("[ERROR] error while checking tag '%s' exists, with error: %s", v, err)
			return nil, err
		}

		if !ok {
			tagsToCreate = append(tagsToCreate, v)
		}
	}

	for _, v := range tagsToCreate {
		t := fmt.Sprint(v)
		data := models.Tag{
			Name: &t,
			Slug: &t,
		}

		params := extras.NewExtrasTagsCreateParams().WithData(&data)
		if _, err := conn.ExtrasTagsCreate(params, nil); err != nil {
			log.Printf("[ERROR] error while creating tag '%s' with error: %s", v, err)
			return nil, err
		}
	}

	return returnTags, nil
}
