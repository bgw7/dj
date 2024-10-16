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
	logger := slog.New(handler).With("serviceName", "dj-roomba")
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
	eg, ctx := errgroup.WithContext(ctx)

	lisenOnTextMsgs(ctx, srv, eg)
	listenAndServe(ctx, mux, eg)
	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, "errgroup err", "error", err)
		os.Exit(1)
	}
}

func lisenOnTextMsgs(ctx context.Context, srv *service.DomainService, eg *errgroup.Group) {
	eg.Go(func() error {
		return srv.RunSmsPoller(ctx)
	})
	eg.Go(func() error {
		return srv.RunPlayNext(ctx)INSERT INTO tracks (
			id,
			url,
			url,
			filename,
			has_played,
			created_by,
			created_at
		  )
		VALUES (
			id:integer,
			'url:character varying',
			'url:character varying',
			'filename:character varying',
			has_played:boolean,
			'created_by:character varying',
			'created_at:timestamp without time zone'
		  );
	})
	eg.Go(func() error {
		<-ctx.Done()
		slog.WarnContext(
			ctx,
			"shutting down sms poller",
			"contextErr",
			ctx.Err(),
		)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return termux.MediaStop(shutdownCtx)
	})
}

func listenAndServe(ctx context.Context, h *http.ServeMux, eg *errgroup.Group) {
	s := &http.Server{
		Addr:    net.JoinHostPort(c.Host, c.Port),
		Handler: h,
	}
	eg.Go(func() error {
		slog.InfoContext(ctx, "starting http server", "address", s.Addr)
		return s.ListenAndServe()
	})

	eg.Go(func() error {
		<-ctx.Done()
		slog.WarnContext(
			ctx,
			"shutting down http server",
			"contextErr",
			ctx.Err(),
		)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return s.Shutdown(shutdownCtx)
	})
}
