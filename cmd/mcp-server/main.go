package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/meshbrow-dev/meshbrow-mcp/internal/mcp"
	"github.com/meshbrow-dev/meshbrow-mcp/internal/tools"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	mode := flag.String("mode", "stdio", "Transport mode: stdio or ws")
	port := flag.Int("port", 9090, "WebSocket listen port (ws mode only)")
	apiURL := flag.String("api-url", "", "Meshbrow API URL (default: https://api.meshbrow.dev)")
	apiKey := flag.String("api-key", "", "Meshbrow API key (or set MESHBROW_API_KEY)")
	logLevel := flag.String("log-level", "info", "Log level: debug, info, warn, error")
	version := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *version {
		fmt.Printf("meshbrow-mcp %s\n  commit: %s\n  built:  %s\n", Version, Commit, Date)
		os.Exit(0)
	}

	// Configure logging
	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	// Resolve API key
	key := *apiKey
	if key == "" {
		key = os.Getenv("MESHBROW_API_KEY")
	}
	if key == "" {
		slog.Error("MESHBROW_API_KEY is required (flag --api-key or environment variable)")
		os.Exit(1)
	}

	// Resolve API URL
	url := *apiURL
	if url == "" {
		url = os.Getenv("MESHBROW_API_URL")
	}
	if url == "" {
		url = "https://api.meshbrow.dev"
	}

	// Create tool registry
	registry := tools.NewRegistry(url, key)

	// Create and run server
	server := mcp.NewServer(registry, Version)

	switch *mode {
	case "stdio":
		slog.Info("starting meshbrow-mcp", "mode", "stdio", "version", Version)
		if err := server.ServeStdio(); err != nil {
			slog.Error("stdio server error", "error", err)
			os.Exit(1)
		}
	case "ws":
		slog.Info("starting meshbrow-mcp", "mode", "ws", "port", *port, "version", Version)
		if err := server.ServeWebSocket(*port); err != nil {
			slog.Error("websocket server error", "error", err)
			os.Exit(1)
		}
	default:
		slog.Error("invalid mode", "mode", *mode)
		os.Exit(1)
	}
}
