package transport

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/mcp-gateway/pkg/catalog"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RemoteTransport creates MCP sessions by connecting to remote servers via HTTP
type RemoteTransport struct{}

// CreateSession creates an MCP session by connecting to a remote server (SSE or Streamable HTTP)
func (t *RemoteTransport) CreateSession(ctx context.Context, client *mcp.Client, server catalog.Server, serverName string) (*mcp.ClientSession, error) {
	// Get URL and transport type
	url := server.Remote.URL
	transport := server.Remote.Transport
	if url == "" && server.SSEEndpoint != "" {
		// Fallback to deprecated SSEEndpoint
		url = server.SSEEndpoint
		transport = "sse"
	}
	if url == "" {
		return nil, fmt.Errorf("no remote URL configured for %s", serverName)
	}

	// Create HTTP client with custom headers if any
	httpClient := &http.Client{}
	if len(server.Remote.Headers) > 0 {
		httpClient.Transport = &headerRoundTripper{
			base:    http.DefaultTransport,
			headers: server.Remote.Headers,
		}
	}

	// Select transport based on type
	var mcpTransport mcp.Transport
	switch strings.ToLower(transport) {
	case "sse":
		mcpTransport = &mcp.SSEClientTransport{
			Endpoint:   url,
			HTTPClient: httpClient,
		}
	case "http", "streamable", "streaming", "streamable-http":
		mcpTransport = &mcp.StreamableClientTransport{
			Endpoint:   url,
			HTTPClient: httpClient,
		}
	default:
		return nil, fmt.Errorf("unsupported remote transport: %s", transport)
	}

	// Connect to remote server
	session, err := client.Connect(ctx, mcpTransport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to remote %s: %w", serverName, err)
	}

	return session, nil
}

// headerRoundTripper is an http.RoundTripper that adds custom headers to all requests
type headerRoundTripper struct {
	base    http.RoundTripper
	headers map[string]string
}

// RoundTrip implements http.RoundTripper
func (h *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	newReq := req.Clone(req.Context())
	// Add custom headers
	for key, value := range h.headers {
		// Don't override Accept header if already set by streamable transport
		if key == "Accept" && newReq.Header.Get("Accept") != "" {
			continue
		}
		newReq.Header.Set(key, value)
	}
	return h.base.RoundTrip(newReq)
}
