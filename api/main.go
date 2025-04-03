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

const shutdownTimeout = 7 * time.Second

func main() {
	// Initialize flags and logger
	var host, port string
	flag.StringVar(&port, "port", "9999", "server port")
	flag.StringVar(&host, "host", "localhost", "server host")
	flag.Parse()

	handler := log.New(os.Stderr)
	logger := slog.New(handler).With("serviceName", "dj-roomba")
	slog.SetDefault(logger)

	// Create a parent context with CancelCauseFunc
	rootCtx, cancelWithCause := context.WithCancelCause(context.Background())
	defer cancelWithCause(nil) // Ensure cleanup

	// Listen for termination signals
	signalCtx, signalCancel := signal.NotifyContext(rootCtx, syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	// Database connection
	conn, err := pgxpool.New(signalCtx, "")
	if err != nil {
		log.Error("Database connection initialization failed", "error", err)
		cancelWithCause(err) // Attach error cause
		os.Exit(1)
	}

	if pingErr := conn.Ping(signalCtx); pingErr != nil {
		log.Error("Database ping failed", "error", pingErr)
		cancelWithCause(pingErr) // Attach error cause
		conn.Close()             // Ensure proper cleanup before exit
		os.Exit(1)
	}
	defer conn.Close()

	// Set up media directory
	mediaDir, ok := os.LookupEnv("MEDIA_DIR")
	if !ok || mediaDir == "" {
		mediaDir = "/data/data/com.termux/files/home/storage/shared/Termux_Downloader/Youtube"
	}

	slog.InfoContext(signalCtx, "Media Directory Set", "mediaDirLocation", mediaDir)

	// Datastore and service initialization
	store := datastore.NewDatastore(conn)
	service := service.NewDomainService(signalCtx, mediaDir, store)

	h := restapi.NewHandler(service, mediaDir)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", h))

	srv := &http.Server{
		Addr:         net.JoinHostPort(host, port),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 20 * time.Second,
		Handler:      mux,
	}

	// Error group to manage goroutines
	eg, ctx := errgroup.WithContext(signalCtx)

	// Start HTTP server in a separate goroutine
	eg.Go(func() error {
		slog.InfoContext(ctx, "Starting HTTP server", "address", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			cancelWithCause(err) // Attach error cause to cancellation
			return err
		}
		return nil
	})

	// Wait for shutdown signal
	<-signalCtx.Done()
	slog.WarnContext(ctx, "Shutdown signal received", "cause", context.Cause(signalCtx))

	// Create a timeout context for graceful shutdown
	gracefulCtx, gracefulCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer gracefulCancel()

	// Gracefully shut down the server
	slog.WarnContext(gracefulCtx, "Shutting down HTTP server")
	if err := srv.Shutdown(gracefulCtx); err != nil {
		slog.ErrorContext(gracefulCtx, "Error during server shutdown", "error", err)
	}

	// Wait for all goroutines to finish before final logging
	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, "Unexpected error during shutdown", "error", err)
	}

	slog.InfoContext(ctx, "Shutdown complete")
}
