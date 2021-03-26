package repo

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"xorm.io/xorm"
)

type Repo struct {
	conn   *pgxpool.Pool
	engine *xorm.Engine
}

func New(conn *pgxpool.Pool, engine *xorm.Engine) *Repo {
	engine.Sync(Auth{})
	return &Repo{conn: conn, engine: engine}
}
