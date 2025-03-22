package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
)

func defineCreateDirectoryTool() mcp.Tool {
	return mcp.NewTool("create_directory",
		mcp.WithDescription(
			"Create a new directory or ensure a directory exists. Can create multiple "+
				"nested directories in one operation. If the directory already exists, "+
				"this operation will succeed silently. Perfect for setting up directory "+
				"structures for projects or ensuring required paths exist. Only works within allowed directories."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path for the new directory")),
	)
}

func createDirectoryHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	path, ok := req.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}
	validPath, err := validatePath(path, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := os.MkdirAll(validPath, 0755); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully created directory %s", path)), nil
}
