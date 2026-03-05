package gateway

import (
	"context"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type contextKey string

const sessionIDKey contextKey = "sessionID"

// getSessionID extracts or generates a session ID from the context
func getSessionID(ctx context.Context) string {
	if sessionID := ctx.Value(sessionIDKey); sessionID != nil {
		if id, ok := sessionID.(string); ok {
			return id
		}
	}
	// this should really not happen, I think a safe-ish option is to return a random UUID
	return uuid.New().String()
}

// sessionMiddleware extracts session ID from request and adds it to context
func sessionMiddleware(next mcp.MethodHandler) mcp.MethodHandler {
	return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		// Extract session ID from request and add to context
		if session := req.GetSession(); session != nil {
			sessionID := session.ID()
			ctx = context.WithValue(ctx, sessionIDKey, sessionID)
		}
		return next(ctx, method, req)
	}
}
