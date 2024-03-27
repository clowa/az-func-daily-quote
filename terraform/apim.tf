################################################################################
# API Management

resource "azurerm_api_management" "apim" {
  name                = "${local.global_prefix}-apim"
  resource_group_name = azurerm_resource_group.this.name
  location            = azurerm_resource_group.this.location
  publisher_name      = "Clowa"
  publisher_email     = "example@clowa.de"

  sku_name = "Consumption_0"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_api_management_logger" "apim_application_insights" {
  name                = "application-insights"
  resource_group_name = azurerm_api_management.apim.resource_group_name
  api_management_name = azurerm_api_management.apim.name
  application_insights {
    instrumentation_key = azurerm_application_insights.this.instrumentation_key
  }
  resource_id = azurerm_application_insights.this.id
}

resource "azurerm_api_management_diagnostic" "apim" {
  resource_group_name      = azurerm_api_management.apim.resource_group_name
  api_management_name      = azurerm_api_management.apim.name
  identifier               = "applicationinsights"
  api_management_logger_id = azurerm_api_management_logger.apim_application_insights.id

  sampling_percentage = 10.0
  always_log_errors   = true
  log_client_ip       = true
  verbosity           = "information"
}

# resource "azurerm_api_management_policy" "apim" {
#   api_management_id = azurerm_api_management.apim.id
#   xml_content       = <<XML
# <policies>
# <cors allow-credentials="false | true" terminate-unmatched-request="true | false">
#     <allowed-origins>
#         <origin>*</origin>
#     </allowed-origins>
#     <allowed-methods preflight-result-max-age="60">
#         <method>GET</method>
#         <method>POST</method>
#     </allowed-methods>
#     <allowed-headers>
#         <header>header name</header>
#     </allowed-headers>
#     <expose-headers>
#         <header>header name</header>
#     </expose-headers>
# </cors>
# XML
# }

################################################################################
# Key Vault

resource "azurerm_key_vault" "apim" {
  name                      = "${local.global_prefix}-apim-kv"
  resource_group_name       = azurerm_resource_group.this.name
  location                  = azurerm_resource_group.this.location
  enable_rbac_authorization = true
  tenant_id                 = data.azurerm_client_config.current.tenant_id
  sku_name                  = "standard"
}

resource "azurerm_role_assignment" "apim_key_vault_secrets_officer" {
  scope                = azurerm_key_vault.apim.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_api_management.apim.identity[0].principal_id
}
