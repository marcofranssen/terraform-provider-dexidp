# Below examples expects the Dex gRPC server to be reachable at localhost

provider "dexidp" {
  host = "127.0.0.1:5557"
}

# This example uses a mTLS certificate to authenticate the Dex gRPC server

provider "dexidp" {
  host = "127.0.0.1:5557"

  tls = {
    ca_cert     = file("certs/ca.crt")
    client_cert = file("certs/client.crt")
    client_key  = file("certs/client.key")
  }
}
