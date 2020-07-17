package netbox

import (
	"github.com/netbox-community/go-netbox/netbox"
	"github.com/netbox-community/go-netbox/netbox/client"
)

type Config struct {
	api_key string
	host    string
}

type Client struct {
	*client.NetBox
}

func (c *Config) Client() (interface{}, error) {
	netconn := Client{netbox.NewNetboxWithAPIKey(c.host, c.api_key)}
	return &netconn, nil
}
