package top

import (
	"context"
	"encoding/json"
	"github.com/djherbis/times"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"time"
)

func DefineGetFileInfoTool() mcp.Tool {
	return mcp.NewTool("get_file_info",
		mcp.WithDescription(
			"Retrieve detailed metadata about a file or directory. Returns comprehensive "+
				"information including size, creation time, last modified time, permissions, "+
				"and type. This tool is perfect for understanding file characteristics "+
				"without reading the actual content. Only works within allowed directories."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to query")),
	)
}

type FileInfo struct {
	Permissions string `json:"permissions"`
	Symlink     string `json:"symlink,omitempty"`
	Size        int64  `json:"size"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified"`
	Changed     string `json:"changed,omitempty"`
	Accessed    string `json:"accessed"`
}

func GetFileInfoHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	path, ok := req.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}
	validPath, err := validatePath(path, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	info, err := os.Lstat(validPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	t, err := times.Stat(validPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fileStats := FileInfo{
		Permissions: info.Mode().String(),
		Size:        info.Size(),
		Modified:    t.ModTime().Format(time.RFC3339),
		Accessed:    t.AccessTime().Format(time.RFC3339),
	}
	if isSymlink(info) {
		linkTarget, err := os.Readlink(validPath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		fileStats.Symlink = linkTarget
	}
	if t.HasBirthTime() {
		fileStats.Created = t.BirthTime().Format(time.RFC3339)
	}
	if t.HasChangeTime() {
		fileStats.Changed = t.ChangeTime().Format(time.RFC3339)
	}
	jsonData, err := json.Marshal(fileStats)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}
