package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"strings"
)

func defineListAllowedDirectoriesTool() mcp.Tool {
	return mcp.NewTool("list_allowed_directories",
		mcp.WithDescription(
			"Returns the list of directories that this server is allowed to access. "+
				"Use this to understand which directories are available before trying to access files."),
	)
}

func listAllowedDirectoriesHandler(_ context.Context, _ mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(fmt.Sprintf("Allowed directories:\n%s", strings.Join(allowedDirs, "\n"))), nil
}
