package datastore

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/bgw7/dj/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//go:embed select.tracks_all.sql
var tracksSelectAll string

//go:embed select.track.sql
var tracksSelectOne string

//go:embed select.tracks.sql
var tracksSelect string

//go:embed insert.tracks.sql
var tracksInsert string

//go:embed update.tracks.sql
var tracksUpdate string

//go:embed delete.votes.sql
var votesDelete string

//go:embed insert.votes.sql
var votesInsert string

func (db *Datastore) GetTrackByID(ctx context.Context, id int) (*internal.Track, error) {
	var track internal.Track
	err := db.conn.QueryRow(
		ctx,
		tracksSelectOne,
		id,
	).Scan(
		&track.ID,
		&track.Url,
		&track.Filename,
	)

	if err == pgx.ErrNoRows {
		return nil, internal.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}
	return &track, nil
}

func (db *Datastore) GetTracks(ctx context.Context) ([]internal.Track, error) {
	rows, err := db.conn.Query(ctx, tracksSelectAll)
	if err != nil {
		return nil, fmt.Errorf("get tracks query failed: %w", err)
	}
	tracks, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[internal.Track])
	if err == pgx.ErrNoRows {
		return nil, internal.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetTracks pgx.CollectRows failed: %w", err)
	}
	slog.InfoContext(ctx, "GetTracks", "trackCount", len(tracks))
	return tracks, nil
}
func (db *Datastore) GetNextTrack(ctx context.Context) (*internal.Track, error) {
	rows, err := db.conn.Query(ctx, tracksSelect)
	if err != nil {
		return nil, fmt.Errorf("get tracks query failed: %w", err)
	}
	track, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[internal.Track])
	if err == pgx.ErrNoRows {
		return nil, internal.ErrRecordNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetTracks pgx.CollectRows failed: %w", err)
	}
	slog.InfoContext(ctx, "GetNextTrack", "track", track)
	return &track, nil
}

func (db *Datastore) CreateTrack(ctx context.Context, track *internal.Track) (*internal.Track, error) {
	err := db.conn.QueryRow(
		ctx,
		tracksInsert,
		track.Url,
		track.Filename,
		track.CreatedBy,
		track.CreatedWith,
	).Scan(
		&track.ID,
		&track.HasPlayed,
	)

	if err != nil {
		if data, ok := err.(*pgconn.PgError); ok && data.Code == "23505" {
			return track, fmt.Errorf("%w: %v", internal.ErrUniqueConstraintViolation, err)
		}
		return nil, fmt.Errorf("CreateTrack failed: %w", err)
	}
	slog.InfoContext(ctx, "track created", "track", track)
	return track, nil
}

func (db *Datastore) UpdateTrack(ctx context.Context, track *internal.Track) error {
	_, err := db.conn.Exec(
		ctx,
		tracksUpdate,
		track.ID,
		track.HasPlayed,
	)
	return err
}

func (db *Datastore) DeleteVote(ctx context.Context, id int) error {
	_, err := db.conn.Exec(
		ctx,
		votesDelete,
		id,
	)
	return err
}

func (db *Datastore) CreateVote(ctx context.Context, v *internal.Vote) error {
	_, err := db.conn.Exec(
		ctx,
		votesInsert,
		v.Filename,
		v.Url,
		v.VoterID,
	)
	return err
}
