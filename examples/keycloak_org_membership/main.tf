terraform {
  required_providers {
    slu = {
      source = "sikalabsx/slu"
    }
  }
}

variable "keycloak_url" {}
variable "keycloak_admin_username" {}
variable "keycloak_admin_password" {}
variable "realm_id" {}
variable "org_id" {}
variable "username" {}

provider "slu" {
  keycloak_url            = var.keycloak_url
  keycloak_admin_username = var.keycloak_admin_username
  keycloak_admin_password = var.keycloak_admin_password
}

resource "slu_keycloak_org_membership" "org_membership" {
  realm_id = var.realm_id
  org_id   = var.org_id
  username = var.username
}
