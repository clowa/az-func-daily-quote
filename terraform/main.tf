locals {
  custom_hostname = "quotes.clowa.dev"
}

################################################################################
# Resource Group

resource "azurerm_resource_group" "this" {
  name     = "${local.global_prefix}-rg"
  location = local.region
  tags     = local.tags
}

################################################################################
# Key Vault

#trivy:ignore:AVD-AZU-0013
resource "azurerm_key_vault" "this" {
  name                       = "${local.global_prefix}-kv"
  location                   = azurerm_resource_group.this.location
  resource_group_name        = azurerm_resource_group.this.name
  tags                       = local.tags
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days = 7
  purge_protection_enabled   = false # ToDo: Enable purge protection
  sku_name                   = "standard"
  enable_rbac_authorization  = true
}

resource "azurerm_role_assignment" "function_key_vault_user" {
  scope                = azurerm_key_vault.this.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_linux_function_app.this.identity[0].principal_id
}

resource "azurerm_role_assignment" "developer_key_vault_administrator" {
  scope                = azurerm_key_vault.this.id
  role_definition_name = "Key Vault Administrator"
  principal_id         = data.azurerm_client_config.current.object_id
}

################################################################################
# Storage Account

resource "azurerm_storage_account" "this" {
  name                = replace("${local.global_prefix}-stac", "-", "")
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags                = local.tags

  account_tier                    = "Standard"
  account_kind                    = "StorageV2"
  access_tier                     = "Hot"
  account_replication_type        = "LRS"
  min_tls_version                 = "TLS1_2"
  default_to_oauth_authentication = true

  blob_properties {
    change_feed_enabled      = false
    last_access_time_enabled = false
    versioning_enabled       = false
  }

  #trivy:ignore:AVD-AZU-0012
  network_rules {
    bypass         = ["AzureServices"]
    default_action = "Allow"
  }
}

resource "azurerm_storage_management_policy" "function_releases" {
  storage_account_id = azurerm_storage_account.this.id

  rule {
    name    = "function-releases"
    enabled = true

    filters {
      prefix_match = ["function-releases/"]
      blob_types   = ["blockBlob"]
    }

    actions {
      base_blob {
        tier_to_cold_after_days_since_creation_greater_than = 7
        delete_after_days_since_creation_greater_than       = 97
      }
    }
  }
}

################################################################################
# Log Analytics / Application Insights

resource "azurerm_log_analytics_workspace" "this" {
  name                = "${local.global_prefix}-log"
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags                = local.tags

  sku = "PerGB2018"
}

data "azapi_resource_id" "app_traces_table" {
  type      = "Microsoft.OperationalInsights/workspaces/tables@2022-10-01"
  name      = "AppTraces"
  parent_id = azurerm_log_analytics_workspace.this.id
}

resource "azapi_resource_action" "app_traces_table_basic" {
  type        = "Microsoft.OperationalInsights/workspaces/tables@2022-10-01"
  resource_id = data.azapi_resource_id.app_traces_table.id
  method      = "PATCH"

  body = jsonencode({
    properties = {
      plan                 = "Basic"
      totalRetentionInDays = 90
    }
  })
}

resource "azurerm_application_insights" "this" {
  name                = "${local.global_prefix}-appi"
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags                = local.tags

  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.this.id
  sampling_percentage = 0
}

################################################################################
# Function App

resource "azurerm_service_plan" "this" {
  name                = "${local.global_prefix}-asp"
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags                = local.tags #

  os_type  = "Linux"
  sku_name = "Y1"
}

resource "azurerm_linux_function_app" "this" {
  name                = "${local.global_prefix}-func"
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags = merge(local.tags, {
    "hidden-link: /app-insights-conn-string"         = azurerm_application_insights.this.connection_string
    "hidden-link: /app-insights-instrumentation-key" = azurerm_application_insights.this.instrumentation_key
    "hidden-link: /app-insights-resource-id"         = azurerm_application_insights.this.id
  })

  service_plan_id            = azurerm_service_plan.this.id
  storage_account_name       = azurerm_storage_account.this.name
  storage_account_access_key = azurerm_storage_account.this.primary_access_key
  # storage_uses_managed_identity = true
  https_only = true

  identity {
    type = "SystemAssigned"
  }

  app_settings = {
    ## Required for Run From Package deployment, eg. if deployed from VS Code or Function Tools
    # AzureWebJobsStorage   = azurerm_storage_account.this.primary_connection_string # Same as storage_account_access_key
    WEBSITE_MOUNT_ENABLED = "1"

    ## Required for MSI to access the storage account
    ## See: https://learn.microsoft.com/en-us/azure/azure-functions/functions-reference?tabs=blob&pivots=programming-language-powershell#connecting-to-host-storage-with-an-identity
    # AzureWebJobsStorage__blobServiceUri  = azurerm_storage_account.this.primary_blob_endpoint
    # AzureWebJobsStorage__queueServiceUri = azurerm_storage_account.this.primary_queue_endpoint
    # AzureWebJobsStorage__tableServiceUri = azurerm_storage_account.this.primary_table_endpoint

    MONGODB_CONNECTION_STRING = "@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.mongodb_primary_connection_string.versionless_id}/)"
    MONGODB_DATABASE          = azurerm_cosmosdb_mongo_database.quotes.name
    MONGODB_COLLECTION        = azurerm_cosmosdb_mongo_collection.quotes.name
  }

  site_config {
    always_on       = false
    app_scale_limit = 200
    ftps_state      = "FtpsOnly"

    application_insights_connection_string = azurerm_application_insights.this.connection_string
    application_insights_key               = azurerm_application_insights.this.instrumentation_key

    application_stack {
      use_custom_runtime = true
    }
    cors {
      allowed_origins = ["https://portal.azure.com"]
    }

    ## Access Restrictions aren't needed if endpoints require token authentication
    # ip_restriction_default_action = "Deny"
    # ip_restriction {
    #   action     = "Allow"
    #   priority   = 400
    #   name       = "AllowFrontendDomain"
    #   ip_address = "0.0.0.0/0"
    #   headers = [{
    #     x_azure_fdid      = []
    #     x_fd_health_probe = []
    #     x_forwarded_for   = []
    #     x_forwarded_host  = ["api.clowa.dev"]
    #   }]
    # }
  }

  sticky_settings {
    app_setting_names = [
      "APPINSIGHTS_INSTRUMENTATIONKEY",
      "APPLICATIONINSIGHTS_CONNECTION_STRING ",
      "APPINSIGHTS_PROFILERFEATURE_VERSION",
      "APPINSIGHTS_SNAPSHOTFEATURE_VERSION",
      "ApplicationInsightsAgent_EXTENSION_VERSION",
      "XDT_MicrosoftApplicationInsights_BaseExtensions",
      "DiagnosticServices_EXTENSION_VERSION",
      "InstrumentationEngine_EXTENSION_VERSION",
      "SnapshotDebugger_EXTENSION_VERSION",
      "XDT_MicrosoftApplicationInsights_Mode",
      "XDT_MicrosoftApplicationInsights_PreemptSdk",
      "APPLICATIONINSIGHTS_CONFIGURATION_CONTENT",
      "XDT_MicrosoftApplicationInsightsJava",
      "XDT_MicrosoftApplicationInsights_NodeJS"
    ]
  }

  lifecycle {
    ignore_changes = [
      # storage_account_access_key,
      # app_settings["AzureWebJobsStorage"],
      app_settings["WEBSITE_RUN_FROM_PACKAGE"],
      tags["hidden-link: /app-insights-resource-id"], # inconsistent formating of Azure API
    ]
  }
}

## Required to use MSI of the function app to access the storage account.
## See: https://learn.microsoft.com/en-us/azure/azure-functions/functions-bindings-storage-blob-trigger?tabs=python-v2%2Cisolated-process%2Cnodejs-v4&pivots=programming-language-powershell#grant-permission-to-the-identity
resource "azurerm_role_assignment" "function_storage_account_permissions" {
  for_each = toset(["Storage Blob Data Owner", "Storage Queue Data Contributor"])

  scope                = azurerm_storage_account.this.id
  role_definition_name = each.value
  principal_id         = azurerm_linux_function_app.this.identity[0].principal_id
}

resource "azurerm_role_assignment" "developer_storage_account_contributor" {
  scope                = azurerm_storage_account.this.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

# resource "azurerm_linux_function_app_slot" "preview" {
#   name                 = "preview"
#   function_app_id      = azurerm_linux_function_app.this.id
#   storage_account_name = azurerm_storage_account.this.name

#   site_config {
#     ftps_state = "Disabled"

#     application_insights_connection_string = azurerm_application_insights.this.connection_string
#     application_insights_key               = azurerm_application_insights.this.instrumentation_key

#     application_stack {
#       powershell_core_version = "7.2"
#     }
#   }
# }

################################################################################
# Custom Domain

resource "azurerm_app_service_custom_hostname_binding" "quotes_clowa_dev" {
  app_service_name    = azurerm_linux_function_app.this.name
  resource_group_name = azurerm_linux_function_app.this.resource_group_name
  hostname            = local.custom_hostname
}

resource "azurerm_app_service_managed_certificate" "quotes_clowa_dev" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.quotes_clowa_dev.id
}

resource "azurerm_app_service_certificate_binding" "quotes_clowa_dev" {
  hostname_binding_id = azurerm_app_service_custom_hostname_binding.quotes_clowa_dev.id
  certificate_id      = azurerm_app_service_managed_certificate.quotes_clowa_dev.id
  ssl_state           = "SniEnabled"
}
