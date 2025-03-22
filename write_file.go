package top

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
)

func DefineWriteFileTool() mcp.Tool {
	return mcp.NewTool("write_file",
		mcp.WithDescription(
			"Create a new file or completely overwrite an existing file with new content. "+
				"Use with caution as it will overwrite existing files without warning. "+
				"Handles text content with proper encoding. Only works within allowed directories."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the file")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Content to write")),
	)
}

func WriteFileHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	path, ok := req.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}
	content, ok := req.Params.Arguments["content"].(string)
	if !ok {
		return mcp.NewToolResultError("content must be a string"), nil
	}
	validPath, err := validatePath(path, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := os.WriteFile(validPath, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully wrote to %s", path)), nil
}
