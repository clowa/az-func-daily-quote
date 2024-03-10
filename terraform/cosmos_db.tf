import {
  id = "/subscriptions/0a0a4299-b306-4dad-94de-862e8405fdbe/resourceGroups/cwa-quotes-prod-rg/providers/Microsoft.DocumentDB/databaseAccounts/cwa-quotes-prod-cosmos"
  to = azurerm_cosmosdb_account.this
}

resource "azurerm_cosmosdb_account" "this" {
  name                = "${local.global_prefix}-cosmos"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  tags                = local.tags

  offer_type = "Standard"
  kind       = "GlobalDocumentDB"

  capabilities {
    name = "EnableServerless"
  }

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    failover_priority = 0
    location          = "westeurope"
    zone_redundant    = false
  }
}

resource "azurerm_key_vault_secret" "cosmos_password" {
  key_vault_id = azurerm_key_vault.this.id
  name         = "cosmos-password"
  value        = azurerm_cosmosdb_account.this.primary_key
}

resource "azurerm_cosmosdb_sql_database" "quotes" {
  name                = "quotes"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = azurerm_cosmosdb_account.this.name
}

resource "azurerm_cosmosdb_sql_container" "quotes" {
  name                = "quotes"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = azurerm_cosmosdb_account.this.name
  database_name       = azurerm_cosmosdb_sql_database.quotes.name
  ## Time is normally a "hot partition key" and should be avoided,
  ##since it can lead to many concurrent requests going to the same pyhsical partition.
  partition_key_path = "/timestamp"
}
