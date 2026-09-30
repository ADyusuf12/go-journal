package main

import (
	"context"
	"fmt"
)

// Define custom unexported types for context keys to prevent collisions
type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
)

func logServiceCall(ctx context.Context, action string) {
	// Extract values from context
	reqID, ok1 := ctx.Value(requestIDKey).(string)
	userID, ok2 := ctx.Value(userIDKey).(string)

	if !ok1 || !ok2 {
		fmt.Println("[Log] Missing required request context metadata!")
		return
	}

	fmt.Printf("[Log] Action: %s | RequestID: %s | UserID: %s\n", action, reqID, userID)
}

func main() {
	// Start with background context
	ctx := context.Background()

	// Inject request-scoped values into context
	ctx = context.WithValue(ctx, requestIDKey, "REQ-88301")
	ctx = context.WithValue(ctx, userIDKey, "USR-4029")

	// Pass context down to downstream logger function
	logServiceCall(ctx, "UPDATE_INVENTORY")
}
