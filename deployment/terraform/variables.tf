variable "resource_group_name" {
  description = "Name of the Azure resource group"
  type        = string
  default     = "Portfolio"
}

variable "location" {
  description = "Azure region to deploy resources"
  type        = string
  default     = "westus3"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "mysql_admin_password" {
  description = "MySQL administrator password"
  type        = string
  sensitive   = true
}