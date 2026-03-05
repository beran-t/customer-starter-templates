package gateway

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/docker/mcp-gateway/pkg/catalog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// dynamicallyListTools discovers and registers tools from all configured MCP servers in parallel
func (g *Gateway) dynamicallyListTools(ctx context.Context) {
	// Lock to signal tools are loading
	g.toolsLoading.Lock()
	defer g.toolsLoading.Unlock()

	// Extract session ID from context
	sessionID := getSessionID(ctx)

	// Create errgroup with context and limit concurrency to numCPU * 2
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU() * 2)

	for configKey := range g.userConfigs {
		// Check if this is a GitHub server (prefix: github/)
		if strings.HasPrefix(configKey, "github/") {
			if catalogServer, ok := g.buildGitHubServer(configKey); ok {
				serverName := configKey
				eg.Go(func() error {
					return discoverAndRegisterTools(ctx, g.pool, g.server, serverName, sessionID, catalogServer)
				})
			}
			continue
		}

		// Handle catalog servers: map config key to actual catalog server name
		actualServerName, ok := GetServerNameFromInstructions(g.instructionMap, configKey)
		if !ok {
			continue
		}

		cServer, ok := g.catalog.Servers[actualServerName]
		if !ok {
			continue
		}

		// Capture loop variables for goroutine
		serverName := actualServerName
		catalogServer := cServer

		eg.Go(func() error {
			return discoverAndRegisterTools(ctx, g.pool, g.server, serverName, sessionID, catalogServer)
		})
	}

	// Wait for all goroutines to complete
	// Individual failures are logged within each goroutine and don't affect others
	if err := eg.Wait(); err != nil {
		zap.L().Error("Some tools failed to load", zap.String("component", "TOOLS"), zap.Error(err))
	}
}

// buildGitHubServer creates a catalog.Server for a GitHub-based MCP server
func (g *Gateway) buildGitHubServer(serverName string) (catalog.Server, bool) {
	userConfig := g.userConfigs[serverName]

	// Extract install and run commands from user config
	var installCmd, runCmd string
	if install, ok := userConfig["installCmd"].(string); ok {
		installCmd = install
	}
	if run, ok := userConfig["runCmd"].(string); ok {
		runCmd = run
	}

	if runCmd == "" {
		zap.L().Warn("GitHub server missing run command",
			zap.String("component", "TOOLS"),
			zap.String("server", serverName))
		return catalog.Server{}, false
	}

	envs := []catalog.Env{}
	if envs, ok := userConfig["envs"].([]any); ok {
		for _, env := range envs {
			if name, ok := env.(string); ok {
				envs = append(envs, catalog.Env{Name: name, Value: name})
			}
		}
	}
	// Create ad-hoc catalog.Server for GitHub transport
	catalogServer := catalog.Server{
		Type:      "github",
		Name:      serverName,
		Command:   []string{runCmd},
		LongLived: false,
		Env:       envs,
	}

	// Store install command in Env as special marker
	if installCmd != "" {
		catalogServer.Env = []catalog.Env{
			{Name: "INSTALL_COMMAND", Value: installCmd},
		}
	}

	return catalogServer, true
}

// discoverAndRegisterTools discovers and registers tools for a single MCP server
func discoverAndRegisterTools(ctx context.Context, clientPool *ClientPool, server *mcp.Server, serverName string, sessionID string, catalogServer catalog.Server) error {
	session, err := clientPool.Acquire(ctx, serverName, sessionID, catalogServer)
	if err != nil {
		zap.L().Error("Failed to acquire session", zap.String("component", "TOOLS"), zap.String("server", serverName), zap.Error(err))
		return nil // Don't fail the entire operation
	}
	defer clientPool.Release(serverName, sessionID)

	tools, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		zap.L().Error("Failed to list tools", zap.String("component", "TOOLS"), zap.String("server", serverName), zap.Error(err))
		return nil // Don't fail the entire operation
	}

	toolHandler := createToolHandler(clientPool, serverName, catalogServer)

	for _, tool := range tools.Tools {
		tool.Name = fmt.Sprintf("%s-%s", serverName, tool.Name)
		server.AddTool(tool, toolHandler)
	}

	return nil
}

// createToolHandler creates a handler function for tool calls
func createToolHandler(clientPool *ClientPool, serverName string, catalogServer catalog.Server) func(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, params *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sessionID := getSessionID(ctx)

		session, err := clientPool.Acquire(ctx, serverName, sessionID, catalogServer)
		if err != nil {
			return &mcp.CallToolResult{}, fmt.Errorf("failed to acquire session: %w", err)
		}
		defer clientPool.Release(serverName, sessionID)

		return session.CallTool(ctx, &mcp.CallToolParams{
			Arguments: params.Params.Arguments,
			Name:      strings.TrimPrefix(params.Params.Name, serverName+"-"),
		})
	}
}
