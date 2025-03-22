package main

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
)

// Tool definitions
func defineReadFileTool() mcp.Tool {
	return mcp.NewTool("read_file",
		mcp.WithDescription(
			"Read the complete contents of a file from the file system. "+
				"Handles various text encodings and provides detailed error messages "+
				"if the file cannot be read. Use this tool when you need to examine "+
				"the contents of a single file. Only works within allowed directories."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the file")),
	)
}

// Tool handlers
func readFileHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	path, ok := req.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}
	validPath, err := validatePath(path, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	content, err := os.ReadFile(validPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(content)), nil
}
