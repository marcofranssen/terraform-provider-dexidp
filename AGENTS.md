# AGENTS.md

## Build & Test
- `go build .` - build binary
- `go install .` - install provider to `$GOBIN`
- `make test` - run unit tests only
- `make acc-test` - runs acceptance tests (requires k8s + Dex helm chart, defined in Makefile)

## Docs
- `go generate ./...` - regenerate provider documentation (not `go doc`)

## Provider Config
- Provider address: `hashicorp.com/marcofranssen/dexidp`
- Framework: `terraform-plugin-framework`

## Commit Style
- Commits in present tense
- One feature per branch
- Rebase on main before PR
- Include docs updates in same commit as code changes

## Local Testing
- Set up `.terraformrc` with dev_overrides pointing to your `$GOBIN`
- Use `tfenv use` to match `.terraform-version`