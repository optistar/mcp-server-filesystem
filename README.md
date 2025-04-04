# Filesystem MCP Server

This is a Go implementation of the Model Context Protocol (MCP) server for filesystem operations,
based on Anthropic's reference implementation from [modelcontextprotocol/servers](https://github.com/modelcontextprotocol/servers/blob/main/src/filesystem/index.ts).

Significant differences from the reference implementation include:

- The `get_file_info` and `directory_tree` commands return JSON data.
- The `search_files` tool supports gitignore-style exclude patterns.

## Installation

Install the server with the following command:

```shell
go install github.com/optistar/mcp-server-filesystem/cmd/mcp-server-filesystem@latest
```

The executable will be installed in the Go bin directory, which is typically `~/go/bin`.
If this directory is not in your PATH, you may need to add it.

You can also set the `GOBIN` environment variable to change the installation directory.

## Testing

A full test suite is included to ensure the server behaves as expected.

- The tests can be run normally with `go test`.
- A `tester` command packages the test suite as an MCP client that can be used to test the MCP server command, 
  or another MCP server implementation.

## Tools

The server provides the following tools for interacting with the filesystem:

- `create_directory`: Create a new directory or ensure a directory exists.
- `directory_tree`: Get a recursive tree view of files and directories as a JSON structure.
- `edit_file`: Make line-based edits to a text file.
- `get_file_info`: Retrieve detailed metadata about a file or directory.
- `list_allowed_directories`: Returns the list of directories that this server is allowed to access.
- `list_directory`: Get a detailed listing of all files and directories in a specified path.
- `move_file`: Move or rename files and directories.
- `read_file`: Read the complete contents of a file from the file system.
- `read_multiple_files`: Read the contents of multiple files simultaneously.
- `search_files`: Recursively search for files and directories matching a pattern.
- `write_file`: Create a new file or completely overwrite an existing file with new content.

Use the [inspector](https://github.com/modelcontextprotocol/inspector) for full details on each tool.

## Other implementations

This repository is unrelated to [mark3labs/mcp-filesystem-server](https://github.com/mark3labs/mcp-filesystem-server.git).

## License

The code in this repository is distributed under the MIT License.
