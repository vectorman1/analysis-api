package entities

import (
	"github.com/jackc/pgx/pgtype"
)

type PrivateRole uint

const (
	Default PrivateRole = iota
)

type User struct {
	ID          uint
	Uuid        pgtype.UUID
	PrivateRole PrivateRole
	Username    string
	Password    string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	DeletedAt   pgtype.Timestamptz
}
