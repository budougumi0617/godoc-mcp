# godoc-mcp

## Overview

godoc-mcp is a server that provides information about Go project packages, types, functions, constants, and variables via the Multi-Command Protocol (MCP). It allows you to flexibly retrieve and utilize Go code documentation and structure from external tools.

## Main Features

- Retrieve a list of Go packages
- List exported structs, functions, and methods in a package
- Get detailed information about structs (fields, methods, comments)
- Get detailed information about functions and methods (signature, comments, examples)
- Get detailed information about constants and variables in a package

## Installation

Go 1.24.2 or later is required.

```sh
make build
```

## Usage

### Start the Server

```sh
./godoc-mcp -root <root directory of your Go project>
```

- Use the `-root` option to specify the root directory of the Go project to analyze.
- If omitted, the current directory or the environment variable `GODOC_MCP_ROOT_DIR` will be used.

### Using as an MCP Tool

You can use the following tools from an MCP client:

- `golang_list_packages`: Get a list of packages and their comments
- `golang_inspect_package`: List structs, functions, and methods in a package
- `golang_get_struct_doc`: Get detailed information about a struct
- `golang_get_func_doc`: Get detailed information about a function
- `golang_get_method_doc`: Get detailed information about a struct method
- `golang_get_const_and_var_doc`: Get detailed information about constants and variables

#### Example: mcp settings for Roo Code

```json
{
  "mcpServers": {
    "godoc-mcp": {
      "type": "stdio",
      "command": "godoc-mcp",
      "env": {
        "GODOC_MCP_ROOT_DIR": "${env:HOME}/go/src/github.com/golang/tools/internal",
        "GOPATH": "${env:HOME}/go",
        "HOME": "${env:HOME}",
        "GOCACHE": "${env:HOME}/Library/Caches/go-build"
      },
    }
  }
}
```

## Test

```sh
go test ./...
```

## Dependencies

- github.com/ktr0731/go-mcp
- golang.org/x/exp/jsonrpc2
- golang.org/x/tools

## Environment Variables

- `GODOC_MCP_ROOT_DIR`: Root directory of the Go project to analyze

## License

MIT License