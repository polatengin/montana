# Montana

Makefile: This file is used to automate common tasks in the project. It includes targets for cleaning the project, building the provider, running unit tests, and applying the Terraform configuration.

main.go: This is the entry point of the application. It sets up and starts the Terraform provider server. The provider server is configured to run in debug mode if the -debug flag is passed when starting the application.

go.mod: This file is used by Go's dependency management system. It lists the modules that the project depends on.

tools/tools.go: This file is used to track tool dependencies of the project. These are dependencies that are not directly used by your code, but are needed for tasks like generating documentation.

internal/provider: This directory contains the implementation of the Terraform provider. It includes the provider configuration and the implementation of the data source and resource.

examples: This directory contains example Terraform configurations that use the provider.

data: This directory contains data files used by the provider.
