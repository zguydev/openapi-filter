# openapi-filter
> Easily filter an OpenAPI spec to include only specific paths, methods, or components

[![Go Reference](https://pkg.go.dev/badge/github.com/zguydev/openapi-filter.svg)](https://pkg.go.dev/github.com/zguydev/openapi-filter)
[![Go Report Card](https://goreportcard.com/badge/github.com/zguydev/openapi-filter?style=flat-square)](https://goreportcard.com/report/github.com/zguydev/openapi-filter)

## Introduction
`openapi-filter` is a command-line tool designed to help you selectively filter out your OpenAPI 3.0 specification files. This is useful when you need to generate client SDKs or mock servers for a subset of an API, or when you want to exclude internal or deprecated parts of your specification before sharing it.

## Install
```shell
go install github.com/zguydev/openapi-filter@latest
```

## Usage
```go
//go:generate go run github.com/zguydev/openapi-filter openapi.yaml filtered.openapi.yaml --config .openapi-filter.yaml
```

## Features
- **Filter by Paths and Methods**: precisely include only specific API paths and their associated HTTP methods (e.g., keep only `GET /users` and `POST /items`). All referenced components (schemas, parameters, etc.) are automatically included to ensure a valid, self-contained spec (applies only to components referenced by `$ref`).
- **Filter by Components**: externally add specified components to filtered OpenAPI spec.
- **Control Top-Level Elements**: choose whether to include top-level elements:
    - Server definitions (`servers`)
    - Global security requirements (`security`)
    - Tag definitions (`tags`)
    - External documentation objects (`externalDocs`)
- **Easy Filter Configuration**: define your filtering rules in a simple config file: `YAML`, `TOML` and `JSON` formats are supported!

### Filter Configuration

The filter configuration file (e.g., `.openapi-filter.yaml`) specifies what parts of the OpenAPI spec to keep. `YAML`, `TOML` and `JSON` formats are supported. Here's an example `YAML` configuration:

```yaml
# .openapi-filter.yaml

# Tool-specific configurations (optional)
x-openapi-filter:
  logger:
    level: info # Log level (e.g., "debug", "info", "warn", "error")
  loader:
    external_refs_allowed: false # Whether to allow external references

# Keep or discard server information (default: false)
servers: true
# Keep or discard global security definitions (default: false)
security: true
# Keep or discard tag definitions (default: false)
tags: true
# Keep or discard external documentation (default: false)
externalDocs: true

# Specify paths and methods to keep.
# If a path is listed, only the specified methods are kept.
paths:
  /pets: [ post, put ]
  /pet/{petId}/uploadImage: [ post ]
  /user/login: [ get ]
  # Paths not listed here will be removed.

# Specify components to keep.
# Referenced components from kept paths are automatically kept.
components:
  schemas:
    - Pet
    - User
    - Error
  securitySchemes:
    - petstore_auth
  # Components not listed (that are not referenced from kept paths) will be removed.
```

## Examples
Explore ready-to-use examples:

| Example Name           | Description                               | Path                                        |
| ---------------------- | ----------------------------------------- | ------------------------------------------- |
| 🤖 **OpenAI Example**   | Filters for OpenAI API schema             | [`examples/OpenAI`](./examples/OpenAI/)     |
| 🐶 **Petstore Example** | Classic Swagger Petstore demo             | [`examples/petstore`](./examples/petstore/) |
| 🦊 **GitLab Example**   | TOML filter example for GitLab API schema | [`examples/gitlab`](./examples/gitlab/)     |

## License

This project is licensed under the terms of the Apache License 2.0. See the [LICENSE](./LICENSE) file for details.
