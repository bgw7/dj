package internal

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Reservation struct {
	ID             *string           `json:"-" db:"-"`
	PublicID       *string           `json:"publicId" db:"public_id"`
	ClientID       *string           `json:"clientId" db:"client_id"`
	VenueID        *string           `json:"venueId" db:"venue_id"`
	StartTimestamp *pgtype.Timestamp `json:"startTimestamp" db:"start_timestamp"`
	EndTimestamp   *pgtype.Timestamp `json:"endTimestamp" db:"end_timestamp"`
	Metadata
}

type Metadata struct {
	CreatedAt   pgtype.Timestamp  `json:"createdAt" db:"created_at"`
	CreatedBy   string            `json:"createdBy" db:"created_by"`
	CreatedWith string            `json:"-" db:"created_with"`
	UpdateAt    *pgtype.Timestamp `json:"updateAt,omitempty" db:"updated_at"`
	UpdateBy    *string           `json:"updateBy,omitempty" db:"updated_by"`
	UpdateWith  *string           `json:"-" db:"updated_with"`
}

type ReservationSearch struct {
	ClientID       *string `json:"clientId"`
	StartTimestamp *string `json:"startTimestamp"`
	Limit          int     `json:"limit"`
	Offset         int     `json:"offset"`
}
