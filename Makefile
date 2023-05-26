K8S_NS := terraform-provider-dexidp

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Documentation:

.PHONY: provider-docs
provider-docs: ## Generate provider documentation
	@go generate ./...

##@ Testing:

.PHONY: test
test: ## Run unit tests
	go test -v -count=1 ./...

.PHONY: acc-test
acc-test: install-dex ## Run acceptance tests
	@echo Running tests…
	@scripts/run-acc-test.sh $(K8S_NS)

##@ Install:

GOBIN := $(shell go env GOBIN)
ifeq ($(strip $(GOBIN)),)
    # If GOBIN is empty, set the variable GOBIN to the $GOPATH/bin bin folder
    GOBIN := $(shell go env GOPATH)/bin
endif

.PHONY: install
install: ## Install the provider to $GOBIN
	@echo Installing provider to $(GOBIN)…
	@go install .

.PHONY: setup-dex-helm-repo
setup-dex-helm-repo:
	@helm repo add dex https://charts.dexidp.io

.PHONY: install-dex
install-dex: setup-dex-helm-repo ## Install dex on k8s using helm chart
	@helm upgrade -n $(K8S_NS) --install --create-namespace \
		--values .github/ci/values.yaml --wait \
		dex-test dex/dex
	@echo

.PHONY: uninstall-dex
uninstall-dex: ## Uninstall dex from k8s
	@helm uninstall -n $(K8S_NS) dex-test
	@echo
