package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"os"
)

type ToolTestFunc func(t T)

var toolTestMap = map[string]ToolTestFunc{
	"read_file":                readFileTester,
	"read_multiple_files":      readMultipleFilesTester,
	"write_file":               writeFileTester,
	"edit_file":                editFileTester,
	"create_directory":         createDirectoryTester,
	"list_directory":           listDirectoryTester,
	"directory_tree":           directoryTreeTester,
	"move_file":                moveFileTester,
	"search_files":             searchFilesTester,
	"get_file_info":            getFileInfoTester,
	"list_allowed_directories": listAllowedDirectoriesTester,
}

func main() {
	ctx := context.Background()

	// Create a temporary directory for all test.
	// Each test will get a subdirectory of this.
	mainTempDir, err := os.MkdirTemp("", "mcp-test")
	if err != nil {
		fmt.Println("Failed to create temp dir:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(mainTempDir)

	makeClient := func(ctx context.Context, args []string) (*mcp.InitializeResult, client.MCPClient) {
		// Assume the first argument is the command to run,
		// which expects allowed directories as arguments.
		cmd := os.Args[1]
		env := os.Environ()
		cli, err := client.NewStdioMCPClient(cmd, env, args...)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = "1.0"
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "tester-client",
			Version: "1.0.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{
			Roots: &struct {
				ListChanged bool `json:"listChanged,omitempty"`
			}{
				ListChanged: true,
			},
		}
		result, err := cli.Initialize(ctx, initRequest)
		if err != nil {
			log.Fatalf("Initialize failed: %v", err)
		}
		return result, cli
	}

	result, cli := makeClient(ctx, []string{mainTempDir})
	log.Printf("Connected to %s v%s\n", result.ServerInfo.Name, result.ServerInfo.Version)
	tools, err := cli.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	ctx = WithMCPClientFactory(ctx, makeClient)
	tc := NewTestContext(ctx, "", mainTempDir)
	_ = cli.Close()

	anyFailed := false
	for _, tool := range tools.Tools {
		if tt, ok := toolTestMap[tool.Name]; ok {
			tc.Run(tool.Name, tt)
			if tc.Failed() {
				anyFailed = true
			}
		} else {
			log.Printf("Tool %s: not tested\n", tool.Name)
		}
	}
	if anyFailed {
		color.Red("Some tests failed")
		os.Exit(1)
	} else {
		color.Green("All tests passed")
	}
}
