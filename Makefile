GREEN := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE := $(shell tput -Txterm setaf 7)
RESET := $(shell tput -Txterm sgr0)


default: build

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  clean           - Clean up"
	@echo "  build (default) - Build the provider"
	@echo "  unittest        - Run unit tests"
	@echo "  test            - Run apply and tests the results"

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf .bin/

.PHONY: build
build:
	make clean
	@echo "Building..."
	@go build -o .bin/ ./...
	mv .bin/montana .bin/terraform-provider-montana

.PHONY: unittest
unittest:
	make build
	@echo "Testing..."
	TF_ACC=1 go test -v ./...

.PHONY: test
test:
	make build
	export TF_CLI_CONFIG_FILE="$(ROOT_DIR)/terraform.tfrc" ; cd "$(ROOT_DIR)/examples/resources/montana_palindrome" && terraform apply -auto-approve
