package netbox

import (
	"errors"
	"log"

	"github.com/netbox-community/go-netbox/netbox/client/tenancy"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func tenantGroupExists(tg *string, conn tenancy.ClientService) (bool, error) {
	// Using the `List` method here instead of `Read`, because `Read` only accepts an ID as input, not the name of the TG.
	params := tenancy.NewTenancyTenantGroupsListParams().WithName(tg)
	exist, err := conn.TenancyTenantGroupsList(params, nil)
	if err != nil {
		log.Printf("[ERROR] error while fetching TenantGroup '%s' with error: %s", *tg, err)
		return false, err
	}

	if len(exist.Payload.Results) < 1 {
		return false, nil
	} else if len(exist.Payload.Results) > 1 {
		log.Printf("[ERROR] duplicated groups found for %s ", *tg)
		return true, errors.New("[ERROR] duplicated groups")
	}

	return true, nil
}

func tenantGroupCreate(tg *string, conn tenancy.ClientService) (int64, error) {
	params := tenancy.NewTenancyTenantGroupsCreateParams()
	params.WithData(&models.WritableTenantGroup{
		Name: tg,
		Slug: tg,
	})

	tenantGroup, err := conn.TenancyTenantGroupsCreate(params, nil)
	if err != nil {
		log.Printf("[ERROR] error creating Tenant Group '%s' with error message: %s", *tg, err)
		return 0, err
	}

	return tenantGroup.Payload.ID, nil
}

func tenantGroupFetchID(tg *string, conn tenancy.ClientService) (int64, error) {
	exists, err := tenantGroupExists(tg, conn)
	if err != nil {
		return 0, err
	}

	if !exists {
		id, err := tenantGroupCreate(tg, conn)
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	params := tenancy.NewTenancyTenantGroupsListParams().WithName(tg)
	fetch, err := conn.TenancyTenantGroupsList(params, nil)
	if err != nil {
		log.Printf("[ERROR] error while fetching TenantGroup '%s' with error: %s", *tg, err)
		return 0, err
	}

	// We can do this because the `tenantGroupExists` function is called earlier at the beginning of this func.
	return fetch.Payload.Results[0].ID, nil
}
