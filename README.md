# NetBox Terraform Provider
This is a Terraform provider for DigitalOcean's NetBox project. Find out more about NetBox [here](https://netbox.readthedocs.io/en/stable/).

## This provider is deprecated!
This provider is far from done, but while in development I decided to deprecate the provider, check the section below for the motivation.

While I was writing this provider and using it in a cloud-native environment I found that using Netbox as the single source of truth for IP management is perhaps not ideal. Using Netbox as the single source of truth for IP documentation and using it as the central interface for various automation toolings means that some team needs to own Netbox and maintain it.
Netbox would need to run in a highly-available fashion, with active back-ups of all the changes made to the database. Maintaining such a tool will also bring quite some TOIL for the owning team.

Because of this, I decided to deprecate the module and fix the - for me - underlying root cause. Soon we will release an IPAM tool that has a strong cloud-native focus and doesn’t require someone to maintain it. Obviously, a Terraform provider will be available for this solution as well.

### Requirements & Quickstart
You need a working NetBox setup in order to use this provider; The underlying SDK of this provider expects you to run NetBox >= 2.9. Check the [docs](https://github.com/netbox-community/netbox-docker/wiki/Getting-Started) how to get started with NetBox.

### Usage
Download the `terraform-provider-netbox` from the release page, or clone and build it locally; Next follow [instruction](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery) of how to install the plugin locally.
