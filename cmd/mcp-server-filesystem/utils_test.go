package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tempDir := t.TempDir()

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Define allowed directories
	allowedDirs := []string{tempDir}

	// Helper function for path validation assertions
	assertPathValidation := func(t *testing.T, path string, expectedError bool) {
		t.Helper()
		validPath, err := validatePath(path, allowedDirs)

		if expectedError {
			if err == nil {
				t.Errorf("Expected error but got none, path: %s, validPath: %s", path, validPath)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}
	}

	// Test cases
	tests := []struct {
		name          string
		path          string
		expectedError bool
	}{
		{
			name:          "Valid file in allowed directory",
			path:          testFile,
			expectedError: false,
		},
		{
			name:          "Valid subdirectory",
			path:          subDir,
			expectedError: false,
		},
		{
			name:          "Non-existent file in allowed directory",
			path:          filepath.Join(tempDir, "nonexistent.txt"),
			expectedError: false, // Should be allowed as parent dir is valid
		},
		{
			name:          "Path outside allowed directories",
			path:          "/etc/passwd",
			expectedError: true,
		},
		{
			name:          "Path traversal attempt",
			path:          filepath.Join(subDir, "../test.txt"),
			expectedError: false, // Should resolve to valid path
		},
		{
			name:          "Path traversal outside allowed dirs",
			path:          filepath.Join(tempDir, "../../../etc/passwd"),
			expectedError: true,
		},
		{
			name:          "Non-existent nested path",
			path:          filepath.Join(tempDir, "nonexistent/nested/path"),
			expectedError: false, // Should be allowed as parent dir is valid
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertPathValidation(t, tc.path, tc.expectedError)
		})
	}
}
