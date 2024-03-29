GREEN := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE := $(shell tput -Txterm setaf 7)
RESET := $(shell tput -Txterm sgr0)

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

ifeq ($(command), plan)
TARGET_COMMAND := plan
TARGET_DIR := data-sources
else
TARGET_DIR := resources
TARGET_COMMAND := apply -auto-approve
endif

default: build

.PHONY: help
help:
	@echo "$(YELLOW)Usage:$(RESET) make $(GREEN)[target]$(RESET)"
	@echo ""
	@echo "$(YELLOW)Targets:$(RESET)"
	@echo "  clean           - Clean up"
	@echo "  build $(YELLOW)(default)$(RESET) - Build the provider"
	@echo "  unittest        - Run unit tests"
	@echo "  test            - Run apply and tests the results"

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf ./.bin/
	@rm -rf ./examples/resources/montana_joke/*.tfstate

.PHONY: build
build:
	$(MAKE) clean
	@echo "Building..."
	@go build -o .bin/ ./...
	mv .bin/montana .bin/terraform-provider-montana

.PHONY: unittest
unittest:
	$(MAKE) build
	@echo "Testing..."
	TF_ACC=1 go test -v ./...

.PHONY: test
test:
	$(MAKE) build
	export TF_CLI_CONFIG_FILE="$(ROOT_DIR)/terraform.tfrc" ; cd "$(ROOT_DIR)/examples/$(TARGET_DIR)/montana_joke" && terraform $(TARGET_COMMAND)
