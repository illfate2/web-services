package repo

import (
	"context"

	"github.com/illfate2/web-services/client-server-with-auth/api-server-with-auth/pkg/entities"
)

func (r *Repo) InsertPerson(person entities.Person) (entities.Person, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO persons(first_name,second_name,middle_name) VALUES($1,$2,$3) RETURNING id`,
		person.FirstName, person.LastName, person.MiddleName).
		Scan(&person.ID)
	if err != nil {
		return entities.Person{}, err
	}
	return person, nil
}

func (r *Repo) FindPerson(id int) (entities.Person, error) {
	var p entities.Person
	err := r.conn.QueryRow(context.TODO(),
		`SELECT id, first_name, second_name, middle_name FROM persons WHERE id = $1`, id).
		Scan(&p.ID, &p.FirstName, &p.LastName, &p.MiddleName)
	return p, err
}
