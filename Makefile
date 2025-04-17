K8S_NS := terraform-provider-dexidp

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

FORCE: ;

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

HELM_RELEASE_NAME ?= dex

terraform-provider-dexidp:
	@echo Compiling binary…
	@go build .

.PHONY: build ## Build the binary
build: terraform-provider-dexidp FORCE

.PHONY: install
install: ## Install the provider to $GOBIN
	@echo Installing provider to $(GOBIN)…
	@go install .

.PHONY: setup-dex-helm-repo
setup-dex-helm-repo:
	@helm repo add dex https://charts.dexidp.io

.PHONY: install-dex
install-dex: setup-dex-helm-repo $(CA_CERT) certs/server.crt certs/server.key ## Install dex on k8s using helm chart
	@kubectl create namespace $(K8S_NS) || true
	@kubectl -n $(K8S_NS) create configmap dex-tls --from-file=certs/server.crt --from-file=certs/server.key --from-file=$(CA_CERT) || true
	@helm upgrade -n $(K8S_NS) --install --create-namespace \
		--values .github/ci/values.yaml --wait \
		$(HELM_RELEASE_NAME) dex/dex
	@echo

.PHONY: uninstall-dex
uninstall-dex: ## Uninstall dex from k8s
	@helm uninstall -n $(K8S_NS) $(HELM_RELEASE_NAME)
	@kubectl delete namespace $(K8S_NS)
	@echo

##@ Certificates:

CA_CERT := certs/ca.crt
CA_KEY := certs/ca.key

certs/ca.key:
	@echo Generating $@…
	@mkdir -p certs
	@openssl genrsa -out $@ 4096

certs/ca.crt: $(CA_KEY)
	@echo Generating $@…
	@openssl req -x509 -new -nodes -days 1825 -sha256 \
		-config openssl.cnf \
		-key $< -out $@ \
		-subj "/CN=Dex CA"

certs/%.key:
	@echo Generating $@…
	@openssl genrsa -out $@ 2048

certs/server.csr: certs/server.key
	@echo Generating $@…
	@openssl req -new -key $< -out $@ -nodes \
		-config openssl.cnf \
		-subj "/CN=localhost"

certs/server.crt: certs/server.csr $(CA_CERT) $(CA_KEY)
	@echo Generating $@…
	@openssl x509 -req -in $< -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -days 365 \
		-extensions server_cert -extfile openssl.cnf \
		-out $@

certs/client.csr: certs/client.key
	@echo Generating $@…
	@openssl req -new -key $< -out $@ -nodes \
		-config openssl.cnf \
		-subj "/CN=Terraform provider Dex"

certs/client.crt: certs/client.csr $(CA_CERT) $(CA_KEY)
	@echo Generating $@…
	@openssl x509 -req -in $< -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial -days 365 \
		-extensions client_cert -extfile openssl.cnf \
		-out $@

.PHONY: generate-certs
generate-certs: certs/ca.crt certs/server.crt certs/client.crt ## Generate mTLS certificates for tests

.PHONY: clean-certs
clean-certs: ## Remove certs folder so you can regenerate all certificates
	@echo Removing certs…
	@rm -f certs/client.crt certs/server.crt
