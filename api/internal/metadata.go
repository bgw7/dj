package internal

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Metadata struct {
	CreatedAt   pgtype.Timestamp  `json:"createdAt" db:"created_at"`
	CreatedBy   string            `json:"createdBy" db:"created_by"`
	CreatedWith string            `json:"-" db:"created_with"`
	UpdateAt    *pgtype.Timestamp `json:"updateAt,omitempty" db:"updated_at"`
	UpdateBy    *string           `json:"updateBy,omitempty" db:"updated_by"`
	UpdateWith  *string           `json:"-" db:"updated_with"`
}
