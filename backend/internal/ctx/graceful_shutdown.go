package ctx

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	ctx         context.Context
	cancel      context.CancelFunc
	mux         sync.Mutex
	initialised = false
)

func GetGracefulShutdownCtx() context.Context {
	mustInitialise()
	return ctx
}

func mustInitialise() {
	mux.Lock()
	defer mux.Unlock()
	if !initialised {
		initialiseContext()
		initialised = true
	}
}

func initialiseContext() {
	// Create Cancel Context for Shutting Down Server
	ctx, cancel = context.WithCancel(context.Background())

	// Create Signal Channel
	sigChan := make(chan os.Signal, 1)
	// Fileter terminate signals
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start Goroutine to Catch Shutdown Signals.
	go func() {
		defer cancel()

		sig := <-sigChan
		// Use default log to reduce package dependency.
		log.Printf("received signal: %s, initiating graceful shutdown", sig.String())

		signal.Reset()
	}()
}
