################################################################################
# Cosmos DB

resource "azurerm_cosmosdb_account" "this" {
  name                = "${local.global_prefix}-cosno"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  tags                = local.tags

  offer_type = "Standard"
  kind       = "GlobalDocumentDB"
  # enable_free_tier = true

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

  backup {
    type                = "Periodic"
    interval_in_minutes = "60"
    retention_in_hours  = "720"
    storage_redundancy  = "Geo"
  }
}

resource "azurerm_key_vault_secret" "cosmos_password" {
  key_vault_id = azurerm_key_vault.this.id
  name         = "cosmos-password"
  value        = azurerm_cosmosdb_account.this.primary_key
}

resource "azurerm_cosmosdb_sql_database" "quotes" {
  name                = "quotes-cosmos" # should end with -cosmos
  resource_group_name = azurerm_resource_group.this.name
  account_name        = azurerm_cosmosdb_account.this.name
}

resource "azurerm_cosmosdb_sql_container" "quotes" {
  name                = "quotes"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = azurerm_cosmosdb_account.this.name
  database_name       = azurerm_cosmosdb_sql_database.quotes.name
  partition_key_path  = "/authorSlug"

  indexing_policy {
    indexing_mode = "consistent"

    included_path { path = "/*" }
    included_path { path = "/creationDate/?" }
    # excluded_path { path = "/_etag" }
  }
}

resource "azurerm_monitor_diagnostic_setting" "cosmos_db" {
  name                           = "${local.global_prefix}-cosmos-db"
  target_resource_id             = azurerm_cosmosdb_account.this.id
  log_analytics_workspace_id     = azurerm_log_analytics_workspace.this.id
  log_analytics_destination_type = "Dedicated"

  dynamic "enabled_log" {
    for_each = ["ControlPlaneRequests", "DataPlaneRequests", "MongoRequests", "QueryRuntimeStatistics", "PartitionKeyStatistics"]
    content {
      category = enabled_log.value
    }
  }
  metric {
    category = "Requests"
    enabled  = false
  }
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
  scope               = "${azurerm_cosmosdb_account.this.id}/dbs/${azurerm_cosmosdb_sql_database.quotes.name}"
  principal_id        = azurerm_linux_function_app.this.identity[0].principal_id
}
