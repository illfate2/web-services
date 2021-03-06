package repo

import (
	"context"

	"github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/entities"
)

func (r *Repo) InsertMuseumSet(set entities.MuseumSet) (entities.MuseumSet, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO museum_item_sets(name)
		VALUES($1)
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
	rows, err := r.conn.Query(context.Background(),
		`SELECT 
      mis.id, mis.name
      FROM museum_item_sets mis
	WHERE mis.id = $1`, id)
	if err != nil {
		return entities.MuseumSetWithDetails{}, err
	}
	var curSet entities.MuseumSetWithDetails
	for rows.Next() {
		var item entities.MuseumItemWithKeeper
		err := rows.Scan(
			&curSet.ID, &curSet.Name,
			//&item.ID, &item.Name, &item.CreationDate.Time, &item.Annotation,
			//&item.Keeper.ID, &item.Keeper.FirstName, &item.Keeper.LastName, &item.Keeper.MiddleName,
		)
		if err != nil {
			return entities.MuseumSetWithDetails{}, err
		}
		curSet.Items = append(curSet.Items, item)
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
