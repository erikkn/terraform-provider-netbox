// Default API details from the Netbox Image
provider "netbox" {
  api_key = "0123456789abcdef0123456789abcdef01234567"
  host    = "127.0.0.1:8000"
}

resource "netbox_child_prefix" "foobar" {
  cidr_prefix_length = "28"

  parent_prefix_tags = [
    "foobar1",
    "foobar2",
  ]
}
resource "netbox_child_prefix" "foobar2" {
  cidr_prefix_length = "24"

  parent_prefix_tags = [
    "foobar3",
    "foobar4",
  ]
}
