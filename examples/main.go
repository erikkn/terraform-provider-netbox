package main

import (
	"context"
	"fmt"

	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
)

func main() {
	c := netbox.NewNetboxWithAPIKey("127.0.0.1:8000", "0123456789abcdef0123456789abcdef01234567").Ipam

	dnsName := "test.com"
	ipAddress := "1.1.1.1/32"

	test0 := ipam.IpamIPAddressesListParams{
		DNSName: &dnsName,
		Address: &ipAddress,
		Context: context.Background(),
	}

	test1, err := c.IpamIPAddressesList(&test0, nil)
	if err != nil {
		fmt.Errorf("Error: ", err)
	}

	//var role *models.IPAddressRole

	//var role *models.IPAddressRole

	for _, v := range test1.Payload.Results {
		fmt.Println("Address is: ", *v.Address)
		fmt.Println("DNS Name is: ", v.DNSName)
		fmt.Println("Family is: ", *v.Family.Label)

		//role = v.Role
		if v.Role == nil {
			fmt.Errorf("ERRRRORRR!")
			fmt.Println("Error bitches!")
		}
	}
}
