provider "netbox" {
  api_key = "0123456789abcdef0123456789abcdef01234567"
  host    = "127.0.0.1:8000"
}

data "netbox_ip_address" "foobar" {
  //cidr_block = "1.1.1.1/32"
  dns_name = "razu.ee"
  //  address_family = "IPv4"
}

//output "foobar" {
//  value = data.netbox_ip_address.foobar.cidr_block
//}
