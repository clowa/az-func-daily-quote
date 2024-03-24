################################################################################
# Cosmos DB

resource "azurerm_cosmosdb_account" "this" {
  name                = "${local.global_prefix}-cosmon"
  location            = azurerm_resource_group.this.location
  resource_group_name = azurerm_resource_group.this.name
  tags                = local.tags

  offer_type = "Standard"
  kind       = "MongoDB"
  # enable_free_tier = true

  capabilities {
    name = "EnableMongo"
  }

  capabilities {
    name = "EnableMongoRoleBasedAccessControl"
  }

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

  # ip_range_filter = "0.0.0.0"
  # ip_range_filter = "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45,52.187.184.26,0.0.0.0"

  cors_rule {
    allowed_headers    = []
    allowed_methods    = []
    allowed_origins    = ["https://${local.custom_hostname}"]
    exposed_headers    = []
    max_age_in_seconds = 60
  }

  backup {
    type                = "Periodic"
    interval_in_minutes = "60"
    retention_in_hours  = "720"
    storage_redundancy  = "Geo"
  }
}

resource "azurerm_key_vault_secret" "mongodb_primary_connection_string" {
  key_vault_id = azurerm_key_vault.this.id
  name         = "primary-mongodb-connection-string"
  value        = azurerm_cosmosdb_account.this.primary_mongodb_connection_string
}

resource "azurerm_cosmosdb_mongo_database" "quotes" {
  name                = "quotes"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = azurerm_cosmosdb_account.this.name
}

resource "azurerm_cosmosdb_mongo_collection" "quotes" {
  name                = "quotes"
  resource_group_name = azurerm_resource_group.this.name
  account_name        = azurerm_cosmosdb_account.this.name
  database_name       = azurerm_cosmosdb_mongo_database.quotes.name
  default_ttl_seconds = -1
  shard_key           = "id"

  index {
    keys   = ["_id"]
    unique = true # Has to be unset during creation
  }

  index {
    keys   = ["id"]
    unique = true
  }

  index {
    keys   = ["author", "content", "creationDate", "tags"]
    unique = false
  }
}

resource "azurerm_monitor_diagnostic_setting" "cosmos_db" {
  name                           = "${local.global_prefix}-cosmos-db"
  target_resource_id             = azurerm_cosmosdb_account.this.id
  log_analytics_workspace_id     = azurerm_log_analytics_workspace.this.id
  log_analytics_destination_type = "Dedicated"

  dynamic "enabled_log" {
    for_each = ["ControlPlaneRequests", "DataPlaneRequests", "MongoRequests", "QueryRuntimeStatistics", "PartitionKeyStatistics", "PartitionKeyRUConsumption"]
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
