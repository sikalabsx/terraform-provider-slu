terraform {
  required_providers {
    slu = {
      source = "sikalabsx/slu"
    }
  }
}

variable "smtp_host" {}
variable "smtp_port" {}
variable "smtp_user" {}
variable "smtp_password" {}
variable "from" {}
variable "to" {}

provider "slu" {
  smtp_host     = var.smtp_host
  smtp_port     = var.smtp_port
  smtp_user     = var.smtp_user
  smtp_password = var.smtp_password
}

resource "slu_mail_send" "example" {
  from    = var.from
  to      = var.to
  subject = "TEST sikalabsx/slu send mail"
  message = <<EOT
TEST

Hi,
sikalabsx/slu send mail from Terraform.

O.
EOT
}
