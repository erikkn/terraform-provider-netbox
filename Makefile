.DEFAULT_GOAL := build
build:
	rm -rf examples/terraform.tfstate
	rm -rf examples/terraform.tfstate.backup
	rm -rf examples/crash.log
	rm -rd examples/.terraform
	go build -o examples/terraform-provider-netbox
