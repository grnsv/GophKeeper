package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/client/app"
	"github.com/grnsv/GophKeeper/internal/client/config"
	"github.com/grnsv/GophKeeper/internal/client/service"
	"github.com/grnsv/GophKeeper/internal/client/storage"
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

func fatalIfErr(prefix string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", prefix, err)
		os.Exit(1)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Fprintf(os.Stdout, "Version: %s\nBuild date: %s\n", buildVersion, buildDate)
			os.Exit(0)
		}
	}

	cfg, err := config.Parse()
	fatalIfErr("config error", err)

	security := service.NewSecuritySource()
	client, err := api.NewClient(cfg.ServerAddress, security)
	fatalIfErr("client error", err)

	srv := service.New(client, security, storage.New)
	defer srv.Close()

	_, err = tea.NewProgram(app.New(srv, buildVersion, buildDate), tea.WithContext(ctx)).Run()
	fatalIfErr("program error", err)
}
