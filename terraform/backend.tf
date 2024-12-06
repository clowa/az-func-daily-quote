locals {
  iac_backend_state_containers = toset([
    "prod"
  ])

  iac_backend_state_owner_oids = toset([
    data.azurerm_client_config.current.object_id
  ])

  iac_backend_tags = merge(local.tags, {
    application = "iac-backend"
  })
}

################################################################################
# Resource Group

resource "azurerm_resource_group" "iac" {
  name     = "${local.global_prefix}-iac-rg"
  location = local.region
  tags     = local.iac_backend_tags

  lifecycle {
    prevent_destroy = true
  }
}


################################################################################
# Storage Account

#trivy:ignore:AVD-AZU-0011
resource "azurerm_storage_account" "iac" {
  name                            = replace(lower("${local.global_prefix}-iac-st"), "-", "")
  resource_group_name             = azurerm_resource_group.iac.name
  location                        = azurerm_resource_group.iac.location
  tags                            = local.iac_backend_tags
  account_tier                    = "Standard"
  account_replication_type        = "ZRS"
  https_traffic_only_enabled      = true
  allow_nested_items_to_be_public = false

  blob_properties {
    versioning_enabled = true
    delete_retention_policy {
      days = 30
    }
    container_delete_retention_policy {
      days = 30
    }
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "azurerm_role_assignment" "iac_storage_blob_data_owner" {
  for_each = local.iac_backend_state_owner_oids

  role_definition_name = "Storage Blob Data Owner"
  scope                = azurerm_storage_account.iac.id
  principal_id         = each.key
}


################################################################################
# Terraform State Backend Blob Containers

resource "azurerm_storage_container" "iac" {
  for_each = local.iac_backend_state_containers

  name                  = each.key
  storage_account_id    = azurerm_storage_account.iac.id
  container_access_type = "private"

  lifecycle {
    prevent_destroy = true
  }
}


################################################################################
# Terraform State Blob Lifecycle Policy

resource "azurerm_storage_management_policy" "iac_blob_lifecycle" {
  storage_account_id = azurerm_storage_account.iac.id

  rule {
    name    = "tfstate-blob-version-expiration"
    enabled = true
    filters {
      blob_types = ["blockBlob"]
    }
    actions {
      version {
        delete_after_days_since_creation = 7
      }
    }
  }
}

################################################################################
# User Assigned Managed Identity (Function App)

resource "azurerm_user_assigned_identity" "func" {
  name                = "${local.global_prefix}-func-id"
  resource_group_name = azurerm_resource_group.iac.name
  location            = azurerm_resource_group.iac.location
  tags                = local.tags
}

resource "azurerm_federated_identity_credential" "func_main" {
  parent_id           = azurerm_user_assigned_identity.func.id
  resource_group_name = azurerm_resource_group.iac.name
  name                = "iac"
  issuer              = "https://token.actions.githubusercontent.com"
  subject             = "repo:clowa/az-func-daily-quote:ref:refs/heads/main"
  audience            = ["api://AzureADTokenExchange"]
}

resource "azurerm_role_assignment" "func_website_contributor" {
  scope                = azurerm_resource_group.quotes.id
  role_definition_name = "Website Contributor"
  principal_id         = azurerm_user_assigned_identity.func.principal_id
}
