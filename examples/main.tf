provider "netbox" {
  api_key = "0123456789abcdef0123456789abcdef01234567"
  host    = "127.0.0.1:8000"
}

resource "netbox_ip_address" "test" {
  address     = "1.1.1.21/32"
  description = "test dinges4"
  dns_name    = "test4.example.com"
  role        = "vrrp"

  tags = ["test5", "test7"]
}
