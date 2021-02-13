package repo

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

func New(conn *pgxpool.Pool) *Repo {
	return &Repo{
		conn: conn,
	}
}

type Repo struct {
	conn *pgxpool.Pool
}
