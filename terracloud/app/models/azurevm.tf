{
    "variable": {
        "client_id": {
            "type": "string"
        },
        "client_secret": {
            "type": "string"
        },
        "tenant_id": {
            "type": "string"
        }
    },
    "provider": {
        "azurerm": {
            "version": "=2.4.0",
            "subscription_id": "3dc3cd1a-d5cd-4e3e-a648-b2253048af83",
            "client_id": "${var.client_id}",
            "client_secret": "${var.client_secret}",
            "tenant_id": "${var.tenant_id}",
            "features": {}
        }
    },
    "module": {
        "azureVM": {
            "source": "app.terraform.io/ClDevTeam/vm/azurerm",
            "version": "1.0.3"
        }
    }
}
