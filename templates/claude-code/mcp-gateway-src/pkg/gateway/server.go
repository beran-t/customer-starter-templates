package gateway

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// SetupHTTPServer creates and configures the HTTP server with all routes
func SetupHTTPServer(mcpServer *mcp.Server, addr string, token string) *http.Server {
	// Create MCP protocol handler
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return mcpServer
	}, nil)

	// Create mux for multiple endpoints
	mux := http.NewServeMux()

	// Health check endpoint (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// MCP protocol endpoint with auth middleware if token is provided
	mcpHandler := http.Handler(handler)
	if token != "" {
		mcpHandler = authMiddleware(token)(mcpHandler)
		zap.L().Info("authentication enabled")
	}
	mux.Handle("/mcp", mcpHandler)

	// Add global middleware (logging and CORS, but not auth)
	var finalHandler http.Handler = mux
	finalHandler = loggingHandler(finalHandler)
	finalHandler = corsHandler(finalHandler)

	return &http.Server{
		Addr:    addr,
		Handler: finalHandler,
	}
}

// Run starts the HTTP server and handles graceful shutdown
func Run(srv *http.Server, clientPool *ClientPool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("Server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-sigChan

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close pool first
	if err := clientPool.Close(); err != nil {
		zap.L().Error("Error closing pool", zap.Error(err))
	}

	// Shutdown HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("Error shutting down HTTP server", zap.Error(err))
	}
}
