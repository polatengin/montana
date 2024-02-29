default: build

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
