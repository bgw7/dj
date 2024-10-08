package database

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/la-viajera/reservation-service/internal"
)

type DBIface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
type Database struct {
	conn DBIface
}

func NewDB(connection DBIface) *Database {
	return &Database{
		conn: connection,
	}
}

//go:embed insert_reservation.sql
var insertReservation string

func (db *Database) SearchVenues(ctx context.Context, rs *internal.VenueSearch) ([]internal.Venue, error) {
	return []internal.Venue{}, nil
}
func (db *Database) GetVenues(ctx context.Context) ([]internal.Venue, error) {
	return []internal.Venue{}, nil
}
func (db *Database) CreateVenue(ctx context.Context, r *internal.Venue) (*internal.Venue, error) {
	return nil, nil
}
func (db *Database) CreateReservation(ctx context.Context, r *internal.Reservation) (*internal.Reservation, error) {
	return r, db.conn.QueryRow(
		ctx,
		insertReservation,
		r.ClientID,
		r.VenueID,
		r.StartTimestamp,
		r.EndTimestamp,
		r.Metadata.CreatedBy,
	).Scan(
		&r.PublicID,
		&r.Metadata.CreatedAt,
	)
}

//go:embed update_reservations.sql
var reservationsUpdate string

func (db *Database) UpdateReservation(ctx context.Context, r *internal.Reservation) (*internal.Reservation, error) {
	args := pgx.NamedArgs{
		"@start_timestamp": r.StartTimestamp,
		"@end_timestamp":   r.EndTimestamp,
		"@updated_by":      r.Metadata.UpdateBy,
		"@public_id":       r.PublicID,
	}
	var updatedR internal.Reservation
	err := db.conn.
		QueryRow(ctx, reservationsUpdate, args).
		Scan(&updatedR.Metadata.UpdateAt)

	if err.Error() == pgx.ErrNoRows.Error() {
		return nil, fmt.Errorf("%s ID not found in database: %w", *r.PublicID, internal.RecordNotFoundErr)
	}
	return &updatedR, err

}

//go:embed select_reservations.sql
var reservationsSelect string

func (db *Database) GetReservations(ctx context.Context) ([]internal.Reservation, error) {
	rows, err := db.conn.Query(
		ctx,
		reservationsSelect,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows[internal.Reservation](
		rows, pgx.RowToStructByPos[internal.Reservation])
}

func (db *Database) SearchReservations(ctx context.Context, rs *internal.ReservationSearch) ([]internal.Reservation, error) {
	return nil, nil
}

func (db *Database) FindReservation(ctx context.Context, id string) (*internal.Reservation, error) {

	rows, err := db.conn.Query(
		ctx,
		`SELECT 
			public_id,
			client_id,
			venue_id,
			start_timestamp,
			end_timestamp,
			created_at,
			created_by,
			created_with,
			updated_at,
			updated_by,
			updated_with
		 FROM la_viajera.reservations
			WHERE public_id=$1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	r, err := pgx.CollectExactlyOneRow[internal.Reservation](rows, pgx.RowToStructByName[internal.Reservation])

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("reservation not found for ID: %s %w", id, internal.RecordNotFoundErr)
	}
	return &r, err
}
