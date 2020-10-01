variable "client_id" {
	 type = string 
}
variable "client_secret" {
	 type = string 
}
variable "tenant_id" {
	 type = string 
}
variable "subscription_id" {
	 type = string 
}

module "vm" {
	version =  "1.0.4"
	source =  "app.terraform.io/ClDevTeam/vm/azurerm"
	location = "eastus"
	vm_name = "euwmvm01"
	resource_group_name = "Testing"
	admin_username = "winuser"
	admin_password = "Password@2020"
	vm_sku = "2016-Datacenter"
	vm_size = "Standard_DS1_v2"
	os_data_disk_size_in_gb = 127
	data_disks = [ 0,120 ]
	vnet_name = "vnet001"
	vnet_r_group = "32943"
	subnet_name = "sbn001"
	tags = { 
		created-by = "32943"
		deployment_id = "ADVMMP"
	}

provider "azurerm" {
	 version =  "=2.4.0"
	 client_id = var.client_id
	 client_secret = var.client_secret
	 subscription_id = "3dc3cd1a-d5cd-4e3e-a648-b2253048af83"
	 tenant_id = var.tenant_id
	 features {}
}
