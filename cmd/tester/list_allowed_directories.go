package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"strings"
)

func listAllowedDirectoriesTester(t T) {
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()

	// Test cases
	testCases := []struct {
		name          string
		allowedDirs   []string
		expectedError bool
		checkContent  func(string) bool
	}{
		{
			name:          "Empty allowed directories",
			allowedDirs:   []string{},
			expectedError: false,
			checkContent: func(content string) bool {
				return content == "Allowed directories:\n"
			},
		},
		{
			name:          "Single allowed directory",
			allowedDirs:   []string{tempDir1},
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Allowed directories:") &&
					strings.Contains(content, tempDir1) &&
					!strings.Contains(content, tempDir2)
			},
		},
		{
			name:          "Multiple allowed directories",
			allowedDirs:   []string{tempDir1, tempDir2},
			expectedError: false,
			checkContent: func(content string) bool {
				return strings.Contains(content, "Allowed directories:") &&
					strings.Contains(content, tempDir1) &&
					strings.Contains(content, tempDir2)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create MCP client
			c := GetMCPClient(t.Context(), tc.allowedDirs)
			defer c.Close()

			// Create empty request (no arguments needed for this handler)
			req := mcp.CallToolRequest{}
			req.Params.Name = "list_allowed_directories"

			// Call handler
			result, err := c.CallTool(t.Context(), req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result using helper function
			assertToolResult(t, result, tc.expectedError, tc.checkContent)
		})
	}
}
