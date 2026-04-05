resource "azurerm_storage_account" "main" {
  name                     = "harkportfoliostore"
  resource_group_name      = azurerm_resource_group.main.name
  location                 = azurerm_resource_group.main.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = var.environment
    project     = "hark-portfolio"
  }
}

resource "azurerm_storage_container" "images" {
  name                  = "art-images"
  storage_account_name  = azurerm_storage_account.main.name
  container_access_type = "blob"
}
