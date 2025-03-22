package tester

import (
	"encoding/json"
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	Permissions string `json:"permissions"`
	Symlink     string `json:"symlink,omitempty"`
	Size        int64  `json:"size"`
	Created     string `json:"created,omitempty"`
	Modified    string `json:"modified"`
	Changed     string `json:"changed,omitempty"`
	Accessed    string `json:"accessed"`
}

func TestGetFileInfo(t T, f MCPClientFactory) {
	tempDir := t.TempDir()
	_, c := f(t.Context(), []string{tempDir})
	defer c.Close()

	// Create subdirectories for testing
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	emptyDir := filepath.Join(tempDir, "emptydir")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}

	// Create test files with different permissions
	regularFile := filepath.Join(tempDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("Regular file content"), 0644); err != nil {
		t.Fatalf("Failed to create regular file: %v", err)
	}

	executableFile := filepath.Join(tempDir, "executable.sh")
	if err := os.WriteFile(executableFile, []byte("#!/bin/sh\necho 'Hello'"), 0755); err != nil {
		t.Fatalf("Failed to create executable file: %v", err)
	}

	// Create a symlink if supported by the OS
	symlinkFile := filepath.Join(tempDir, "symlink.txt")
	symlinkCreated := false
	if err := os.Symlink(regularFile, symlinkFile); err == nil {
		symlinkCreated = true
	}

	// Test cases
	testCases := []struct {
		name          string
		path          string
		expectedError bool
		checkContent  func(string) bool
	}{
		{
			name:          "Regular file info",
			path:          regularFile,
			expectedError: false,
			checkContent: func(content string) bool {
				var fi FileInfo
				if err := json.Unmarshal([]byte(content), &fi); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				return strings.HasPrefix(fi.Permissions, "-rw") && fi.Size > 0
			},
		},
		{
			name:          "Executable file info",
			path:          executableFile,
			expectedError: false,
			checkContent: func(content string) bool {
				var fi FileInfo
				if err := json.Unmarshal([]byte(content), &fi); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				return strings.HasPrefix(fi.Permissions, "-rwx") && fi.Size > 0
			},
		},
		{
			name:          "Directory info",
			path:          subDir,
			expectedError: false,
			checkContent: func(content string) bool {
				var fi FileInfo
				if err := json.Unmarshal([]byte(content), &fi); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				return strings.HasPrefix(fi.Permissions, "drwx")
			},
		},
		{
			name:          "Non-existent file",
			path:          filepath.Join(tempDir, "nonexistent.txt"),
			expectedError: true,
			checkContent: func(content string) bool {
				return strings.Contains(content, "no such file or directory")
			},
		},
		{
			name:          "Path outside allowed directories",
			path:          "/etc/passwd",
			expectedError: true,
			checkContent: func(content string) bool {
				return strings.Contains(content, "access denied")
			},
		},
		{
			name:          "Path traversal attempt",
			path:          filepath.Join(subDir, "../regular.txt"),
			expectedError: false, // Should resolve to valid path
			checkContent: func(content string) bool {
				var fi FileInfo
				if err := json.Unmarshal([]byte(content), &fi); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				return strings.HasPrefix(fi.Permissions, "-rw") && fi.Size > 0
			},
		},
	}

	// Add symlink test if it was created successfully
	if symlinkCreated {
		testCases = append(testCases, struct {
			name          string
			path          string
			expectedError bool
			checkContent  func(string) bool
		}{
			name:          "Symlink file info",
			path:          symlinkFile,
			expectedError: false,
			checkContent: func(content string) bool {
				var fi FileInfo
				if err := json.Unmarshal([]byte(content), &fi); err != nil {
					t.Fatalf("Failed to unmarshal JSON: %v", err)
				}
				return strings.HasPrefix(fi.Permissions, "L") && fi.Symlink != ""
			},
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create request with path
			req := mcp.CallToolRequest{}
			req.Params.Name = "get_file_info"
			req.Params.Arguments = map[string]interface{}{
				"path": tc.path,
			}

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
