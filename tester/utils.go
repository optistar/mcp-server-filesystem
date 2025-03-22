package tester

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// Helper functions for test assertions
func assertToolResult(t T, result *mcp.CallToolResult, expectedError bool, contentCheck func(string) bool) {
	t.Helper()

	if expectedError {
		if !result.IsError {
			t.Errorf("Expected error but got result: %v", result.Content)
		}
		return
	}

	if result.IsError {
		t.Errorf("Unexpected error in result: %v", result.Content)
		return
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Errorf("Expected text content but got: %v", result.Content)
		return
	}

	if contentCheck != nil && !contentCheck(textContent.Text) {
		t.Errorf("Content check failed. Got: %s", textContent.Text)
	}
}

func expectExactText(expected string) func(string) bool {
	return func(actual string) bool {
		return actual == expected
	}
}
