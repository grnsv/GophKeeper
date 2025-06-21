package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grnsv/GophKeeper/internal/server/app"
)

// buildVersion is set at compile time using -ldflags.
// Example:
//
//	go build -ldflags "-X 'main.buildVersion=1.0.0'"
var buildVersion string

// buildDate is the date of the build, injected via -ldflags.
// Example:
//
//	go build -ldflags "-X 'main.buildDate=2025-05-02'"
var buildDate string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := app.New(ctx, buildVersion, buildDate)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	go app.Run()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
