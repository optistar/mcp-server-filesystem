package tester

import (
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

func TestCreateDirectory(t T, f MCPClientFactory) {
	tempDir := t.TempDir()
	_, c := f(t.Context(), []string{tempDir})
	defer c.Close()

	// Test cases
	testCases := []struct {
		name          string
		path          string
		expectedError bool
		checkResult   func(string) bool
		verifyDir     func(string) bool
	}{
		{
			name:          "Create new directory",
			path:          filepath.Join(tempDir, "newdir"),
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully created directory")
			},
			verifyDir: func(path string) bool {
				info, err := os.Stat(path)
				return err == nil && info.IsDir()
			},
		},
		{
			name:          "Create nested directory",
			path:          filepath.Join(tempDir, "parent/child/grandchild"),
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully created directory")
			},
			verifyDir: func(path string) bool {
				info, err := os.Stat(path)
				return err == nil && info.IsDir()
			},
		},
		{
			name:          "Create directory that already exists",
			path:          tempDir, // This already exists
			expectedError: false,   // Should not error, MkdirAll is idempotent
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully created directory")
			},
			verifyDir: func(path string) bool {
				info, err := os.Stat(path)
				return err == nil && info.IsDir()
			},
		},
		{
			name:          "Path outside allowed directories",
			path:          "/etc/newdir",
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "access denied")
			},
			verifyDir: func(path string) bool {
				_, err := os.Stat(path)
				return os.IsNotExist(err) // Directory should not exist
			},
		},
		{
			name:          "Create directory with file name conflict",
			path:          filepath.Join(tempDir, "file-conflict"),
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "not a directory") ||
					strings.Contains(result, "file exists")
			},
			verifyDir: func(path string) bool {
				info, err := os.Stat(path)
				return err == nil && !info.IsDir() // Should be a file, not a directory
			},
		},
	}

	// Create a file to test name conflict
	fileConflictPath := filepath.Join(tempDir, "file-conflict")
	if err := os.WriteFile(fileConflictPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create request with path
			req := mcp.CallToolRequest{}
			req.Params.Name = "create_directory"
			req.Params.Arguments = map[string]interface{}{
				"path": tc.path,
			}

			// Call handler
			result, err := c.CallTool(t.Context(), req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result using helper function
			assertToolResult(t, result, tc.expectedError, tc.checkResult)

			// Verify directory was created or not as expected
			if !tc.verifyDir(tc.path) {
				t.Errorf("Directory verification failed for path: %s", tc.path)
			}
		})
	}
}
