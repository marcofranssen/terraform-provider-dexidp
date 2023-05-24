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
