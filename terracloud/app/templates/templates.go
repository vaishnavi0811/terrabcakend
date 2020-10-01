package templates

type Designs struct {
	ID            string
	Resourceid    string
	Dependson     interface{}
	Configuration interface{}
}

type Resource struct {
	ResourceID   string `json:"resource_id" validate:"required"`
	ResourceType string `json:"resource_type" validate:"required"`
	Config       interface{}
	DependsOn    []string
}

type Resources struct {
	Resources interface{}
}
type ApplyPlan struct {
	ApplyMessage string
}
type MVMVARS struct {
	Location                string            `json:"location" validate:"required"`
	VMName                  string            `json:"vm_name" validate:"required"`
	ResourceGroupName       string            `json:"rg_name" validate:"required"`
	AdminUsername           string            `json:"admin_username" validate:"required"`
	AdminPassword           string            `json:"admin_password" validate:"required"`
	VMSku                   string            `json:"vm_sku" validate:"required"`
	VMSize                  string            `json:"vm_size" validate:"required"`
	OSDataDiskSizeInGB      int               `json:"osdatadisksizeingb" validate:"required,gte=127"`
	DataDisks               []int             `json:"data_disks" validate:"required,gt=0"`
	BootDiagStorageEndpoint string            `json:"boot_diag_storage_account_id"`
	VnetName                string            `json:"vnet_name" validate:"required"`
	VnetRGroup              string            `json:"vnet_rgroup" validate:"required"`
	SubnetName              string            `json:"subnet_name" validate:"required"`
	AvailabilitySet         string            `json:"availability_set"`
	Identity                string            `json:"identity_id"`
	Tags                    map[string]string `json:"tags" validate:"required"`
	SubscriptionID          string            `json:"subscription_id" validate:"required"`
}
type RGVARS struct {
	Location string            `json:"location" validate:"required"`
	Rgname   string            `json:"rg_name" validate:"required"`
	Tags     map[string]string `json:"tags" validate:"required"`
}
