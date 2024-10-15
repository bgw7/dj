package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/bgw7/dj/internal"
	"github.com/jackc/pgx/v5"
)

//go:embed select.tracks.sql
var tracksSelect string

//go:embed insert.tracks.sql
var tracksInsert string

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
	t, err := db.conn.Exec(
		ctx,
		tracksInsert,
		track.Url,
		track.Filename,
		track.CreatedBy,
	)
	fmt.Print(t)
	return err
}

func (db *Database) DeleteVote(ctx context.Context, trackId int, userId string) error {
	_, err := db.conn.Exec(
		ctx,
		votesDelete,
		trackId,
		userId,
	)
	return err
}

func (db *Database) CreateVote(ctx context.Context, trackId int, userId string) error {
	_, err := db.conn.Exec(ctx, votesInsert, trackId, userId)
	return err
}
