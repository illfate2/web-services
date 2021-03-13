package repo

import (
	"context"

	"github.com/illfate2/web-services/server-with-html-serve/pkg/entities"
)

func (r *Repo) InsertMuseumSet(set entities.MuseumSet) (entities.MuseumSet, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO museum_item_sets(name)
		VALUES($1)
		ON CONFLICT (name)
		DO UPDATE SET name=EXCLUDED.name
		RETURNING id`,
		set.Name).
		Scan(&set.ID)
	if err != nil {
		return entities.MuseumSet{}, err
	}
	return set, nil
}

func (r *Repo) FindMuseumSets() ([]entities.MuseumSet, error) {
	rows, err := r.conn.Query(context.Background(),
		`SELECT 
			id, name
			FROM museum_item_sets
`)
	if err != nil {
		return nil, err
	}
	var sets []entities.MuseumSet
	for rows.Next() {
		var set entities.MuseumSet
		err := rows.Scan(
			&set.ID, &set.Name,
		)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}
	return sets, nil
}

func (r *Repo) FindMuseumSetByName(name string) (entities.MuseumSet, error) {
	var set entities.MuseumSet
	err := r.conn.QueryRow(context.TODO(), `SELECT id,name FROM museum_item_sets WHERE name = $1`, name).
		Scan(&set.ID, &set.Name)
	return set, err
}

func (r *Repo) FindMuseumSet(id int) (entities.MuseumSetWithDetails, error) {
	var curSet entities.MuseumSetWithDetails
	err := r.conn.QueryRow(context.Background(),
		`SELECT 
      mis.id, mis.name
      FROM museum_item_sets mis
	WHERE mis.id = $1`, id).Scan(&curSet.ID, &curSet.Name)
	if err != nil {
		return entities.MuseumSetWithDetails{}, err
	}
	return curSet, nil
}

func (r *Repo) UpdateMuseumSet(set entities.MuseumSet) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE museum_item_sets 
			SET name = $1
			WHERE id = $2`,
		set.Name, set.ID)
	return err
}

func (r *Repo) DeleteMuseumSet(id int) error {
	_, err := r.conn.Exec(context.TODO(), `DELETE FROM museum_item_sets WHERE id = $1`, id)
	return err
}
