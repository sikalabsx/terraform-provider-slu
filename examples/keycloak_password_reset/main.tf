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
variable "realm" {}
variable "username" {}
variable "new_password" {}


provider "slu" {
  keycloak_url            = var.keycloak_url
  keycloak_admin_username = var.keycloak_admin_username
  keycloak_admin_password = var.keycloak_admin_password
}

resource "slu_keycloak_password_reset" "password_reset" {
  realm        = var.realm
  username     = var.username
  new_password = var.new_password
}
