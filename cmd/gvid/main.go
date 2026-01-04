// Command gvid is the Gastown Viewer Intent daemon.
// It provides an HTTP API for querying Beads issue data.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/intent-solutions-io/gastown-viewer-intent/internal/api"
	"github.com/intent-solutions-io/gastown-viewer-intent/internal/beads"
)

// version is set by goreleaser ldflags at build time
var version = "dev"

func main() {
	// Parse flags
	port := flag.Int("port", 7070, "HTTP server port")
	host := flag.String("host", "localhost", "HTTP server host")
	workDir := flag.String("dir", "", "Working directory (default: current directory)")
	townRoot := flag.String("town", "", "Gas Town workspace root (default: ~/gt)")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("gvid version %s\n", version)
		os.Exit(0)
	}

	// Create beads adapter
	adapter := beads.NewCLIAdapter(*workDir)

	// Create server config
	config := api.DefaultConfig()
	config.Port = *port
	config.Host = *host
	config.Version = version
	config.TownRoot = *townRoot

	// Create and start server
	server := api.NewServer(config, adapter)

	// Handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-done
		log.Println("Shutting down...")
		os.Exit(0)
	}()

	// Start server
	log.Printf("Gastown Viewer Intent daemon v%s", version)
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
