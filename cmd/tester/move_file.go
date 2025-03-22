package main

import (
	"github.com/mark3labs/mcp-go/mcp"
	"os"
	"path/filepath"
	"strings"
)

func moveFileTester(t T) {
	tempDir := t.TempDir()
	c := GetMCPClient(t.Context(), []string{tempDir})
	defer c.Close()

	// Create subdirectories for testing
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	emptyDir := filepath.Join(tempDir, "empty")

	for _, dir := range []string{srcDir, destDir, emptyDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create test files
	testFiles := map[string]string{
		filepath.Join(srcDir, "file1.txt"): "Content of file 1",
		filepath.Join(srcDir, "file2.txt"): "Content of file 2",
		filepath.Join(tempDir, "root.txt"): "Root file content",
	}

	for path, content := range testFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	// Test cases
	testCases := []struct {
		name          string
		source        string
		destination   string
		expectedError bool
		checkResult   func(string) bool
		verifyMove    func() bool
	}{
		{
			name:          "Move file to new location",
			source:        filepath.Join(srcDir, "file1.txt"),
			destination:   filepath.Join(destDir, "moved_file1.txt"),
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully moved")
			},
			verifyMove: func() bool {
				// Source should not exist, destination should exist with correct content
				_, srcErr := os.Stat(filepath.Join(srcDir, "file1.txt"))
				destContent, destErr := os.ReadFile(filepath.Join(destDir, "moved_file1.txt"))
				return os.IsNotExist(srcErr) &&
					destErr == nil &&
					string(destContent) == "Content of file 1"
			},
		},
		{
			name:          "Rename file in same directory",
			source:        filepath.Join(srcDir, "file2.txt"),
			destination:   filepath.Join(srcDir, "renamed_file2.txt"),
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully moved")
			},
			verifyMove: func() bool {
				_, srcErr := os.Stat(filepath.Join(srcDir, "file2.txt"))
				destContent, destErr := os.ReadFile(filepath.Join(srcDir, "renamed_file2.txt"))
				return os.IsNotExist(srcErr) &&
					destErr == nil &&
					string(destContent) == "Content of file 2"
			},
		},
		{
			name:          "Move directory",
			source:        emptyDir,
			destination:   filepath.Join(destDir, "moved_empty"),
			expectedError: false,
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully moved")
			},
			verifyMove: func() bool {
				_, srcErr := os.Stat(emptyDir)
				destInfo, destErr := os.Stat(filepath.Join(destDir, "moved_empty"))
				return os.IsNotExist(srcErr) &&
					destErr == nil &&
					destInfo.IsDir()
			},
		},
		{
			name:          "Source does not exist",
			source:        filepath.Join(srcDir, "nonexistent.txt"),
			destination:   filepath.Join(destDir, "should_not_create.txt"),
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "no such file or directory")
			},
			verifyMove: func() bool {
				_, destErr := os.Stat(filepath.Join(destDir, "should_not_create.txt"))
				return os.IsNotExist(destErr)
			},
		},
		{
			name:          "Destination already exists",
			source:        filepath.Join(tempDir, "root.txt"),
			destination:   filepath.Join(destDir, "moved_file1.txt"), // This was created in first test
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "file exists") ||
					strings.Contains(result, "already exists")
			},
			verifyMove: func() bool {
				// Source should still exist, destination should not be overwritten
				srcContent, srcErr := os.ReadFile(filepath.Join(tempDir, "root.txt"))
				destContent, destErr := os.ReadFile(filepath.Join(destDir, "moved_file1.txt"))
				return srcErr == nil &&
					destErr == nil &&
					string(srcContent) == "Root file content" &&
					string(destContent) == "Content of file 1" // Original content
			},
		},
		{
			name:          "Source outside allowed directories",
			source:        "/etc/passwd",
			destination:   filepath.Join(destDir, "passwd"),
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "access denied")
			},
			verifyMove: func() bool {
				_, destErr := os.Stat(filepath.Join(destDir, "passwd"))
				return os.IsNotExist(destErr)
			},
		},
		{
			name:          "Destination outside allowed directories",
			source:        filepath.Join(tempDir, "root.txt"),
			destination:   "/tmp/should_not_create.txt",
			expectedError: true,
			checkResult: func(result string) bool {
				return strings.Contains(result, "access denied")
			},
			verifyMove: func() bool {
				// Source should still exist
				_, srcErr := os.Stat(filepath.Join(tempDir, "root.txt"))
				return srcErr == nil
			},
		},
		{
			name:          "Path traversal attempt in source",
			source:        filepath.Join(srcDir, "../root.txt"),
			destination:   filepath.Join(destDir, "traversal_test.txt"),
			expectedError: false, // Should resolve to valid path
			checkResult: func(result string) bool {
				return strings.Contains(result, "Successfully moved")
			},
			verifyMove: func() bool {
				_, srcErr := os.Stat(filepath.Join(tempDir, "root.txt"))
				destContent, destErr := os.ReadFile(filepath.Join(destDir, "traversal_test.txt"))
				return os.IsNotExist(srcErr) &&
					destErr == nil &&
					string(destContent) == "Root file content"
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t T) {
			// Create request with source and destination
			req := mcp.CallToolRequest{}
			req.Params.Name = "move_file"
			req.Params.Arguments = map[string]interface{}{
				"source":      tc.source,
				"destination": tc.destination,
			}

			// Call handler
			result, err := c.CallTool(t.Context(), req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result using helper function
			assertToolResult(t, result, tc.expectedError, tc.checkResult)

			// Verify move operation
			if !tc.verifyMove() {
				t.Errorf("Move verification failed for %s -> %s", tc.source, tc.destination)
			}
		})
	}
}
