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

resource "azurerm_role_assignment" "developer_key_vault_administrator" {
  scope                = azurerm_key_vault.this.id
  role_definition_name = "Key Vault Administrator"
  principal_id         = data.azurerm_client_config.current.object_id
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
# App Service Plan

resource "azurerm_service_plan" "this" {
  name                = "${local.global_prefix}-asp"
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  tags                = local.tags #

  os_type  = "Linux"
  sku_name = "Y1"
}
