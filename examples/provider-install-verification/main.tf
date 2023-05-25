terraform {
  required_providers {
    dexidp = {
      source = "marcofranssen/dexidp"
    }
  }
}

provider "dexidp" {
  host = "127.0.0.1:5557"
}

resource "dexidp_client" "my_oidc_client" {
  name      = "Awesome Test Client"
  client_id = "Awesome"
  secret    = "T0pS3cr3t"
  redirect_uris = [
    "https://openidconnect.net/callback",
  ]
}
