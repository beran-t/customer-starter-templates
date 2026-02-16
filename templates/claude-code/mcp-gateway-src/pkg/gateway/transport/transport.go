package transport

import (
	"context"

	"github.com/docker/mcp-gateway/pkg/catalog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Transport defines the interface for creating MCP client sessions.
// Different implementations support different connection types (Docker, Remote, etc.)
//
// To add a new transport type:
//  1. Create a new type that implements this interface
//  2. Add a case in GetTransport()
//  3. Example: see docker.go and remote.go
type Transport interface {
	// CreateSession establishes a new MCP session using this transport
	CreateSession(ctx context.Context, client *mcp.Client, server catalog.Server, serverName string) (*mcp.ClientSession, error)
}

// GetTransport returns the appropriate transport implementation for the given server type
func GetTransport(serverType string) Transport {
	switch serverType {
	case "remote":
		return &RemoteTransport{}
	case "github":
		return &GitHubTransport{}
	default:
		return &DockerTransport{}
	}
}
