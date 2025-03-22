package top

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestApplyFileEdits_ExactMatch(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	originalContent := "This is line one.\nThis is line two.\nThis is line three.\n"
	err := os.WriteFile(testFilePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Define an edit with exact match
	edits := []Edit{
		{
			OldText: "This is line two.",
			NewText: "This is line TWO - modified.",
		},
	}

	// Apply edits in dry-run mode
	diff, err := applyFileEdits("test.txt", testFilePath, edits, true)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify diff contains expected changes
	if !strings.Contains(diff, "-This is line two.") {
		t.Errorf("Diff doesn't contain removed line: %s", diff)
	}
	if !strings.Contains(diff, "+This is line TWO - modified.") {
		t.Errorf("Diff doesn't contain added line: %s", diff)
	}

	// Verify file wasn't actually changed
	contentAfterDryRun, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after dry run: %v", err)
	}
	if string(contentAfterDryRun) != originalContent {
		t.Errorf("File was modified during dry run")
	}

	// Apply edits for real
	_, err = applyFileEdits("test.txt", testFilePath, edits, false)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify file was actually changed
	contentAfterEdit, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after edit: %v", err)
	}
	expectedContent := "This is line one.\nThis is line TWO - modified.\nThis is line three.\n"
	if string(contentAfterEdit) != expectedContent {
		t.Errorf("File content doesn't match expected.\nGot: %s\nWant: %s",
			string(contentAfterEdit), expectedContent)
	}
}

func TestApplyFileEdits_NoMatch(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	originalContent := "This is line one.\nThis is line two.\nThis is line three.\n"
	err := os.WriteFile(testFilePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Define an edit with no match
	edits := []Edit{
		{
			OldText: "This text doesn't exist in the file.",
			NewText: "Replacement text.",
		},
	}

	// Apply edits in dry-run mode
	_, err = applyFileEdits("test.txt", testFilePath, edits, true)
	if err == nil {
		t.Errorf("Expected error for non-matching text, but got none")
	}
}

func TestApplyFileEdits_WhitespaceIndentation(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test file with indentation
	testFilePath := filepath.Join(tempDir, "test.txt")
	originalContent := "function test() {\n    if (condition) {\n        doSomething();\n    }\n}\n"
	err := os.WriteFile(testFilePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Define an edit with whitespace differences
	edits := []Edit{
		{
			OldText: "if (condition) {\n        doSomething();\n    }",
			NewText: "if (condition) {\n        doSomethingElse();\n        doAnotherThing();\n    }",
		},
	}

	// Apply edits in dry-run mode
	diff, err := applyFileEdits("test.txt", testFilePath, edits, true)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify diff contains expected changes
	if !strings.Contains(diff, "-        doSomething();") {
		t.Errorf("Diff doesn't contain removed line: %s", diff)
	}
	if !strings.Contains(diff, "+        doSomethingElse();") {
		t.Errorf("Diff doesn't contain first added line: %s", diff)
	}
	if !strings.Contains(diff, "+        doAnotherThing();") {
		t.Errorf("Diff doesn't contain second added line: %s", diff)
	}

	// Apply edits for real
	_, err = applyFileEdits("test.txt", testFilePath, edits, false)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify file was actually changed
	contentAfterEdit, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after edit: %v", err)
	}
	expectedContent := "function test() {\n    if (condition) {\n        doSomethingElse();\n        doAnotherThing();\n    }\n}\n"
	if string(contentAfterEdit) != expectedContent {
		t.Errorf("File content doesn't match expected.\nGot: %s\nWant: %s",
			string(contentAfterEdit), expectedContent)
	}
}

func TestApplyFileEdits_RelativeIndentation(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test file with mixed indentation
	testFilePath := filepath.Join(tempDir, "test.txt")
	originalContent := "class Example {\n  constructor() {\n    this.value = 0;\n  }\n  method() {\n    // do something\n  }\n}\n"
	err := os.WriteFile(testFilePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Define an edit with different relative indentation
	edits := []Edit{
		{
			OldText: "method() {\n    // do something\n  }",
			NewText: "method() {\n      // increased indent\n      return this.value;\n  }",
		},
	}

	// Apply edits in dry-run mode
	diff, err := applyFileEdits("test.txt", testFilePath, edits, true)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify diff contains expected changes
	if !strings.Contains(diff, "-    // do something") {
		t.Errorf("Diff doesn't contain removed line: %s", diff)
	}
	if !strings.Contains(diff, "+      // increased indent") {
		t.Errorf("Diff doesn't contain first added line: %s", diff)
	}
	if !strings.Contains(diff, "+      return this.value;") {
		t.Errorf("Diff doesn't contain second added line: %s", diff)
	}

	// Apply edits for real
	_, err = applyFileEdits("test.txt", testFilePath, edits, false)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify file was actually changed
	contentAfterEdit, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after edit: %v", err)
	}
	expectedContent := "class Example {\n  constructor() {\n    this.value = 0;\n  }\n  method() {\n      // increased indent\n      return this.value;\n  }\n}\n"
	if string(contentAfterEdit) != expectedContent {
		t.Errorf("File content doesn't match expected.\nGot: %s\nWant: %s",
			string(contentAfterEdit), expectedContent)
	}
}

func TestApplyFileEdits_MultipleEdits(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	originalContent := "function example() {\n  // First block\n  console.log('hello');\n  \n  // Second block\n  console.log('world');\n}\n"
	err := os.WriteFile(testFilePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Define multiple edits
	edits := []Edit{
		{
			OldText: "// First block\n  console.log('hello');",
			NewText: "// First block - modified\n  console.log('HELLO');",
		},
		{
			OldText: "// Second block\n  console.log('world');",
			NewText: "// Second block - modified\n  console.log('WORLD');\n  console.log('Done!');",
		},
	}

	// Apply edits in dry-run mode
	diff, err := applyFileEdits("test.txt", testFilePath, edits, true)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify diff contains expected changes from both edits
	if !strings.Contains(diff, "-  // First block") {
		t.Errorf("Diff doesn't contain first removed line: %s", diff)
	}
	if !strings.Contains(diff, "+  // First block - modified") {
		t.Errorf("Diff doesn't contain first added line: %s", diff)
	}
	if !strings.Contains(diff, "-  // Second block") {
		t.Errorf("Diff doesn't contain second removed line: %s", diff)
	}
	if !strings.Contains(diff, "+  // Second block - modified") {
		t.Errorf("Diff doesn't contain second added line: %s", diff)
	}
	if !strings.Contains(diff, "+  console.log('Done!');") {
		t.Errorf("Diff doesn't contain additional line: %s", diff)
	}

	// Apply edits for real
	_, err = applyFileEdits("test.txt", testFilePath, edits, false)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify file was actually changed
	contentAfterEdit, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after edit: %v", err)
	}
	expectedContent := "function example() {\n  // First block - modified\n  console.log('HELLO');\n  \n  // Second block - modified\n  console.log('WORLD');\n  console.log('Done!');\n}\n"
	if string(contentAfterEdit) != expectedContent {
		t.Errorf("File content doesn't match expected.\nGot: %s\nWant: %s",
			string(contentAfterEdit), expectedContent)
	}
}

func TestApplyFileEdits_LineEndingNormalization(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test file with Windows line endings
	testFilePath := filepath.Join(tempDir, "test.txt")
	originalContent := "Line one.\r\nLine two.\r\nLine three.\r\n"
	err := os.WriteFile(testFilePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Define an edit with Unix line endings
	edits := []Edit{
		{
			OldText: "Line two.",
			NewText: "Line 2 - modified.",
		},
	}

	// Apply edits
	_, err = applyFileEdits("test.txt", testFilePath, edits, false)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify file was changed correctly
	contentAfterEdit, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after edit: %v", err)
	}

	// The function should preserve the original line endings
	expectedContent := "Line one.\r\nLine 2 - modified.\r\nLine three.\r\n"
	if string(contentAfterEdit) != expectedContent {
		t.Errorf("File content doesn't match expected.\nGot: %s\nWant: %s",
			string(contentAfterEdit), expectedContent)
	}
}

func TestApplyFileEdits_MixedLineEndings(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test file with mixed line endings (more CRLF than LF)
	testFilePath := filepath.Join(tempDir, "mixed.txt")
	originalContent := "Line one.\r\nLine two.\r\nLine three.\nLine four.\r\nLine five.\r\n"
	err := os.WriteFile(testFilePath, []byte(originalContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Define an edit
	edits := []Edit{
		{
			OldText: "Line three.",
			NewText: "Line 3 - modified.",
		},
	}

	// Apply edits
	_, err = applyFileEdits("mixed.txt", testFilePath, edits, false)
	if err != nil {
		t.Fatalf("Failed to apply edits: %v", err)
	}

	// Verify file was changed correctly and uses CRLF (majority rule)
	contentAfterEdit, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read file after edit: %v", err)
	}

	// Should use CRLF for all line endings since that's the majority
	expectedContent := "Line one.\r\nLine two.\r\nLine 3 - modified.\r\nLine four.\r\nLine five.\r\n"
	if string(contentAfterEdit) != expectedContent {
		t.Errorf("File content doesn't match expected or doesn't use correct line endings.\nGot: %s\nWant: %s",
			string(contentAfterEdit), expectedContent)
	}
}
