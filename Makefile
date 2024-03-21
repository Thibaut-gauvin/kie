##-----------------------
## Available make targets
##-----------------------
##

default: help
help: ## Display this message
	@grep -E '(^[a-zA-Z0-9_.-]+:.*?##.*$$)|(^##)' Makefile | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | \
		sed -e 's/\[32m##/[33m/'

artifact: ## Generate binary in dist folder
	@goreleaser build --clean --snapshot --single-target

##
## ----------------------
## Q.A
## ----------------------
##

qa: lint lint.yaml lint.chart test ## Run all QA process

lint: ## Lint source code
	@golangci-lint run -v

lint.yaml: ## Lint yaml file
	@yamllint .

lint.chart: ## Lint Helm chart
	ct lint --charts=./charts/kubernetes-image-exporter/

PKG := "./..."
RUN := ".*"
RED := $(shell tput setaf 1)
GREEN := $(shell tput setaf 2)
BLUE := $(shell tput setaf 4)
RESET := $(shell tput sgr0)

.PHONY: test
test: ## Run tests
	@go test -v -race -failfast -coverprofile coverage.output -run $(RUN) $(PKG) | \
        sed 's/RUN/$(BLUE)RUN$(RESET)/g' | \
        sed 's/CONT/$(BLUE)CONT$(RESET)/g' | \
        sed 's/PAUSE/$(BLUE)PAUSE$(RESET)/g' | \
        sed 's/PASS/$(GREEN)PASS$(RESET)/g' | \
        sed 's/FAIL/$(RED)FAIL$(RESET)/g'

##
## ----------------------
## Development
## ----------------------
##

export KUBECONFIG := test/local/kubeconfig.yml
export HELM_CONFIG_HOME := test/local/helm

start: ## Start project locally
	go run ./cmd/kubernetes-image-exporter serve -l debug

kind.create: ## Create Kind dev cluster
	kind create cluster --config=.kind.yaml
	kind get clusters
	kubectl cluster-info
	kubectl config get-contexts
	kubectl get node -o wide

kind.provision: ## Setup Kind dev cluster (install prometheus-stack with Helm)
	helm repo update && \
	helm upgrade -i \
		prometheus-stack prometheus-community/kube-prometheus-stack \
		--version 57.0.3 \
		--create-namespace \
		--namespace monitoring \
		--values test/local/helm/prometheus-stack/values.yml \
		--debug \
		--wait

kind.delete: ## Delete Kind dev cluster
	kind delete cluster --name dev

app.deploy: ## Deploy app on Kind dev cluster
	docker build -t kie:dev .
	kind load docker-image kie:dev --name dev
	helm upgrade -i \
		kie ./charts/kubernetes-image-exporter \
		--values test/local/helm/kie_local/value.yml \
		--debug \
		--wait
