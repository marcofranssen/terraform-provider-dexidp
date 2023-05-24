# Provider install verification

Using this example you can test the provider implementation.

## Prerequisites

To contribute to this Terraform provider you should setup `.terraformrc` with following contents.

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

## Use

```shell
terraform apply
```

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

No providers.

## Modules

No modules.

## Resources

No resources.

## Inputs

No inputs.

## Outputs

No outputs.
<!-- END_TF_DOCS -->
