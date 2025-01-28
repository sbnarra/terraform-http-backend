terraform {
  backend "http" {
    address        = "http://localhost:9944/states/test/deployment"
    lock_address   = "http://localhost:9944/locks/test/deployment"
    unlock_address = "http://localhost:9944/locks/test/deployment"
    
    username = "user"  # Optional: If AUTH_USERNAME is defined
    password = "pass"  # Optional: If AUTH_PASSWORD is defined
  }
}

resource "null_resource" "example" {}