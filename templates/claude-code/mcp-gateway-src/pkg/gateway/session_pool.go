package gateway

import (
	"context"
	"fmt"
	"sync"

	"e2b.dev/mcp-gateway/pkg/gateway/transport"
	"github.com/docker/mcp-gateway/pkg/catalog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// ClientPool is a thread-safe pool of MCP client sessions
type ClientPool struct {
	mu        sync.RWMutex
	sessions  map[string]*mcp.ClientSession
	longLived map[string]bool // tracks which sessions are long-lived
}

// NewClientPool creates a new client pool
func NewClientPool() *ClientPool {
	return &ClientPool{
		sessions:  make(map[string]*mcp.ClientSession),
		longLived: make(map[string]bool),
	}
}

// Acquire gets or creates an MCP session for the given key and session ID
func (p *ClientPool) Acquire(ctx context.Context, mcpKey string, sessionID string, cServer catalog.Server) (*mcp.ClientSession, error) {
	key := fmt.Sprintf("%s:%s", mcpKey, sessionID)

	// Check if session already exists
	p.mu.RLock()
	if session, ok := p.sessions[key]; ok {
		p.mu.RUnlock()
		return session, nil
	}
	p.mu.RUnlock()

	// Create new session using appropriate transport
	session, err := p.createSession(ctx, mcpKey, cServer)
	if err != nil {
		zap.L().Error("Failed to create session",
			zap.String("component", "POOL"),
			zap.String("key", key),
			zap.Error(err))
		return nil, err
	}

	// Store in pool
	p.mu.Lock()
	p.sessions[key] = session
	if cServer.LongLived {
		p.longLived[key] = true
	}
	p.mu.Unlock()

	return session, nil
}

// createSession creates a new MCP session using the appropriate transport
func (p *ClientPool) createSession(ctx context.Context, serverName string, server catalog.Server) (*mcp.ClientSession, error) {
	// Create MCP client
	client := mcp.NewClient(&mcp.Implementation{
		Name: server.Name,
	}, nil)

	// Get appropriate transport and create session
	t := transport.GetTransport(server.Type)
	return t.CreateSession(ctx, client, server, serverName)
}

// Release removes and closes a session from the pool
// Long-lived sessions are kept alive and not closed
func (p *ClientPool) Release(mcpKey string, sessionID string) error {
	key := fmt.Sprintf("%s:%s", mcpKey, sessionID)

	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if this is a long-lived session
	if p.longLived[key] {
		// Don't delete or close long-lived sessions
		return nil
	}

	if session, ok := p.sessions[key]; ok {
		delete(p.sessions, key)
		return session.Close()
	}
	return nil
}

// Close closes all sessions in the pool, including long-lived ones
// This should be called during graceful shutdown
func (p *ClientPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var firstErr error
	for key, session := range p.sessions {
		if err := session.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		delete(p.sessions, key)
		delete(p.longLived, key)
	}

	return firstErr
}
