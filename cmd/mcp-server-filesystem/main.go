// Based on https://github.com/modelcontextprotocol/servers/blob/main/src/filesystem/index.ts
// This server provides a secure filesystem API that restricts access to a set of allowed directories.
// It supports reading, writing, moving, and listing files and directories, as well as searching for files by pattern.

package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	top "github.com/optistar/mcp-server-filesystem"
	"os"
	"path/filepath"
)

func main() {
	// Check command line arguments
	if len(os.Args) < 1 {
		fmt.Println("Usage: mcp-server-filesystem [<allowed-directory> ...]")
		os.Exit(1)
	}

	// Normalize allowed directories
	allowedDirectories := make([]string, 0, len(os.Args)-1)
	for _, dir := range os.Args[1:] {
		absPath, err := filepath.Abs(top.ExpandHome(dir))
		if err != nil {
			fmt.Printf("Error resolving path %s: %v\n", dir, err)
			os.Exit(1)
		}
		allowedDirectories = append(allowedDirectories, absPath)
	}

	// Create MCP server
	s := server.NewMCPServer(
		"secure-filesystem-server",
		"0.2.0",
	)

	// Register tools with handlers
	for _, t := range top.Tools {
		handler := t.Handler // Capture in closure
		s.AddTool(t.Tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handler(ctx, req, allowedDirectories)
		})
	}

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
