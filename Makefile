# Variables
PROVIDER_NAME       := sonarcloud
PROVIDER_PARENT     := m-yosefpor
NAMESPACE           := local
VERSION             := 1.0.0
OS_ARCH             := darwin_arm64
PLUGIN_DIR          := ~/.terraform.d/plugins/$(NAMESPACE)/$(PROVIDER_PARENT)/$(PROVIDER_NAME)/$(VERSION)/$(OS_ARCH)
BUILD_DIR           := ./bin
SOURCE_DIR          := .
BINARY_NAME         := terraform-provider-$(PROVIDER_NAME)

# Default target
.PHONY: all
all: build plugin-init plugin-cleanup terraform-init terraform-plan

# Build the provider binary
.PHONY: build
build:
	@echo "Building the Terraform provider binary for macOS ARM64..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

# Create the plugin directory and copy the binary
.PHONY: plugin-init
plugin-init: build
	@echo "Creating the plugin directory and copying the binary..."
	mkdir -p $(PLUGIN_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(PLUGIN_DIR)/
	chmod +x $(PLUGIN_DIR)/$(BINARY_NAME)
	@echo "Plugin binary is ready at $(PLUGIN_DIR)/$(BINARY_NAME)"


.PHONY: plugin-cleanup
plugin-cleanup:
	@echo "Cleaning up the plugin directory..."
	rm -rf ./terraform.lock.hcl
	@echo "Plugin directory cleaned up."

# Terraform initialization
.PHONY: terraform-init
terraform-init:
	@echo "Initializing Terraform..."
	terraform init

# Terraform validation
.PHONY: terraform-validate
terraform-validate:
	@echo "Validating Terraform configuration..."
	terraform validate

# Terraform plan
.PHONY: terraform-plan
terraform-plan: terraform-init
	@echo "Planning Terraform changes..."
	terraform plan

# Terraform apply
.PHONY: terraform-apply
terraform-apply: terraform-init
	@echo "Applying Terraform changes..."
	terraform apply -auto-approve

# Terraform destroy
.PHONY: terraform-destroy
terraform-destroy:
	@echo "Destroying Terraform resources..."
	terraform destroy -auto-approve

# Clean up
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "Cleaned up."

