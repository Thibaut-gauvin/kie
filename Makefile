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

qa: lint test ## Run all QA process

lint: lint.go lint.yaml lint.hadolint lint.chart ## Run all linters

lint.go: ## Lint go source code
	golangci-lint run -v
	@echo ""

lint.yaml: ## Lint yaml file
	yamllint .
	@echo ""

lint.hadolint: ## Lint dockerfiles
	hadolint -- Dockerfile
	@echo ""

lint.chart: ## Lint Helm chart
	ct lint --charts=./charts/kubernetes-image-exporter/
	@echo ""

PKG := "./..."
RUN := ".*"
RED := $(shell tput setaf 1)
GREEN := $(shell tput setaf 2)
BLUE := $(shell tput setaf 4)
RESET := $(shell tput sgr0)

.PHONY: test
test: ## Run Go tests
	@go test -v -race -failfast -coverprofile coverage.output -run $(RUN) $(PKG) | \
        sed 's/RUN/$(BLUE)RUN$(RESET)/g' | \
        sed 's/CONT/$(BLUE)CONT$(RESET)/g' | \
        sed 's/PAUSE/$(BLUE)PAUSE$(RESET)/g' | \
        sed 's/PASS/$(GREEN)PASS$(RESET)/g' | \
        sed 's/FAIL/$(RED)FAIL$(RESET)/g'
	@echo ""

##
## ----------------------
## Development
## ----------------------
##

export KUBECONFIG := test/local/kubeconfig.yml
export HELM_CONFIG_HOME := test/local/helm

start: ## Start project locally (go run ...)
	go run ./cmd/kubernetes-image-exporter serve -l debug -k test/local/kubeconfig.yml

kind.create: ## Create Kind dev cluster
	kind create cluster --config=.kind.yaml
	@echo ""
	kind get clusters
	kubectl cluster-info
	kubectl config get-contexts
	kubectl get node -o wide

kind.delete: ## Delete Kind dev cluster
	kind delete cluster --name dev

kind.provision: ## Setup Kind dev cluster (install prometheus-stack with Helm)
	helm repo update
	@echo ""
	helm upgrade -i \
		prometheus-stack prometheus-community/kube-prometheus-stack \
		--version 57.0.3 \
		--create-namespace \
		--namespace monitoring \
		--values test/local/helm/prometheus-stack/values.yml \
		--debug \
		--wait

app.build:
	docker build -t kie:dev .
	@echo ""
#	trivy image --format spdx-json --output kie.spdx.json kie:dev
#	trivy image --format cyclonedx --output kie.cyclonedx.json kie:dev
#	trivy image --output kie.trivy.json kie:dev

app.deploy: app.build ## Deploy app on Kind dev cluster
	kind load docker-image kie:dev --name dev
	@echo ""
	helm upgrade -i \
		kie ./charts/kubernetes-image-exporter \
		--values test/local/helm/kie_local/value.yml \
		--debug \
		--wait

app.test: app.build ## Install & test Helm chart using Kind dev cluster
	kind load docker-image kie:dev --name dev
	@echo ""
	ct install \
		--all \
		--helm-extra-set-args '--set=image.tag=dev'
