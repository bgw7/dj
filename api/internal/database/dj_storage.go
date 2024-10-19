package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/bgw7/dj/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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

func (db *Database) ListTracks(ctx context.Context) ([]internal.Track, error) {
	rows, err := db.conn.Query(ctx, tracksSelect)
	if err != nil {
		return nil, fmt.Errorf("list tracks query failed: %w", err)
	}
	tracks, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[internal.Track])
	if err != nil {
		return nil, fmt.Errorf("ListTracks pgx.CollectRows failed: %w", err)
	}
	return tracks, nil
}

func (db *Database) CreateTrack(ctx context.Context, track *internal.Track) error {
	err := db.conn.QueryRow(
		ctx,
		tracksInsert,
		track.Url,
		track.Filename,
		track.CreatedBy,
	).Scan(
		&track.ID,
		&track.HasPlayed,
	)

	if err != nil {
		if data, ok := err.(*pgconn.PgError); ok && data.Code == "23505" && !track.HasPlayed {
			return db.CreateVote(ctx, track.ID, track.CreatedBy)
		}
		return fmt.Errorf("CreateTrack failed: %w", err)
	}
	return err
}

func (db *Database) UpdateTrack(ctx context.Context, track *internal.Track) error {
	_, err := db.conn.Exec(
		ctx,
		tracksUpdate,
		track.ID,
		track.HasPlayed,
	)
	return err
}

func (db *Database) DeleteVote(ctx context.Context, trackID int, userId string) error {
	_, err := db.conn.Exec(
		ctx,
		votesDelete,
		trackID,
		userId,
	)
	return err
}

func (db *Database) CreateVote(ctx context.Context, trackID int, userId string) error {
	_, err := db.conn.Exec(ctx, votesInsert, trackID, userId)
	return err
}
