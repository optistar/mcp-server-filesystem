package main

import (
	"context"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPClientFactory func(ctx context.Context, args []string) (*mcp.InitializeResult, client.MCPClient)

var contextKey = struct{}{}

func WithMCPClientFactory(ctx context.Context, factory MCPClientFactory) context.Context {
	return context.WithValue(ctx, contextKey, factory)
}

func GetMCPClient(ctx context.Context, allowedDirs []string) client.MCPClient {
	if factory, ok := ctx.Value(contextKey).(MCPClientFactory); ok {
		_, c := factory(ctx, allowedDirs)
		return c
	}
	return nil
}
