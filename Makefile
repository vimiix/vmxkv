APP=vmxkv
BLD_DIR=dist

##@ Targets

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make <target>\033[36m\033[0m\n"} /^[1-9a-zA-Z_- ]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[0m%s:\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

clean: ## Clean all temp files
	@echo ">>> cleaning..."
	@rm -rf run

fmt: ## Go fmt source files
	@echo ">>> formatting..."
	@go fmt ./...

build: clean fmt ## build vmxkv binary file
	go build -o $(BLD_DIR)/$(APP) cmd/vmxkvd/*.go

benchmark: ## run benchmark test
	go test -bench=. ./server

changelog: ## Generate CHANGELOG.md
	@git-chglog -o CHANGELOG.md

.PHONY: help clean fmt changelog build benchmark