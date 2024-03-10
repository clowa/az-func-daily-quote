################################################################################
# Cosmos DB

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

################################################################################
# Database Role

resource "azurerm_cosmosdb_sql_role_definition" "cosmos_db_data_contributor" {
  role_definition_id  = "cec938ce-34cc-4275-8240-c533197e37c2"
  account_name        = azurerm_cosmosdb_account.this.name
  resource_group_name = azurerm_cosmosdb_account.this.resource_group_name
  name                = "Cosmos DB Data Contributor"
  assignable_scopes   = ["${azurerm_cosmosdb_account.this.id}/dbs"]

  permissions {
    ## See: https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-setup-rbac#permission-model
    data_actions = [
      "Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/create",
      "Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read",
      "Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/upsert",
    ]
  }
}

resource "azurerm_cosmosdb_sql_role_assignment" "function_cosmos_db_data_contributor" {
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.cosmos_db_data_contributor.id
  account_name        = azurerm_cosmosdb_account.this.name
  resource_group_name = azurerm_cosmosdb_account.this.resource_group_name
  scope               = "${azurerm_cosmosdb_account.this.id}/dbs/quotes"
  principal_id        = azurerm_linux_function_app.this.identity[0].principal_id
}
