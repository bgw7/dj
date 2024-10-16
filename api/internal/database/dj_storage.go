package database

import (
	"context"
	_ "embed"

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
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[internal.Track])
}

func (db *Database) CreateTrack(ctx context.Context, track *internal.Track) error {
	err := db.conn.QueryRow(
		ctx,
		tracksInsert,
		track.Url,
		track.Filename,
		track.CreatedBy,
	).Scan()
	if err != nil {
		if data, ok := err.(*pgconn.PgError); ok && data.Code == "23505" {
			return db.CreateVote(ctx, track.Url, track.CreatedBy)
		}
		return err

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

func (db *Database) DeleteVote(ctx context.Context, url, userId string) error {
	_, err := db.conn.Exec(
		ctx,
		votesDelete,
		url,
		userId,
	)
	return err
}

func (db *Database) CreateVote(ctx context.Context, url, userId string) error {
	_, err := db.conn.Exec(ctx, votesInsert, url, userId)
	return err
}