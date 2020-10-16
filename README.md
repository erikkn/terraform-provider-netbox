# NetBox Terraform Provider
This is a Terraform provider for DigitalOcean's NetBox project. Find out more about NetBox [here](https://netbox.readthedocs.io/en/stable/).

The goal of the NetBox project is to help people manage and document their network infrastructures. NetBox has a strong focus on physical computer networks, allowing you to document VLANs, equipment racks, etc. . While this is all great, the provider in this repository tries to focus on offering an IPAM provider for Cloud networks. This means that it won’t offer certain resources, like a resource for equipment racks.

## Preface
Special thanks to [Taavi Tuisk](https://github.com/taavituisk) and [Taras Burko](https://github.com/tburko) who, over the past years, mentored and taught me a lot. Thank you, guys :cocktail: :bike: !

## Requirements & Quickstart
You need a working NetBox setup in order to use this provider; The underlying SDK of this provider expects you to run NetBox >= 2.9. Check the [docs](https://github.com/netbox-community/netbox-docker/wiki/Getting-Started) how to get started with NetBox.

## Usage
Download the `terraform-provider-netbox` from the release page, or clone and build it locally; Next follow [instruction](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery) of how to install the plugin locally.
