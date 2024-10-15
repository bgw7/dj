package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/bgw7/dj/internal/database"
	"github.com/bgw7/dj/internal/restapi"
	"github.com/bgw7/dj/internal/service"
	"github.com/bgw7/dj/internal/termux"
	"github.com/charmbracelet/log"
)

const shutdownTimeout = 3 * time.Second

type serverConfig struct {
	Host string
	Port string
}

var c serverConfig

func initFlags() {
	handler := log.New(os.Stderr)
	logger := slog.New(handler).With("serviceName", "reservation-service")
	slog.SetDefault(logger)
	flag.StringVar(&c.Port, "port", "9999", "port used in http server's address")
	flag.StringVar(&c.Host, "host", "localhost", "host used in http server's address")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()
	initFlags()

	conn, err := pgxpool.New(ctx, "") //using default var names. see https://www.postgresql.org/docs/current/libpq-envars.html
	if err != nil {
		slog.ErrorContext(ctx, "pgxpool.New() database connection failed", "error", err)
		os.Exit(1)
	}
	err = conn.Ping(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "database ping failed", "error", err)
		os.Exit(1)
	}

	defer conn.Close()

	db := database.NewDB(conn)
	srv := service.NewDomainService(db)
	h := restapi.NewHandler(srv)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", h))

	if err := lisenOnTextMsgs(ctx, srv); err != nil {
		slog.ErrorContext(ctx, "sms poller error", "error", err)
		os.Exit(1)
	}

	if err := listenAndServe(ctx, mux); err != nil {
		slog.ErrorContext(ctx, "listenAndServe() err", "error", err)
		os.Exit(1)
	}
}

func lisenOnTextMsgs(ctx context.Context, srv *service.DomainService) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return srv.RunSmsPoller(ctx)
})
	

	eg.Go(func() error {
		<- ctx.Done()
slog.InfoContext(
			ctx,
			"context is done. shutting sms poller",
			"contextErr",
			ctx.Err(),
			"timeout",
			shutdownTimeout,
		)
shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
defer cancel()
return termux.MediaStop(shutdownCtx)
)

return eg.Wait()


func listenAndServe(ctx context.Context, h *http.ServeMux) error {
	s := &http.Server{
		Addr:    net.JoinHostPort(c.Host, c.Port),
		Handler: h,
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		slog.InfoContext(ctx, "starting http server", "address", s.Addr)
		return s.ListenAndServe()
	})

	eg.Go(func() error {
		<-ctx.Done()
		slog.InfoContext(
			ctx,
			"http server context is done. shutting down server",
			"contextErr",
			ctx.Err(),
			"timeout",
			shutdownTimeout,
		)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return s.Shutdown(shutdownCtx)
	})

	return eg.Wait()
}
