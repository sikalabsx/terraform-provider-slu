terraform {
  required_providers {
    slu = {
      source = "sikalabsx/slu"
    }
  }
}


resource "slu_random_password" "password" {
  count = 5
}

output "passwords" {
  value     = slu_random_password.password[*].result
  sensitive = true
}

output "passwords_nonsensitive" {
  value = [for password in slu_random_password.password : nonsensitive(password.result)]
}
