package repo

import (
	"context"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
)

func (r *Repo) InsertUser(user entities.User) (entities.User, error) {
	err := r.conn.QueryRow(context.TODO(),
		`INSERT INTO users(email, password) VALUES($1,$2) RETURNING id`,
		user.Email, user.Password,
	).Scan(&user.ID)
	return user, err
}

func (r Repo) FindUserByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.conn.QueryRow(context.TODO(),
		`SELECT id,email, password FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}

func (r Repo) FindUser(id int) (entities.User, error) {
	var user entities.User
	err := r.conn.QueryRow(context.TODO(),
		`SELECT id,email, password FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.Password)
	return user, err
}
