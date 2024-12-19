# Terraform Provider for SonarCloud

This repository contains a Terraform provider for managing resources in SonarCloud.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.16+ (to build the provider)

## Building The Provider

1. Clone the repository
   ```sh
   git clone https://github.com/m-yosefpor/terraform-provider-sonarcloud.git
   cd terraform-provider-sonarcloud
   ```

2. Build the provider
   ```sh
   make build
   ```

3. Copy the provider binary to the Terraform plugins directory
   ```sh
   make plugin-init
   ```

## Using The Provider

To use the provider, add it to your Terraform configuration (see main.tf as an example)

## Variables

- `sonarcloud_token`: The SonarCloud token used for authentication.
- `organization`: The SonarCloud organization.

## Resources

- `sonarcloud_project`: Manages a SonarCloud project.
- `sonarcloud_qualitygates_select`: Manages the selection of a quality gate for a project.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature-branch`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature-branch`)
5. Create a new Pull Request

## License

This project is licensed under the Apache 2 License - see the LICENSE file for details.
