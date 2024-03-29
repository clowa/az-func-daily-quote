terraform {
  # Start with this block commented out to bootstrap using a local terraform state
  # Uncomment after initial deployment, run "terraform init -migrate-state" to migrate the state to the new backend
  #
  # backend "azurerm" {
  #   tenant_id       = "452df130-4ad1-437a-8c0e-9be535aeb732" # Clowa
  #   subscription_id = "0a0a4299-b306-4dad-94de-862e8405fdbe" # teq-free-msdn-sandbox-sub
  #   storage_account_name = ""
  #   container_name       = "prod"
  #   key                  = "terraform.tfstate"
  #   use_azuread_auth     = true
  # }

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.94"
    }
    azapi = {
      source  = "Azure/azapi"
      version = "~> 1.12"
    }
  }
}

provider "azapi" {
  tenant_id       = "452df130-4ad1-437a-8c0e-9be535aeb732" # Clowa
  subscription_id = "0a0a4299-b306-4dad-94de-862e8405fdbe" # teq-free-msdn-sandbox-sub
}

provider "azurerm" {
  tenant_id       = "452df130-4ad1-437a-8c0e-9be535aeb732" # Clowa
  subscription_id = "0a0a4299-b306-4dad-94de-862e8405fdbe" # teq-free-msdn-sandbox-sub

  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}
data "azurerm_subscription" "current" {}

locals {
  region = "West Europe"

  company     = "cwa"
  solution    = "quotes"
  environment = "prod"

  global_prefix = "${local.company}-${local.solution}-${local.environment}"

  tags = {
    datadog-metrics  = true
    solution         = local.solution
    environment      = local.environment
    deploymentMethod = "terraform"
  }
}
