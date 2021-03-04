package db

import (
	"github.com/jackc/pgx/pgtype"
)

type Historical struct {
	ID         uint        `json:"id"`
	SymbolUuid pgtype.UUID `json:"symbol_uuid"`

	Values pgtype.Float4Array `json:"values"`

	ForDate   pgtype.Timestamptz `json:"for_date"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
