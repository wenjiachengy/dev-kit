package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/dev-kit/tools"
)

func main() {
	envFile := flag.String("env", ".env", "Path to environment file")
	protocol := flag.String("protocol", "stdio", "Protocol to use (stdio, sse)")
	flag.Parse()

	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil {
			fmt.Printf("Warning: Error loading env file %s: %v\n", *envFile, err)
		}
	}

	mcpServer := server.NewMCPServer(
		"Dev Kit",
		"1.0.0",
		server.WithLogging(),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	enableTools := strings.Split(os.Getenv("ENABLE_TOOLS"), ",")
	allToolsEnabled := len(enableTools) == 1 && enableTools[0] == ""

	isEnabled := func(toolName string) bool {
		return allToolsEnabled || slices.Contains(enableTools, toolName)
	}

	if isEnabled("confluence") {
		tools.RegisterConfluenceTool(mcpServer)
	}

	if isEnabled("jira") {
		tools.RegisterJiraTool(mcpServer)
	}

	if isEnabled("gitlab") {
		tools.RegisterGitLabTool(mcpServer)
	}

	if isEnabled("github") {
		tools.RegisterGitHubTool(mcpServer)
	}

	if isEnabled("script") {
		tools.RegisterScriptTool(mcpServer)
	}

	if isEnabled("codereview") {
		tools.RegisterCodeReviewTool(mcpServer)
	}

	if *protocol == "stdio" {

		if err := server.ServeStdio(mcpServer); err != nil {
			panic(fmt.Sprintf("Server error: %v", err))
		}
	} else if *protocol == "sse" {
		port := os.Getenv("PORT")
		if port == "" {	
			port = "8080"
		}

		sseServer := server.NewSSEServer(mcpServer)
		if err := sseServer.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
