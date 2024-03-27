################################################################################
# API Management Ganeral

moved {
  from = azurerm_api_management.this
  to   = azurerm_api_management.apim
}

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

################################################################################
# APIs

data "azurerm_function_app_host_keys" "quotes" {
  name                = azurerm_linux_function_app.this.name
  resource_group_name = azurerm_linux_function_app.this.resource_group_name
}

resource "azurerm_key_vault_secret" "apim_quotes_function_key" {
  key_vault_id = azurerm_key_vault.apim.id
  name         = "${azurerm_linux_function_app.this.name}-key"
  value        = data.azurerm_function_app_host_keys.quotes.default_function_key
}

resource "azurerm_api_management_named_value" "apim_quotes_function_key" {
  name                = "${azurerm_linux_function_app.this.name}-key"
  resource_group_name = azurerm_api_management.apim.resource_group_name
  api_management_name = azurerm_api_management.apim.name
  display_name        = "quotes-func-key"
  secret              = true #! Has to be set for key vault secrets
  value_from_key_vault {
    secret_id = azurerm_key_vault_secret.apim_quotes_function_key.versionless_id
  }
  tags = ["key", "function"]
}

# import {
#   id = "/subscriptions/0a0a4299-b306-4dad-94de-862e8405fdbe/resourceGroups/cwa-quotes-prod-rg/providers/Microsoft.ApiManagement/service/cwa-quotes-prod-apim/backends/cwa-quotes-prod-func"
#   to = azurerm_api_management_backend.quotes
# }

resource "azurerm_api_management_backend" "quotes" {
  name                = azurerm_linux_function_app.this.name
  api_management_name = azurerm_api_management.apim.name
  resource_group_name = azurerm_api_management.apim.resource_group_name
  protocol            = "http"
  resource_id         = "https://management.azure.com${azurerm_linux_function_app.this.id}"
  url                 = "https://${azurerm_app_service_custom_hostname_binding.quotes_clowa_dev.hostname}"

  credentials {
    header = {
      x-functions-key = "{{${azurerm_api_management_named_value.apim_quotes_function_key.name}}}"
    }
  }
}

resource "azurerm_api_management_api" "quotes" {
  name                = "quotes"
  resource_group_name = azurerm_resource_group.this.name
  api_management_name = azurerm_api_management.apim.name
  revision            = "1"
  display_name        = "Quotes API"
  path                = ""
  protocols           = ["https"]
  import {
    content_format = "openapi"
    content_value  = file("./../function/openapi.yaml")
  }
}

resource "azurerm_api_management_api_diagnostic" "quotes" {
  resource_group_name      = azurerm_api_management.apim.resource_group_name
  api_management_name      = azurerm_api_management.apim.name
  api_name                 = azurerm_api_management_api.quotes.name
  identifier               = "applicationinsights"
  api_management_logger_id = azurerm_api_management_logger.apim_application_insights.id

  sampling_percentage = 100
  always_log_errors   = true
  log_client_ip       = true
  verbosity           = "information"
}

resource "azurerm_api_management_api_policy" "quotes" {
  resource_group_name = azurerm_api_management.apim.resource_group_name
  api_management_name = azurerm_api_management.apim.name
  api_name            = azurerm_api_management_api.quotes.name
  xml_content         = <<XML
<policies>
    <inbound>
        <base />
        <set-backend-service id="terraform-generated-policy" backend-id="${azurerm_api_management_backend.quotes.name}" />
    </inbound>
    <backend>
        <base />
    </backend>
    <outbound>
        <base />
    </outbound>
    <on-error>
        <base />
    </on-error>
</policies>
XML
}
