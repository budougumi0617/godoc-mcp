# Go parameters
BINARY_NAME=godoc-mcp
CMD_SERVER_PATH=./cmd/server
CMD_MCPGEN_PATH=./cmd/mcpgen

# Build the application
build:
	go build -o $(BINARY_NAME) $(CMD_SERVER_PATH)

# Clean the binary
clean:
	rm -f $(BINARY_NAME)

.PHONY: build clean