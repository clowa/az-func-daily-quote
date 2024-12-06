locals {
  custom_hostname = "quotes.clowa.dev"
}

################################################################################
# Resource Group

resource "azurerm_resource_group" "quotes" {
  name     = "${local.global_prefix}-func-rg"
  location = local.region
  tags     = local.tags
}

################################################################################
# Key Vault Access

resource "azurerm_role_assignment" "function_key_vault_user" {
  scope                = azurerm_key_vault.this.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_linux_function_app.quotes.identity[0].principal_id
}

################################################################################
# Storage Account

resource "azurerm_storage_account" "quotes" {
  name                = replace("${local.global_prefix}-stac", "-", "")
  resource_group_name = azurerm_resource_group.quotes.name
  location            = azurerm_resource_group.quotes.location
  tags                = local.tags

  cross_tenant_replication_enabled = false
  account_tier                     = "Standard"
  account_kind                     = "StorageV2"
  access_tier                      = "Hot"
  account_replication_type         = "LRS"
  min_tls_version                  = "TLS1_2"
  allow_nested_items_to_be_public  = false
  default_to_oauth_authentication  = true

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

################################################################################
# Function App

resource "azurerm_linux_function_app" "quotes" {
  name                = "${local.global_prefix}-func"
  resource_group_name = azurerm_resource_group.quotes.name
  location            = azurerm_resource_group.quotes.location
  tags = merge(local.tags, {
    "hidden-link: /app-insights-conn-string"         = azurerm_application_insights.this.connection_string
    "hidden-link: /app-insights-instrumentation-key" = azurerm_application_insights.this.instrumentation_key
    "hidden-link: /app-insights-resource-id"         = azurerm_application_insights.this.id
  })

  service_plan_id            = azurerm_service_plan.this.id
  storage_account_name       = azurerm_storage_account.quotes.name
  storage_account_access_key = azurerm_storage_account.quotes.primary_access_key
  # storage_uses_managed_identity = true
  https_only = true

  identity {
    type = "SystemAssigned"
  }

  app_settings = {
    ## Required for Run From Package deployment, eg. if deployed from VS Code or Function Tools
    # AzureWebJobsStorage   = azurerm_storage_account.quotes.primary_connection_string # Same as storage_account_access_key
    WEBSITE_MOUNT_ENABLED = "1"

    ## Required for MSI to access the storage account
    ## See: https://learn.microsoft.com/en-us/azure/azure-functions/functions-reference?tabs=blob&pivots=programming-language-powershell#connecting-to-host-storage-with-an-identity
    # AzureWebJobsStorage__blobServiceUri  = azurerm_storage_account.quotes.primary_blob_endpoint
    # AzureWebJobsStorage__queueServiceUri = azurerm_storage_account.quotes.primary_queue_endpoint
    # AzureWebJobsStorage__tableServiceUri = azurerm_storage_account.quotes.primary_table_endpoint

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
      tags["hidden-link: /app-insights-conn-string"], # inconsistent formating of Azure API
    ]
  }
}

## Required to use MSI of the function app to access the storage account.
## See: https://learn.microsoft.com/en-us/azure/azure-functions/functions-bindings-storage-blob-trigger?tabs=python-v2%2Cisolated-process%2Cnodejs-v4&pivots=programming-language-powershell#grant-permission-to-the-identity
resource "azurerm_role_assignment" "function_storage_account_permissions" {
  for_each = toset(["Storage Blob Data Owner", "Storage Queue Data Contributor"])

  scope                = azurerm_storage_account.quotes.id
  role_definition_name = each.value
  principal_id         = azurerm_linux_function_app.quotes.identity[0].principal_id
}

resource "azurerm_role_assignment" "developer_storage_account_contributor" {
  scope                = azurerm_storage_account.quotes.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

# resource "azurerm_linux_function_app_slot" "preview" {
#   name                 = "preview"
#   function_app_id      = azurerm_linux_function_app.quotes.id
#   storage_account_name = azurerm_storage_account.quotes.name

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
  app_service_name    = azurerm_linux_function_app.quotes.name
  resource_group_name = azurerm_linux_function_app.quotes.resource_group_name
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
