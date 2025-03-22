package tester

import (
	"context"
	"errors"
	"github.com/mark3labs/mcp-go/mcp"
)

type ToolHandler struct {
	Tool    mcp.Tool
	Handler func(context.Context, mcp.CallToolRequest, []string) (*mcp.CallToolResult, error)
}

type MCPClient interface {
	CallTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
	Close() error
}

type MCPClientFactory func(ctx context.Context, args []string) (*mcp.InitializeResult, MCPClient)

func BypassFactory(tools []ToolHandler) MCPClientFactory {
	return func(ctx context.Context, allowedDirs []string) (*mcp.InitializeResult, MCPClient) {
		return &mcp.InitializeResult{}, &bypassClient{
			tools:       tools,
			allowedDirs: allowedDirs,
		}
	}
}

type bypassClient struct {
	tools       []ToolHandler
	allowedDirs []string
}

var toolNotFoundError = errors.New("tool not found")

func (c *bypassClient) CallTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	for _, tool := range c.tools {
		if tool.Tool.Name == req.Params.Name {
			return tool.Handler(ctx, req, c.allowedDirs)
		}
	}
	return nil, toolNotFoundError
}

func (c *bypassClient) Close() error {
	return nil
}
