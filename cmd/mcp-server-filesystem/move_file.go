package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
)

func defineMoveFileTool() mcp.Tool {
	return mcp.NewTool("move_file",
		mcp.WithDescription(
			"Move or rename files and directories. Can move files between directories "+
				"and rename them in a single operation. If the destination exists, the "+
				"operation will fail. Works across different directories and can be used "+
				"for simple renaming within the same directory. Both source and destination must be within allowed directories."),
		mcp.WithString("source", mcp.Required(), mcp.Description("Source path")),
		mcp.WithString("destination", mcp.Required(), mcp.Description("Destination path")),
	)
}

func moveFileHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	source, ok := req.Params.Arguments["source"].(string)
	if !ok {
		return mcp.NewToolResultError("source must be a string"), nil
	}
	dest, ok := req.Params.Arguments["destination"].(string)
	if !ok {
		return mcp.NewToolResultError("destination must be a string"), nil
	}
	validSource, err := validatePath(source, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	validDest, err := validatePath(dest, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	// Check if destination exists.
	// Not atomic, but os.Rename cannot generally be expected to be anyway.
	if _, err := os.Stat(validDest); err == nil {
		return mcp.NewToolResultError("Destination already exists"), nil
	}
	if err := os.Rename(validSource, validDest); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully moved %s to %s", source, dest)), nil
}
