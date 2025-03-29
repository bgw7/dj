package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bgw7/dj/internal/datastore"
	"github.com/bgw7/dj/internal/restapi"
	"github.com/bgw7/dj/internal/service"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

const shutdownTimeout = 5 * time.Second

func main() {
	// Initialize flags and logger
	var host, port string
	flag.StringVar(&port, "port", "9999", "server port")
	flag.StringVar(&host, "host", "localhost", "server host")
	flag.Parse()

	handler := log.New(os.Stderr)
	logger := slog.New(handler).With("serviceName", "dj-roomba")
	slog.SetDefault(logger)

	// Set up context with cancellation
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Database connection
	conn, err := pgxpool.New(ctx, "") // Default connection string
	if err != nil || conn.Ping(ctx) != nil {
		log.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	var mediaDir string
	var ok bool
	if mediaDir, ok = os.LookupEnv("YT_OUT_DIR"); !ok {
		mediaDir = "/data/data/com.termux/files/home/storage/shared/Termux_Downloader/Youtube"
	}

	// Datastore and channels
	store := datastore.NewDatastore(conn)
	service := service.NewDomainService(ctx, mediaDir, store)

	h := restapi.NewHandler(service)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", h))

	srv := &http.Server{
		Addr:         net.JoinHostPort(host, port),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		Handler:      mux,
	}

	// Error group to manage goroutines
	eg, ctx := errgroup.WithContext(ctx)

	// Start HTTP server in a separate goroutine
	eg.Go(func() error {
		slog.InfoContext(ctx, "Starting HTTP server", "address", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err // Only return an error if it's not a normal shutdown
		}
		return nil
	})

	// Listen for shutdown signals
	<-ctx.Done()
	slog.WarnContext(ctx, "Shutdown signal received")

	// Create a new context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Gracefully shut down the server
	slog.InfoContext(shutdownCtx, "Shutting down HTTP server")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(shutdownCtx, "Error during server shutdown", "error", err)
	}

	// Wait for all goroutines to exit
	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, "Unexpected error during shutdown", "error", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Server shutdown complete")
}
