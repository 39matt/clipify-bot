package globalctx

import (
	"context"
	"sync"
	"time"
)

var (
	globalCtx context.Context
	ctxMutex  sync.RWMutex
	once      sync.Once
)

// Initialize sets up the global context with cancellation
func Initialize() (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc

	once.Do(func() {
		ctxMutex.Lock()
		defer ctxMutex.Unlock()

		globalCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	})

	return globalCtx, cancel
}

// Get returns the global context
func Get() context.Context {
	ctxMutex.RLock()
	defer ctxMutex.RUnlock()

	if globalCtx == nil {
		return context.Background()
	}

	return globalCtx
}

func ForRequest() (context.Context, context.CancelFunc) {
	ctx := Get()
	return context.WithTimeout(ctx, 5*time.Minute)
}
