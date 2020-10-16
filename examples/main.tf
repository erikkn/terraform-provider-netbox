provider "netbox" {
  api_key = "0123456789abcdef0123456789abcdef01234567"
  host    = "127.0.0.1:8000"
}

resource "netbox_ip_address" "foobar" {
  address     = "1.1.1.1/32"
  description = "Example IP address"
  dns_name    = "foobar.example.com"
  role        = "vrrp"

  tags = [
    "foobar1",
    "foobar2",
    "foobar3",
  ]
}
