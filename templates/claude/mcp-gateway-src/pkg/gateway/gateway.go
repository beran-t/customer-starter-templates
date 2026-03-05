package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/docker/mcp-gateway/pkg/catalog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// UserConfig represents the user's flattened configuration
type UserConfig map[string]any

// Gateway holds the application state that can be hot-reloaded
type Gateway struct {
	catalog        catalog.Catalog
	server         *mcp.Server
	pool           *ClientPool
	instructionMap InstructionMap
	userConfigs    map[string]UserConfig
	toolsLoading   sync.Mutex // Held while dynamicallyListTools is running
}

// New creates a new Gateway instance
func New(ctx context.Context, catalogURLs []string, mappingPath string) (*Gateway, error) {
	instructionMap, err := LoadInstructionMap(mappingPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load instruction map: %w", err)
	}

	cat, err := catalog.ReadFrom(ctx, catalogURLs)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalog: %w", err)
	}

	g := &Gateway{
		pool:           NewClientPool(),
		instructionMap: instructionMap,
		catalog:        cat,
	}

	g.server = g.setupMCPServer()

	return g, nil
}

// setupMCPServer creates and configures the MCP server with middleware
func (g *Gateway) setupMCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "e2b-mcp-gateway", Version: "v0.0.1"}, &mcp.ServerOptions{HasTools: true})

	// Add session middleware and tools/list middleware
	server.AddReceivingMiddleware(sessionMiddleware, func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			if method == "tools/list" {
				// Wait for tools to finish loading before proceeding
				func() {
					g.toolsLoading.Lock()
					defer g.toolsLoading.Unlock()
				}()
			}
			return next(ctx, method, req)
		}
	})

	return server
}

// Server returns the MCP server instance
func (g *Gateway) Server() *mcp.Server {
	return g.server
}

// Pool returns the client pool
func (g *Gateway) Pool() *ClientPool {
	return g.pool
}

// LoadConfig loads the configuration from JSON bytes and updates the application state
func (g *Gateway) LoadConfig(ctx context.Context, configJSON []byte) error {
	var userConfigs map[string]UserConfig
	if err := json.Unmarshal(configJSON, &userConfigs); err != nil {
		return fmt.Errorf("failed to parse user configs: %w", err)
	}
	g.userConfigs = userConfigs

	if err := MergeUserConfigsIntoCatalog(g.catalog, g.instructionMap, g.userConfigs); err != nil {
		return fmt.Errorf("failed to merge user configs: %w", err)
	}

	// Dynamically load tools after config is loaded
	go g.dynamicallyListTools(ctx)

	return nil
}
