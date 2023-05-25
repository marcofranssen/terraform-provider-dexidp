# CONTRIBUTING

Thanks for your interest in this project. In this document we outline a few guidelines to ease the way your contributions flow into this project.

## Commit style

Ensure you have clear and concise commits, written in the present tense. Preferably those commits are transactional so each individual commit results into a releasable provider. See [Kubernetes commit message guidelines](https://www.kubernetes.dev/docs/guide/pull-requests/#commit-message-guidelines) for a more detailed explanation of our approach.

```diff
+ git commit -m "Add new oauth2 resource to the provider"
- git commit -m "Added new oauth2 resource to the provider"
```

## PRs

Stick with one feature per branch. This allows us to make small controlled releases of the module and makes it easy for us to review PRs.

Ensure your branch is rebased on top of `main` before issuing your PR. This to keep a clean Git history and to ensure your changes are working with the latest `main` branch changes.

```shell
git checkout main
git pull
git checkout «your-branch»
git rebase main
```

## Documentation

To update any documentation related to the terraform provider and/or resources you should run following command.

```shell
go generate ./...
```

Preferably you commit these documentation updates in the same commit as your code changes (see [commit style](#commit-style)).

See [this turorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation) how to document certain aspects of the provider.

## Setup development environment

### Terraform

To run the examples it is beneficial to test with the same version as we do (see `.terraform-version`). To easily switch terraform versions you can make use of `tfenv` that easily allows you to use the version defined in `.terraform-version`.

```shell
tfenv use
```

### Code style

This repository has a `.editorconfig` which should help align some file formatting to prevent unnecessary diffs in the code base. Consider installing an `.editorconfig` plugin for your preferred IDE of choice.

### Testing

To contribute to this terraform provider you should setup `.terraformrc` with following contents so you can locally test the provider.

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
