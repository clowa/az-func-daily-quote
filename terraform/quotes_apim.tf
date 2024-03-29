################################################################################
# Function App Backend

data "azurerm_function_app_host_keys" "quotes" {
  name                = azurerm_linux_function_app.quotes.name
  resource_group_name = azurerm_linux_function_app.quotes.resource_group_name
}

resource "azurerm_key_vault_secret" "apim_quotes_function_key" {
  key_vault_id = azurerm_key_vault.apim.id
  name         = "${azurerm_linux_function_app.quotes.name}-key"
  value        = data.azurerm_function_app_host_keys.quotes.default_function_key
}

resource "azurerm_api_management_named_value" "apim_quotes_function_key" {
  name                = "${azurerm_linux_function_app.quotes.name}-key"
  resource_group_name = azurerm_api_management.apim.resource_group_name
  api_management_name = azurerm_api_management.apim.name
  display_name        = "quotes-func-key"
  secret              = true #! Has to be set for key vault secrets
  value_from_key_vault {
    secret_id = azurerm_key_vault_secret.apim_quotes_function_key.versionless_id
  }
  tags = ["key", "function"]
}

resource "azurerm_api_management_backend" "quotes" {
  name                = azurerm_linux_function_app.quotes.name
  api_management_name = azurerm_api_management.apim.name
  resource_group_name = azurerm_api_management.apim.resource_group_name
  protocol            = "http"
  resource_id         = "https://management.azure.com${azurerm_linux_function_app.quotes.id}"
  url                 = "https://${azurerm_app_service_custom_hostname_binding.quotes_clowa_dev.hostname}"

  credentials {
    header = {
      x-functions-key = "{{${azurerm_api_management_named_value.apim_quotes_function_key.name}}}"
    }
  }
}

################################################################################
# Product: Quotes Silver

moved {
  from = azurerm_api_management_product.quotes
  to   = azurerm_api_management_product.quotes_silver
}

resource "azurerm_api_management_product" "quotes_silver" {
  product_id            = "quotes-silver"
  resource_group_name   = azurerm_api_management.apim.resource_group_name
  api_management_name   = azurerm_api_management.apim.name
  display_name          = "Quotes Silver"
  subscription_required = true
  subscriptions_limit   = 1
  approval_required     = true
  published             = true
}

resource "azurerm_api_management_product_policy" "quotes_silver" {
  product_id          = azurerm_api_management_product.quotes_silver.product_id
  api_management_name = azurerm_api_management_product.quotes_silver.api_management_name
  resource_group_name = azurerm_api_management_product.quotes_silver.resource_group_name
  xml_content         = <<XML
<policies>
    <inbound>
        <base />
        <rate-limit calls="120" renewal-period="60" remaining-calls-header-name="x-ratelimit-remaining" retry-after-header-name="x-ratelimit-retry-after">
          <api id="${azurerm_api_management_api.quotes.name}" name="${azurerm_api_management_api.quotes.display_name}">
            <operation id="post-quote" calls="1" renewal-period="300" name="Write a quote to the database" />
          </api>
        </rate-limit>
    </inbound>
    <outbound>
        <base />
    </outbound>
</policies>
XML
}

################################################################################
# APIs

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

moved {
  from = azurerm_api_management_product_api.quotes_quotes
  to   = azurerm_api_management_product_api.quotes_quotes_silver
}

resource "azurerm_api_management_product_api" "quotes_quotes_silver" {
  resource_group_name = azurerm_api_management.apim.resource_group_name
  api_management_name = azurerm_api_management.apim.name
  api_name            = azurerm_api_management_api.quotes.name
  product_id          = azurerm_api_management_product.quotes_silver.product_id
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
