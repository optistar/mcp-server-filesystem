package top

import (
	"context"
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
)

func DefineDirectoryTreeTool() mcp.Tool {
	return mcp.NewTool("directory_tree",
		mcp.WithDescription(
			"Get a recursive tree view of files and directories as a JSON structure. "+
				"Each entry includes 'name', 'type' (file/directory), and 'children' for directories. "+
				"Files have no children array, while directories always have a children array (which may be empty). "+
				"Only works within allowed directories."),
		mcp.WithString("path", mcp.Required(), mcp.Description("Root path for the tree")),
		mcp.WithBoolean("pretty", mcp.Description("Format the output with 2-space indentation for readability. "), mcp.DefaultBool(true)),
		mcp.WithNumber("maxDepth", mcp.Description("Maximum depth of recursion"), mcp.DefaultNumber(100)),
	)
}

func DirectoryTreeHandler(_ context.Context, req mcp.CallToolRequest, allowedDirs []string) (*mcp.CallToolResult, error) {
	path, ok := req.Params.Arguments["path"].(string)
	if !ok {
		return mcp.NewToolResultError("path must be a string"), nil
	}
	pretty, _ := req.Params.Arguments["pretty"].(bool)
	maxDepthNum, _ := req.Params.Arguments["maxDepth"].(float64)
	maxDepth := int(maxDepthNum)
	if maxDepth == 0 {
		maxDepth = 100
	} else if maxDepth <= 0 {
		return mcp.NewToolResultError("maxDepth must be a positive integer"), nil
	}
	validPath, err := validatePath(path, allowedDirs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	type TreeEntry struct {
		Name     string      `json:"name"`
		Type     string      `json:"type"`
		Children []TreeEntry `json:"children,omitempty"`
	}

	var buildTree func(string, int) ([]TreeEntry, error)
	buildTree = func(currentPath string, depth int) ([]TreeEntry, error) {
		entries, err := os.ReadDir(currentPath)
		if err != nil {
			return nil, err
		}
		var result []TreeEntry
		for _, entry := range entries {
			entryData := TreeEntry{
				Name: entry.Name(),
				Type: "file",
			}
			if entry.IsDir() {
				entryData.Type = "directory"
				if depth < maxDepth {
					children, err := buildTree(filepath.Join(currentPath, entry.Name()), depth+1)
					if err != nil {
						return nil, err
					}
					entryData.Children = children
				}
			}
			result = append(result, entryData)
		}
		return result, nil
	}

	tree, err := buildTree(validPath, 0)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	// If tree is nil, it will serialize to null; we want [] instead.
	if tree == nil {
		tree = []TreeEntry{}
	}
	indent := ""
	if pretty {
		indent = "  "
	}
	jsonData, err := json.MarshalIndent(tree, "", indent)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}
