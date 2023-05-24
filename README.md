# Terraform provider to interact with Dex

Using this Terraform provider you can provision [dexidp.io][] clients.

To do so [dexidp.io][] uses a gRPC api which is consumed by this Terraform provider.

## Development

To contribute to this terraform provider you should setup `.terraformrc` with following contents.

*Ensure you put the path to your `go env GOBIN` path.*

```terraform
provider_installation {
  dev_overrides {
    "marcofranssen/dexidp" = "/Users/marco/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Now run `go install .` to install the provider to your `GOBIN` path.

Now you can test the provider using following.

```shell
$ cd examples/provider-install-verification
$ terraform plan
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - marcofranssen/curl in /Users/marco/go/bin
│  - marcofranssen/dexidp in /Users/marco/go/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
```

[dexidp.io]: https://dexidp.io/ "A Federated OpenID Connect Provider"
