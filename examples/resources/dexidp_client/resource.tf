resource "dexidp_client" "my_oidc_client" {
  name      = "Awesome Test Client"
  client_id = "awesome-client"
  secret    = "T0pS3cr3t"
  redirect_uris = [
    "https://openidconnect.net/callback",
  ]
}
